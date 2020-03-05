package counters

import "sync"

type Counters struct {
	counters []uint32

	sync.RWMutex
}

const (
	totalCounters = 150
)

// CounterTypes custom counter error type
type CounterTypes int

// WARNING: Append any new counters at the end of the list.
// DO NOT CHANGE EXISTING ORDER.
const (
	ErrUnknownError = iota
	ErrNonPUTraffic
	ErrNoConnFound
	ErrRejectPacket

	ErrMarkNotFound
	ErrPortNotFound
	ErrContextIDNotFound
	ErrInvalidProtocol

	ErrConnectionsProcessed
	ErrEncrConnectionsProcessed

	ErrUDPDropFin
	ErrUDPSynDroppedInvalidToken
	ErrUDPSynAckInvalidToken
	ErrUDPAckInvalidToken

	ErrUDPConnectionsProcessed
	ErrUDPContextIDNotFound
	ErrUDPDropQueueFull
	ErrUDPDropInNfQueue

	// 2
	// Processors
	ErrAppServicePreProcessorFailed
	ErrAppServicePostProcessorFailed
	ErrNetServicePreProcessorFailed
	ErrNetServicePostProcessorFailed

	// Syn
	ErrSynTokenFailed
	ErrSynDroppedInvalidToken
	ErrSynDroppedTCPOption
	ErrSynDroppedInvalidFormat
	ErrSynRejectPacket
	ErrSynDroppedExternalService
	ErrSynUnexpectedPacket
	ErrInvalidNetSynState
	ErrNetSynNotSeen

	// Synack
	ErrSynAckTokenFailed
	ErrOutOfOrderSynAck
	ErrInvalidSynAck
	ErrSynAckInvalidToken
	ErrSynAckMissingToken
	ErrSynAckNoTCPAuthOption
	ErrSynAckInvalidFormat
	ErrSynAckClaimsMisMatch
	ErrSynAckRejected
	ErrSynAckDroppedExternalService
	ErrInvalidNetSynAckState

	// Ack
	ErrAckTokenFailed
	ErrAckRejected
	ErrAckTCPNoTCPAuthOption
	ErrAckInvalidFormat
	ErrAckInvalidToken
	ErrAckInUnknownState
	ErrInvalidNetAckState

	// UDP Processors
	ErrUDPAppPreProcessingFailed // 50
	ErrUDPAppPostProcessingFailed
	ErrUDPNetPreProcessingFailed
	ErrUDPNetPostProcessingFailed

	// UDP Syn
	ErrUDPSynInvalidToken
	ErrUDPSynMissingClaims
	ErrUDPSynDroppedPolicy

	// UDP SynAck
	ErrUDPSynAckNoConnection
	ErrUDPSynAckPolicy

	// Dropped packets
	ErrDroppedTCPPackets
	ErrDroppedUDPPackets
	ErrDroppedICMPPackets
	ErrDroppedDNSPackets
	ErrDroppedDHCPPackets
	ErrDroppedNTPPackets

	// Connections expired
	ErrTCPConnectionsExpired
	ErrUDPConnectionsExpired

	//3

	ErrSynTokenEncodeFailed
	ErrSynTokenHashFailed
	ErrSynTokenSignFailed
	ErrSynSharedSecretMissing
	ErrSynInvalidSecret
	ErrSynInvalidTokenLength
	ErrSynMissingSignature
	ErrSynInvalidSignature
	ErrSynCompressedTagMismatch
	ErrSynDatapathVersionMismatch
	ErrSynTokenDecodeFailed
	ErrSynTokenExpired
	ErrSynSharedKeyHashFailed
	ErrSynPublicKeyFailed

	ErrSynAckTokenEncodeFailed
	ErrSynAckTokenHashFailed
	ErrSynAckTokenSignFailed
	ErrSynAckSharedSecretMissing
	ErrSynAckInvalidSecret
	ErrSynAckInvalidTokenLength
	ErrSynAckMissingSignature
	ErrSynAckInvalidSignature
	ErrSynAckCompressedTagMismatch
	ErrSynAckDatapathVersionMismatch
	ErrSynAckTokenDecodeFailed
	ErrSynAckTokenExpired
	ErrSynAckSharedKeyHashFailed
	ErrSynAckPublicKeyFailed

	ErrAckTokenEncodeFailed
	ErrAckTokenHashFailed
	ErrAckTokenSignFailed
	ErrAckSharedSecretMissing
	ErrAckInvalidSecret
	ErrAckInvalidTokenLength //50
	ErrAckMissingSignature
	ErrAckCompressedTagMismatch
	ErrAckDatapathVersionMismatch
	ErrAckTokenDecodeFailed
	ErrAckTokenExpired
	ErrAckSignatureMismatch

	// udp 3
	ErrUDPSynTokenFailed
	ErrUDPSynTokenEncodeFailed
	ErrUDPSynTokenHashFailed
	ErrUDPSynTokenSignFailed
	ErrUDPSynSharedSecretMissing
	ErrUDPSynInvalidSecret
	ErrUDPSynInvalidTokenLength
	ErrUDPSynMissingSignature
	ErrUDPSynInvalidSignature
	ErrUDPSynCompressedTagMismatch
	ErrUDPSynDatapathVersionMismatch
	ErrUDPSynTokenDecodeFailed
	ErrUDPSynTokenExpired
	ErrUDPSynSharedKeyHashFailed
	ErrUDPSynPublicKeyFailed

	ErrUDPSynAckTokenFailed
	ErrUDPSynAckTokenEncodeFailed
	ErrUDPSynAckTokenHashFailed
	ErrUDPSynAckTokenSignFailed
	ErrUDPSynAckSharedSecretMissing
	ErrUDPSynAckInvalidSecret
	ErrUDPSynAckInvalidTokenLength
	ErrUDPSynAckMissingSignature
	ErrUDPSynAckInvalidSignature
	ErrUDPSynAckCompressedTagMismatch
	ErrUDPSynAckDatapathVersionMismatch
	ErrUDPSynAckTokenDecodeFailed
	ErrUDPSynAckTokenExpired
	ErrUDPSynAckSharedKeyHashFailed
	ErrUDPSynAckPublicKeyFailed

	ErrUDPAckTokenFailed
	ErrUDPAckTokenEncodeFailed
	ErrUDPAckTokenHashFailed
	ErrUDPAckSharedSecretMissing
	ErrUDPAckInvalidSecret
	ErrUDPAckInvalidTokenLength
	ErrUDPAckMissingSignature
	ErrUDPAckCompressedTagMismatch
	ErrUDPAckDatapathVersionMismatch
	ErrUDPAckTokenDecodeFailed
	ErrUDPAckTokenExpired
	ErrUDPAckSignatureMismatch //48
)
