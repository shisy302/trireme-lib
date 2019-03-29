package podmonitor

import (
	"context"
	errs "errors"
	"time"

	"k8s.io/client-go/tools/record"

	"go.aporeto.io/trireme-lib/common"
	"go.aporeto.io/trireme-lib/monitor/config"
	"go.aporeto.io/trireme-lib/monitor/extractors"
	"go.aporeto.io/trireme-lib/policy"
	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	// ErrHandlePUStartEventFailed is the error sent back if a start event fails
	ErrHandlePUStartEventFailed = errs.New("Aporeto Enforcer start event failed")

	// ErrNetnsExtractionMissing is the error when we are missing a PID or netns path after successful metadata extraction
	ErrNetnsExtractionMissing = errs.New("Aporeto Enforcer missed to extract PID or netns path")

	// ErrHandlePUStopEventFailed is the error sent back if a stop event fails
	ErrHandlePUStopEventFailed = errs.New("Aporeto Enforcer stop event failed")

	// ErrHandlePUDestroyEventFailed is the error sent back if a create event fails
	ErrHandlePUDestroyEventFailed = errs.New("Aporeto Enforcer destroy event failed")
)

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, handler *config.ProcessorConfig, metadataExtractor extractors.PodMetadataExtractor, netclsProgrammer extractors.PodNetclsProgrammer, nodeName string, enableHostPods bool) *ReconcilePod {
	return &ReconcilePod{
		client:            mgr.GetClient(),
		scheme:            mgr.GetScheme(),
		recorder:          mgr.GetRecorder("trireme-pod-controller"),
		handler:           handler,
		metadataExtractor: metadataExtractor,
		netclsProgrammer:  netclsProgrammer,
		nodeName:          nodeName,
		enableHostPods:    enableHostPods,

		// TODO: might move into configuration
		handlePUEventTimeout:   5 * time.Second,
		metadataExtractTimeout: 3 * time.Second,
		netclsProgramTimeout:   2 * time.Second,
	}
}

// addController adds a new Controller to mgr with r as the reconcile.Reconciler
func addController(mgr manager.Manager, r *ReconcilePod, eventsCh <-chan event.GenericEvent) error {
	// Create a new controller
	c, err := controller.New("trireme-pod-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// we use this mapper in both of our event sources
	mapper := &WatchPodMapper{
		client:         mgr.GetClient(),
		nodeName:       r.nodeName,
		enableHostPods: r.enableHostPods,
	}

	// use the our watch pod mapper which filters pods before we reconcile
	if err := c.Watch(
		&source.Kind{Type: &corev1.Pod{}},
		&handler.EnqueueRequestsFromMapFunc{ToRequests: mapper},
	); err != nil {
		return err
	}

	// we pass in a custom channel for events generated by resync
	if err := c.Watch(
		&source.Channel{Source: eventsCh},
		&handler.EnqueueRequestsFromMapFunc{ToRequests: mapper},
	); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcilePod{}

// ReconcilePod reconciles a Pod object
type ReconcilePod struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client            client.Client
	scheme            *runtime.Scheme
	recorder          record.EventRecorder
	handler           *config.ProcessorConfig
	metadataExtractor extractors.PodMetadataExtractor
	netclsProgrammer  extractors.PodNetclsProgrammer
	nodeName          string
	enableHostPods    bool

	metadataExtractTimeout time.Duration
	handlePUEventTimeout   time.Duration
	netclsProgramTimeout   time.Duration
}

// Reconcile reads that state of the cluster for a pod object
func (r *ReconcilePod) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	puID := request.NamespacedName.String()

	// Fetch the corresponding pod object.
	pod := &corev1.Pod{}
	if err := r.client.Get(ctx, request.NamespacedName, pod); err != nil {
		if errors.IsNotFound(err) {
			zap.L().Debug("Pod IsNotFound", zap.String("puID", puID))
			handlePUCtx, handlePUCancel := context.WithTimeout(ctx, r.handlePUEventTimeout)
			defer handlePUCancel()

			// NOTE: We should not need to call Stop first, a Destroy should also do a Stop if necessary.
			//       This should be fixed in the policy engine.
			// try to call stop, but don't fail if that errors
			err := r.handler.Policy.HandlePUEvent(
				handlePUCtx,
				puID,
				common.EventStop,
				policy.NewPURuntimeWithDefaults(),
			)
			if err != nil {
				zap.L().Warn("failed to handle stop event during destroy", zap.String("puID", puID), zap.Error(err))
			}

			// call destroy regardless if stop succeeded
			if err := r.handler.Policy.HandlePUEvent(
				handlePUCtx,
				puID,
				common.EventDestroy,
				policy.NewPURuntimeWithDefaults(),
			); err != nil {
				zap.L().Error("failed to handle destroy event", zap.String("puID", puID), zap.Error(err))
				return reconcile.Result{}, ErrHandlePUDestroyEventFailed
			}
			return reconcile.Result{}, nil
		}
		// Otherwise, we retry.
		return reconcile.Result{}, err
	}

	// abort immediately if this is a HostNetwork pod, but we don't want to activate them
	// NOTE: is already done in the mapper, however, this additional check does not hurt
	if pod.Spec.HostNetwork && !r.enableHostPods {
		zap.L().Debug("Pod is a HostNetwork pod, but enableHostPods is false", zap.String("puID", puID))
		return reconcile.Result{}, nil
	}

	// abort if this pod is terminating / shutting down
	//
	// NOTE: a pod that is terminating, is going to reconcile as well in the PodRunning phase,
	// however, it will have the deletion timestamp set which is an indicator for us that it is
	// shutting down. It means for us, that we don't have to do anything yet. We can safely stop
	// the PU when the phase is PodSucceeded/PodFailed, or delete it once we reconcile on a pod
	// that cannot be found any longer.
	if pod.DeletionTimestamp != nil {
		zap.L().Debug("Pod is terminating, deletion timestamp found", zap.String("puID", puID), zap.String("deletionTimestamp", pod.DeletionTimestamp.String()))
		return reconcile.Result{}, nil
	}

	// every HandlePUEvent call gets done in this context
	handlePUCtx, handlePUCancel := context.WithTimeout(ctx, r.handlePUEventTimeout)
	defer handlePUCancel()

	switch pod.Status.Phase {
	case corev1.PodPending:
		// we will create the PU already, but not start it
		zap.L().Debug("PodPending", zap.String("puID", puID))

		// now try to do the metadata extraction
		extractCtx, extractCancel := context.WithTimeout(ctx, r.metadataExtractTimeout)
		defer extractCancel()
		puRuntime, err := r.metadataExtractor(extractCtx, r.client, r.scheme, pod, false)
		if err != nil {
			zap.L().Warn("failed to extract metadata", zap.String("puID", puID), zap.Error(err))
		}
		if puRuntime == nil {
			puRuntime = policy.NewPURuntimeWithDefaults()
		}

		// now create/update the PU
		if err := r.handler.Policy.HandlePUEvent(
			handlePUCtx,
			puID,
			common.EventUpdate,
			puRuntime,
		); err != nil {
			if policy.IsErrPUNotFound(err) {
				// create the PU if it does not exist yet
				if err := r.handler.Policy.HandlePUEvent(
					handlePUCtx,
					puID,
					common.EventCreate,
					puRuntime,
				); err != nil {
					zap.L().Error("failed to handle create event", zap.String("puID", puID), zap.Error(err))
					return reconcile.Result{}, err
				}
				return reconcile.Result{}, nil
			}
			zap.L().Error("failed to handle update event", zap.String("puID", puID), zap.Error(err))
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil

	case corev1.PodRunning:
		// try to find out if any of the containers have been started yet
		//
		// NOTE: the spec says for PodRunning that at least one container has been starting,
		// The reality shows us though that we need to check the container states as well:
		// the containers might still be starting up, and are not necessarily in the running state
		// Only containers in the running state can guarantee to us that we are going to have access
		// to the namespace with our remote enforcer.
		var started bool
		for _, status := range pod.Status.InitContainerStatuses {
			if status.State.Running != nil {
				started = true
				break
			}
		}
		if !started {
			for _, status := range pod.Status.ContainerStatuses {
				if status.State.Running != nil {
					started = true
					break
				}
			}
		}
		zap.L().Debug("PodRunning", zap.String("puID", puID), zap.Bool("anyContainerStarted", started))
		if started {
			// now do the metadata extraction
			extractCtx, extractCancel := context.WithTimeout(ctx, r.metadataExtractTimeout)
			defer extractCancel()
			puRuntime, err := r.metadataExtractor(extractCtx, r.client, r.scheme, pod, true)
			if err != nil {
				zap.L().Error("failed to extract metadata", zap.String("puID", puID), zap.Error(err))
				r.recorder.Eventf(pod, "Warning", "PUStart", "PU '%s' failed to extract metadata: %s", puID, err.Error())
				return reconcile.Result{}, err
			}

			// if the metadata extractor is missing the PID or nspath, we need to try again
			// we need it for starting the PU. However, only require this if we are not in host network mode.
			// NOTE: this can happen for example if the containers are not in a running state on their own
			if !pod.Spec.HostNetwork && len(puRuntime.NSPath()) == 0 && puRuntime.Pid() == 0 {
				zap.L().Error("Kubernetes thinks a container is running, however, we failed to extract a PID or NSPath with the metadata extractor. Requeueing...", zap.String("puID", puID))
				r.recorder.Eventf(pod, "Warning", "PUStart", "PU '%s' failed to extract netns", puID)
				return reconcile.Result{}, ErrNetnsExtractionMissing
			}

			// now start (and maybe even still create) the PU
			if err := r.handler.Policy.HandlePUEvent(
				handlePUCtx,
				puID,
				common.EventStart,
				puRuntime,
			); err != nil {
				if policy.IsErrPUAlreadyActivated(err) {
					// abort early if this PU has already been activated before
					zap.L().Debug("PU has already been activated", zap.String("puID", puID), zap.Error(err))
				} else {
					zap.L().Error("failed to handle start event", zap.String("puID", puID), zap.Error(err))
					r.recorder.Eventf(pod, "Warning", "PUStart", "PU '%s' failed to start: %s", puID, err.Error())
					return reconcile.Result{Requeue: true, RequeueAfter: 100 * time.Millisecond}, ErrHandlePUStartEventFailed
				}
			} else {
				r.recorder.Eventf(pod, "Normal", "PUStart", "PU '%s' started successfully", puID)
			}

			// if this is a host network pod, we need to program the net_cls cgroup
			if pod.Spec.HostNetwork {
				netclsProgramCtx, netclsProgramCancel := context.WithTimeout(ctx, r.netclsProgramTimeout)
				defer netclsProgramCancel()
				if err := r.netclsProgrammer(netclsProgramCtx, pod, puRuntime); err != nil {
					if extractors.IsErrNetclsAlreadyProgrammed(err) {
						zap.L().Warn("net_cls cgroup has already been programmed previously", zap.String("puID", puID), zap.Error(err))
					} else if extractors.IsErrNoHostNetworkPod(err) {
						zap.L().Error("net_cls cgroup programmer told us that this is no host network pod. Aborting.", zap.String("puID", puID), zap.Error(err))
						return reconcile.Result{}, nil
					} else {
						zap.L().Error("failed to program net_cls cgroup of pod", zap.String("puID", puID), zap.Error(err))
						r.recorder.Eventf(pod, "Warning", "PUStart", "Host Network PU '%s' failed to program its net_cls cgroups: %s", puID, err.Error())
						return reconcile.Result{}, err
					}
				} else {
					zap.L().Debug("net_cls cgroup has been successfully programmed for trireme", zap.String("puID", puID))
					r.recorder.Eventf(pod, "Normal", "PUStart", "Host Network PU '%s' has successfully programmed its net_cls cgroups", puID)
				}
			}
		}
		return reconcile.Result{}, nil

	case corev1.PodSucceeded:
		fallthrough
	case corev1.PodFailed:
		zap.L().Debug("PodSucceeded / PodFailed", zap.String("puID", puID))
		err := r.handler.Policy.HandlePUEvent(
			handlePUCtx,
			puID,
			common.EventStop,
			policy.NewPURuntimeWithDefaults(),
		)
		if err != nil {
			if policy.IsErrPUNotFound(err) || policy.IsErrPUCreateFailed(err) {
				// not found means nothing needed stopping
				// just return
				zap.L().Debug("failed to handle stop event (IsErrPUNotFound==true)", zap.String("puID", puID), zap.Error(err))
				return reconcile.Result{}, nil
			}
			zap.L().Error("failed to handle stop event", zap.String("puID", puID), zap.Error(err))
			r.recorder.Eventf(pod, "Warning", "PUStop", "PU '%s' failed to stop: %s", puID, err.Error())
			return reconcile.Result{}, ErrHandlePUStopEventFailed

		}
		r.recorder.Eventf(pod, "Normal", "PUStop", "PU '%s' has been successfully stopped", puID)
		return reconcile.Result{}, nil

	case corev1.PodUnknown:
		zap.L().Error("pod is in unknown state", zap.String("puID", puID))

		// we don't need to retry, there is nothing *we* can do about it to fix this
		return reconcile.Result{}, nil
	default:
		zap.L().Error("unknown pod phase", zap.String("puID", puID), zap.String("podPhase", string(pod.Status.Phase)))

		// we don't need to retry, there is nothing *we* can do about it to fix this
		return reconcile.Result{}, nil
	}
}
