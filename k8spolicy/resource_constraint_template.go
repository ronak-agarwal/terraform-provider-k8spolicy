package k8spolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
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
		Exists: resourceExistsConstraintTemplate,
		Delete: resourceDeleteConstraintTemplate,
		Update: resourceUpdateConstraintTemplate,

		Schema: map[string]*schema.Schema{
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
	templateName := strings.ToLower(d.Get("constraint_crd_name").(string))
	if d.Get("parameters").(string) == "" {
		constraintTemplate = NewConstraintTemplateWithoutParams(templateName, d.Get("constraint_crd_name").(string), d.Get("rego_defination").(string))
	} else {
		var params map[string]interface{}
		json.Unmarshal([]byte(d.Get("parameters").(string)), &params)
		constraintTemplate = NewConstraintTemplateWithParams(templateName, d.Get("constraint_crd_name").(string), d.Get("rego_defination").(string), params)
	}
	result, err := client.Resource(constraintTemplateGVR).Create(context.TODO(), constraintTemplate, metav1.CreateOptions{})
	errExit(fmt.Sprintf("Failed to create NewConstraintTemplate %#v", constraintTemplate), err)
	log.Printf("Created NewConstraintTemplate %s", result)

	// Capture the UID at time of creation
	id := string(result.GetUID())
	d.SetId(id)
	// Add wait time of creation
	time.Sleep(10 * time.Second)

	return resourceReadConstraintTemplate(d, m)
}

func resourceReadConstraintTemplate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).Client
	templateName := strings.ToLower(d.Get("constraint_crd_name").(string))
	result, err := client.Resource(constraintTemplateGVR).Get(context.TODO(), templateName, metav1.GetOptions{})
	errExit(fmt.Sprintf("Failed to read ConstraintTemplate %#v", templateName), err)
	// Capture the UID at time of creation
	id := string(result.GetUID())
	d.SetId(id)

	return nil
}

func resourceExistsConstraintTemplate(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*Config).Client
	templateName := strings.ToLower(d.Get("constraint_crd_name").(string))
	_, err := client.Resource(constraintTemplateGVR).Get(context.TODO(), templateName, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		errExit(fmt.Sprintf("Failed to read ConstraintTemplate Exists %#v", templateName), err)
	}

	return true, err
}

func resourceDeleteConstraintTemplate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).Client
	templateName := strings.ToLower(d.Get("constraint_crd_name").(string))
	err := client.Resource(constraintTemplateGVR).Delete(context.TODO(), templateName, metav1.DeleteOptions{})
	errExit(fmt.Sprintf("Failed to delete ConstraintTemplate %#v", templateName), err)
	// Success remove it from state
	d.SetId("")
	return nil
}

func resourceUpdateConstraintTemplate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("constraint_crd_name") {
		log.Printf("Cannot Update Existing ConstraintTemplate CRD !!")
		oldV, _ := d.GetChange("constraint_crd_name")
		//log.Printf("new ConstraintTemplate CRD %s", newV)
		d.Set("constraint_crd_name", strings.ToLower(oldV.(string)))
		return resourceReadConstraintTemplate(d, m)
	}
	client := m.(*Config).Client
	log.Printf("Updating NewConstraintTemplate")
	templateName := strings.ToLower(d.Get("constraint_crd_name").(string))
	getObj, err := client.Resource(constraintTemplateGVR).Get(context.TODO(), templateName, metav1.GetOptions{})
	errExit(fmt.Sprintf("Failed to get ConstraintTemplate %#v", getObj), err)
	var constraintTemplate *unstructured.Unstructured
	if d.HasChange("rego_defination") || d.HasChange("parameters") {
		if d.Get("parameters").(string) == "" {
			constraintTemplate = NewConstraintTemplateWithoutParams(templateName, d.Get("constraint_crd_name").(string), d.Get("rego_defination").(string))
		} else {
			var params map[string]interface{}
			json.Unmarshal([]byte(d.Get("parameters").(string)), &params)
			constraintTemplate = NewConstraintTemplateWithParams(templateName, d.Get("constraint_crd_name").(string), d.Get("rego_defination").(string), params)
		}
		constraintTemplate.SetResourceVersion(getObj.GetResourceVersion())
		result, err := client.Resource(constraintTemplateGVR).Update(context.TODO(), constraintTemplate, metav1.UpdateOptions{})
		errExit(fmt.Sprintf("Failed to update NewConstraintTemplate %#v", constraintTemplate), err)

		// Capture the UID at time of creation
		id := string(result.GetUID())
		d.SetId(id)
		// Add wait time of creation
		time.Sleep(10 * time.Second)
	}
	return resourceReadConstraintTemplate(d, m)
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
