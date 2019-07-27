package openvswitch

import "github.com/hashicorp/terraform/helper/schema"

func resourceBridge() *schema.Resource {
	return &schema.Resource{
		Create: resourceBridgeCreate,
		Read:   resourceBridgeRead,
		Update: resourceBridgeUpdate,
		Delete: resourceBridgeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceBridgeCreate(d *schema.ResourceData, m interface{}) error {
	return resourceBridgeRead(d, m)
}

func resourceBridgeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceBridgeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBridgeRead(d, m)
}

func resourceBridgeDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
