// +build !windows

package iptablesctrl

import (
	"fmt"

	"go.uber.org/zap"
)

// addContainerChain adds a chain for the specific container and redirects traffic there
// This simplifies significantly the management and makes the iptable rules more readable
// All rules related to a container are contained within the dedicated chain
func (i *iptables) addContainerChain(cfg *ACLInfo) error {
	appChain := cfg.AppChain
	netChain := cfg.NetChain
	if err := i.impl.NewChain(appPacketIPTableContext, appChain); err != nil {
		return fmt.Errorf("unable to add chain %s of context %s: %s", appChain, appPacketIPTableContext, err)
	}

	// if err := i.impl.NewChain(appProxyIPTableContext, appChain); err != nil {
	// 	return fmt.Errorf("unable to add chain %s of context %s: %s", appChain, appPacketIPTableContext, err)
	// }

	if err := i.impl.NewChain(netPacketIPTableContext, netChain); err != nil {
		return fmt.Errorf("unable to add netchain %s of context %s: %s", netChain, netPacketIPTableContext, err)
	}

	return nil
}

// removeGlobalHooksPre is called before we jump into template driven rules.This is best effort
// no errors if these things fail.
func (i *iptables) removeGlobalHooksPre() {
	rules := [][]string{
		{
			"nat",
			"PREROUTING",
			"-p", "tcp",
			"-m", "addrtype",
			"--dst-type", "LOCAL",
			"-m", "set", "!", "--match-set", "TRI-Excluded", "src",
			"-j", "TRI-Redir-Net",
		},
		{
			"nat",
			"OUTPUT",
			"-m", "set", "!", "--match-set", "TRI-Excluded", "dst",
			"-j", "TRI-Redir-App",
		},
	}

	for _, rule := range rules {
		if err := i.impl.Delete(rule[0], rule[1], rule[2:]...); err != nil {
			zap.L().Debug("Error while delete rules", zap.Strings("rule", rule))
		}
	}

}

func transformACLRules(aclRules [][]string, cfg *ACLInfo, rulesBucket *rulesInfo, isAppAcls bool) [][]string {
	// pass through on linux
	return aclRules
}