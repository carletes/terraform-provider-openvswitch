provider "openvswitch" {}

resource "openvswitch_bridge" "sample_bridge" {
  name = "testbr0"
}

resource "openvswitch_port" "sample_port" {
  count     = 2
  name      = "p${count.index}"
  bridge_id = openvswitch_bridge.sample_bridge.id
}
