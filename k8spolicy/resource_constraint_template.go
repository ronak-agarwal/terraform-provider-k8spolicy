package k8spolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	constraintTemplateGVR = k8sschema.GroupVersionResource{
		Group:    "templates.gatekeeper.sh",
		Version:  "v1beta1",
		Resource: "constrainttemplates",
	}
)

func resourceConstraintTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateConstraintTemplate,
		Read:   resourceReadConstraintTemplate,
		Delete: resourceDeleteConstraintTemplate,
		Update: resourceUpdateConstraintTemplate,

		Schema: map[string]*schema.Schema{
			"constraint_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"constraint_crd_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"rego_defination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// resourceCreateConstraintTemplate ...
func resourceCreateConstraintTemplate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).Client
	log.Printf("Creating NewConstraintTemplate")
	var constraintTemplate *unstructured.Unstructured

	if d.Get("parameters").(string) == "" {
		constraintTemplate = NewConstraintTemplateWithoutParams(d.Get("constraint_template_name").(string), d.Get("constraint_crd_name").(string), d.Get("rego_defination").(string))
	} else {
		var params map[string]interface{}
		json.Unmarshal([]byte(d.Get("parameters").(string)), &params)
		constraintTemplate = NewConstraintTemplateWithParams(d.Get("constraint_template_name").(string), d.Get("constraint_crd_name").(string), d.Get("rego_defination").(string), params)
	}
	result, err := client.Resource(constraintTemplateGVR).Create(context.TODO(), constraintTemplate, metav1.CreateOptions{})
	errExit(fmt.Sprintf("Failed to create NewConstraintTemplate %#v", constraintTemplate), err)
	log.Printf("Created NewConstraintTemplate %s", result)

	return resourceReadConstraintTemplate(d, m)
}

func resourceReadConstraintTemplate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceDeleteConstraintTemplate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceUpdateConstraintTemplate(d *schema.ResourceData, m interface{}) error {

	return nil
}

// NewConstraintTemplateWithoutParams ...
// New Class to create ConstraintTemplate
func NewConstraintTemplateWithoutParams(name, crd, rego string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "ConstraintTemplate",
			"apiVersion": constraintTemplateGVR.Group + "/v1beta1",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"crd": map[string]interface{}{
					"spec": map[string]interface{}{
						"names": map[string]interface{}{
							"kind": crd,
						},
					},
				},
				"targets": []map[string]interface{}{
					{
						"target": "admission.k8s.gatekeeper.sh",
						"rego":   rego,
					},
				},
			},
		},
	}
}

// NewConstraintTemplateWithParams ...
// New Class to create ConstraintTemplate
func NewConstraintTemplateWithParams(name, crd, rego string, params map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "ConstraintTemplate",
			"apiVersion": constraintTemplateGVR.Group + "/v1beta1",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"crd": map[string]interface{}{
					"spec": map[string]interface{}{
						"names": map[string]interface{}{
							"kind": crd,
						},
						"validation": map[string]interface{}{
							"openAPIV3Schema": map[string]interface{}{
								"properties": params,
							},
						},
					},
				},
				"targets": []map[string]interface{}{
					{
						"target": "admission.k8s.gatekeeper.sh",
						"rego":   rego,
					},
				},
			},
		},
	}
}

func errExit(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %#v", msg, err)
	}
}
