package openvswitch

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a schema.Provider for OpenVSwitch.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{},

		ResourcesMap: map[string]*schema.Resource{
			"openvswitch_bridge": resourceBridge(),
		},

		DataSourcesMap: map[string]*schema.Resource{},
	}
}
