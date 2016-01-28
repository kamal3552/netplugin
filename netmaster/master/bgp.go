/***
Copyright 2014 Cisco Systems Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package master

import (
	log "github.com/Sirupsen/logrus"
	"github.com/contiv/netplugin/core"
	"github.com/contiv/netplugin/netmaster/intent"
	"github.com/contiv/netplugin/netmaster/mastercfg"
)

//AddBgpNeighbors adds to the etcd state
func AddBgpNeighbors(stateDriver core.StateDriver, bgpCfg *intent.ConfigBgp) error {

	log.Infof("Adding bgp neighbor {%v}", bgpCfg)
	bgpState := &mastercfg.CfgBgpState{}
	bgpState.Hostname = bgpCfg.Hostname
	bgpState.As = bgpCfg.As
	bgpState.Neighbor = bgpCfg.Neighbor
	bgpState.StateDriver = stateDriver
	bgpState.ID = bgpCfg.Hostname
	err := bgpState.Write()

	if err != nil {
		return err
	}
	return nil
}

//DeleteBgpNeighbors deletes from etcd state
func DeleteBgpNeighbors(stateDriver core.StateDriver, hostname string) error {
	log.Infof("Deleting bgp neighbor for {%v}", hostname)
	bgpState := &mastercfg.CfgBgpState{}
	bgpState.StateDriver = stateDriver
	err := bgpState.Read(hostname)
	if err != nil {
		log.Errorf("Error reading bgp config for hostname %s", hostname)
		return err
	}
	err = bgpState.Clear()
	if err != nil {
		log.Errorf("Error deleing Bgp config for hostname %s", hostname)
		return err
	}
	return nil

}
