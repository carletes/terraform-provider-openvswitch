# Terraform provider for Open vSwitch

This Terraform provider manages local Open vSwitch bridges and ports.


## Sample usage

From [examples/sample-bridge](./examples/sample-bridge/):

```
provider "openvswitch" {}

resource "openvswitch_bridge" "sample_bridge" {
  name = "testbr0
  // Optional Parameters
  // OpenFlow10, OpenFlow11, OpenFlow12, OpenFlow14, OpenFlow15
  ofversion = "OpenFlow13"
}

resource "openvswitch_port" "sample_port" {
  count     = 2
  name      = "p${count.index}"
  ofversion = "OpenFlow13"
  bridge_id = openvswitch_bridge.sample_bridge.
  // Optional Field
  action	= "up"
}
```

## Important notes
- The ip, ovs-vsctl, ovs-ofctl commands all require sudo or root access
- Error handling is currently broken

## Installation from source

Requirements:

* Go 1.15.x or later
* GNU Make
* Terraform v0.12.* (This doesn't work on v0.13 yet)

Clone this repo, and then do the fhe following:

```
$ make
$ cp ${GOPATH}/bin/terraform-provider-openvswitch ~/.terraform.d/plugins/
```
