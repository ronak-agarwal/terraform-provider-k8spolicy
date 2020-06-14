package k8spolicy

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceConstraint() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateConstraint,
		Read:   resourceReadConstraint,
		Delete: resourceDeleteConstraint,
		Update: resourceUpdateConstraint,

		Schema: map[string]*schema.Schema{
			"constraint_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"constraint_crd_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"applyon_apigroups": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"applyon_kinds": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCreateConstraint(d *schema.ResourceData, m interface{}) error {

	return nil
}
func resourceReadConstraint(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceDeleteConstraint(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceUpdateConstraint(d *schema.ResourceData, m interface{}) error {

	return nil
}
