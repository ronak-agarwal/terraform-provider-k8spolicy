package k8spolicy

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCustom() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreate,
		Read:   resourceRead,
		Delete: resourceDelete,

		Schema: map[string]*schema.Schema{
			"constraint_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"constraint_crd_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"template_defination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCreate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceDelete(d *schema.ResourceData, m interface{}) error {

	return nil
}
