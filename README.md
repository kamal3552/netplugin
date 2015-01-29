## Netplugin

Generic network plugin (experimental) is designed to handle networking use cases in clustered multi-host systems, beyond what is offered by the docker default plugin. Specifically it is designed to handle:
- Instantiating policies/acl/qos associated with containers
- Multi-tenant environment where disjoint networks are offered to containers on the same host
- SDN applications
- Interoperability with non container environment
- Multicast or multi-destination dependent applications
- Integration with existing IPAM tools for migrating customers

The overall design is _not_ assumed to be complete, because of ongoing work in the docker community with regards to the suitable APIs to interface with network extensions like this. Regardless, flexibility in the design has been taken into consideration to allow using a different state driver for key-value synchronization, or a different flavor of a soft-switch i.e. linux-bridge, macvlan, or openvswitch

The ability to specify the intent succintly is the primary goal of the design and thus some of the specificed user interface will change, and in some cases functionality will be enhanced to accomodate the same. Design details and future work to be captured in a separate document.


###Building and Testing

`vagrant up`

Note: make sure virtualbox is installed

`vagrant ssh default`

`sudo -s`

`source /etc/profile.d/envvar.sh`

`cd $GOSRC/github.com/contiv/netplugin`

`make unit-test`

###Trying it out 

The netplugin produces two binaries, a netplugin daemon and a netdcli tool to interact with it.

####A quick example

1. Start netplugin
`netplugin`
2. Specificy global pool of resources (can be executed on any host) Details later on using vxlans instead of vlans, 
`netdcli -oper=create -construct=global -subnet=11.1.0.0/16 -vlans "11-20" global`
3. Create a `orange` network
`netdcli -oper create -construct network orange`
4. Attach containers `myContainer1` and `myContainer2` to the `orange` network using
`netdcli -oper create -construct endpoint -net-id=orange -container-id=myContainer1 orange-endpoint1` and `netdcli -oper create -construct endpoint -net-id=orange -container-id=myContainer2 orange-endpoint2` respectively
5. `myContainer1` and `myContainer2` are interconnected
6. The configuration remains persistent, i.e. myContainer1 and myContainer2 can go and come back, the configuration is restored

There are many variations to the above configuration, like creating multiple networks,
across multiple hosts, use of vlans, use of vxlan, custom overrides for 
ip/subnet/vlan/vxlan allocation on per network/endpoint basis

####Auto-allocaiton of IP addresses
The plugin can automatically manage the IP address pools and assign an appropriate IP address based on the subnet that was associated with the network. However this doesn't take away the flexibility to keep a specific IP address of a container, which can always be specified as shown earlier. To automatically allocate the IP address, just avoid specifying the IP address during endpoint creation, for example in the previous example:
`netdcli -oper create -construct endpoint -net-id tenant1-net1 -container-id myContainer2 tenant1-net1-ep2`

With this, associating containers with networks will ensure a unique IP address is assigned to the container

While auto-allocation is allowed, per endpoint override to use a specific IP address 
is avialable.

####Auto-allocaiton of Subnets
The plugin can automatically manage the assignment of IP subnets to be used for various networks that are created. This would require configuring the global pool of ip-subnets to pick the subnet allocation from. The implementation will allow distributed atomicity to avoid conflicts

While auto-allocation is allowed, per network override to use a specific subnet 
is avialable.

####Auto-allocation of vlans and vxlan ids
Allocation of vlan-ids is specifically useful to allow interacting containers with 
non containerized applications. In many cases the default deployment choice of 
vlan/vxlan can be specified once as part of global configuration along with the 
allowed range (to avoid possiblee conflict), etc.

Auto allocation of vlan-ids and vxlan-id will be done if the network is not specified with the vlan/vxlan id, and a global pool is available.

While auto-allocation is allowed, per network override to use a specific vlan or vxlan-id is avialable to handle specific cases

####Bring up the netplugin daemon

`vagrant ssh default`

`sudo -s`

`source /etc/profile.d/envvar.sh`

`cd $GOSRC/github.com/contiv/netplugin`

`make build`

Ensure that $GOBIN is included in $PATH, then start the daemon as:

`netplugin`

####Create a network

Acquire another terminal to execute netdcli commands to ensure that the logs from netplugin does not mix with netdcli output. 

First we start with defining a network (could be json input), for now let's use cli and specify the tag (default tag type is 'vlan') to use and subnet mask of '/24'. Let's call it tenant1-net1

`netdcli -oper create -construct network -tag 12 -subnet 0/24 tenant1-net1`

The oepration state of network can be read using 

`netdcli -oper get -construct network tenant1-net1`

####Create an endpoint (an endpoint is an interface to be associated with container)

`netdcli -oper create -construct endpoint -net-id tenant1-net1 -ip-address="11.1.1.1" tenant1-net1-ep1`

Reading back the endpoint operation state can be done using

`netdcli -oper get -construct network tenant1-net1`

`netdcli -oper get -construct endpoint tenant1-net1-ep1`

####Associate an endpoint to a running container (this can also be done during endpoint create)

In order to associate a container to, create a container first
`docker run -it --name=myContainer1 --hostname=myContainer1 ubuntu /bin/bash`

Then, attach the container to the endpoint. Even if the association was done earlier it would work exactly the same
`netdcli -oper attach -construct endpoint -container-id myContainer1 tenant1-net1-ep1`

To associate the container during endpoint creation just pass `cont-id` parameter

To detach an endpoint from a container use detach command
`netdcli -oper detach -construct endpoint -container-id myContainer1 tenant1-net1-ep1`

####Ensure that all is operational

Ensure that a port got added to the ovs bridge named vlanBr

`sudo ovsctl show`

- verify that a linux device is also created for the port added above

`ip link show`

####Add more containers and make sure they can talk to each other
Let's start another container
`docker run -it --name=myContainer2 --hostname=myContainer2 ubuntu /bin/bash`

Add the newly added container to the same tenant's network and attach it to a container. This time, instead of using ep create and attach ep to container, let's specify all parameters during ep creation
`netdcli -oper create -construct endpoint -net-id tenant1-net1 -ip-address="11.1.1.2" -container-id myContainer2 tenant1-net1-ep2`

At this point both containers would have been configured with IP address in a dedicated network called 'tenant1-net1' with an IP address allocated from the subnet/mask associated with the network. Therefore, if a ping test is done from either myContainer1 or myContainer2, it would succeed. IP address can overlap in various networks as long as outbound rules are non overlapping.

####Delete the endpoint

`netdcli -oper delete -construct endpoint tenant1-net1-ep1`

Read the network and endpoint state to verify that they are updated

### How to Contribute
We welcome patches and contributions, please hit the github page to open an issue or to submit patches send pull rquests. 
Happy hacking!
