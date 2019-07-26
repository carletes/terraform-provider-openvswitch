package main

import (
	"github.com/carletes/terraform-provider-openvswitch/openvswitch"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openvswitch.Provider})
}
