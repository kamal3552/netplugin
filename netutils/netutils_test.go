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

package netutils

import (
    "testing"
)

type testSubnetInfo struct {
    subnetIp    string
    subnetLen   uint
    hostId      uint
    hostIp      string
}

var testSubnets = []testSubnetInfo {
        {subnetIp:"11.2.1.0", subnetLen:24, hostId:5, hostIp:"11.2.1.5" },
        {subnetIp:"10.123.16.0", subnetLen:22, hostId:513, hostIp:"10.123.18.1" },
        {subnetIp:"172.12.0.0", subnetLen:16, hostId:261, hostIp:"172.12.1.5" },
    }

func TestGetSubnetIp(t *testing.T) {
    for _, te := range testSubnets {
        hostIp, err := GetSubnetIp(te.subnetIp, te.subnetLen, te.hostId)
        if err != nil {
            t.Fatalf("error getting host ip from subnet %s/%d for hostid %d - err '%s'",
                te.subnetIp, te.subnetLen, te.hostId, err)
        }
        if hostIp != te.hostIp {
            t.Fatalf("obtained ip %s doesn't match expected ip %s for subnet %s/%d\n",
                hostIp, te.hostIp, te.subnetIp, te.subnetLen)
        }
    }
}

var testInvalidSubnets = []testSubnetInfo {
        {subnetIp:"11.2.1.0", subnetLen:32, hostId:5, hostIp:"11.2.1.5" },
        {subnetIp:"10.123.16.0", subnetLen:22, hostId:1025, hostIp:"10.123.18.1" },
        {subnetIp:"172.12.0.0", subnetLen:4, hostId:261, hostIp:"172.12.1.5" },
    }

func TestInvalidGetSubnetIp(t *testing.T) {
    for _, te := range testInvalidSubnets {
        _, err := GetSubnetIp(te.subnetIp, te.subnetLen, te.hostId)
        if err == nil {
            t.Fatalf("Expecting error on invalid config subnet %s/%d for hostid %d",
                te.subnetIp, te.subnetLen, te.hostId)
        }
    }
}

func TestGetIpNumber(t *testing.T) {
    for _, te := range testSubnets {
        hostId, err := GetIpNumber(te.subnetIp, te.subnetLen, te.hostIp)
        if err != nil {
            t.Fatalf("error getting host ip from subnet %s/%d for hostid %d ",
                te.subnetIp, te.subnetLen, te.hostId)
        }
        if hostId != te.hostId {
            t.Fatalf("obtained ip %d doesn't match with expected ip %d \n",
                hostId, te.hostId)
        }
    }
}
