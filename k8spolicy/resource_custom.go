package k8spolicy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCustom() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreate,
		Read:   resourceRead,
		Delete: resourceDelete,
		Update: resourceUpdate,

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

func policytemplate() string {
	return fmt.Sprintf(`
		{
		  "apiVersion": "templates.gatekeeper.sh/v1beta1",
		  "kind": "ConstraintTemplate",
		  "metadata": {
		    "name": "XX"
		  },
		  "spec": {
		    "crd": {
		      "spec": {
		        "names": {
		          "kind": "XX"
		        },
		        "validation": {
		          "openAPIV3Schema": {
		            "properties": {}
		          }
		        }
		      }
		    },
		    "targets": [
		      {
		        "target": "admission.k8s.gatekeeper.sh",
		        "rego": "XX"
		      }
		    ]
		  }
		}
  `)
}

func resourceCreate(d *schema.ResourceData, m interface{}) error {
	srcJSON := policytemplate()
	u, err := parseJSON(srcJSON)
	if err != nil {
		return fmt.Errorf("ResourceCreate: %s", err)
	}
	u.SetName("test")
	u.SetKind("test1")
	fmt.Println(u)
	return nil
}

func resourceRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceDelete(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}
