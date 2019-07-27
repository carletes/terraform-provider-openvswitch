# Terraform provider for Open vSwitch

This Terraform provider manages local Open vSwitch bridges and ports.


## Sample usage

From [examples/sample-bridge](./examples/sample-bridge/):

```
provider "openvswitch" {}

resource "openvswitch_bridge" "sample_bridge" {
  name = "testbr0"
}

resource "openvswitch_port" "sample_port" {
  count     = 2
  name      = "p${count.index}"
  bridge_id = openvswitch_bridge.sample_bridge.id
}
```


## Installation from source

Requirements:

* Go 1.11.x or later
* GNU Make

Clone this repo, and then do the fhe following:

```
$ make
$ cp ${GOPATH}/bin/terraform-provider-openvswitch ~/.terraform.d/plugins/
```
