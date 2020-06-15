package k8spolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
)

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
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"applyon_kinds": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"parameters_values": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCreateConstraint(d *schema.ResourceData, m interface{}) error {
	constraintGVR := k8sschema.GroupVersionResource{
		Group:    "constraints.gatekeeper.sh",
		Version:  "v1beta1",
		Resource: strings.ToLower(d.Get("constraint_crd_name").(string)),
	}

	client := m.(*Config).Client
	log.Printf("Creating NewConstraint")
	var constraint *unstructured.Unstructured

	var applyonAPIGroups = []string{""}
	if v, ok := d.GetOk("applyon_apigroups"); ok && len(v.([]interface{})) > 0 {
		applyonAPIGroups = expandStringList(v.([]interface{}))
	}

	var applyonKinds = []string{""}
	if v, ok := d.GetOk("applyon_kinds"); ok && len(v.([]interface{})) > 0 {
		applyonKinds = expandStringList(v.([]interface{}))
	}

	if d.Get("parameters_values").(string) == "" {
		NewConstraintWithoutParams(d.Get("constraint_name").(string), d.Get("constraint_crd_name").(string), constraintGVR.Group, applyonAPIGroups, applyonKinds)
	} else {
		var params map[string]interface{}
		json.Unmarshal([]byte(d.Get("parameters_values").(string)), &params)
		constraint = NewConstraintWithParams(d.Get("constraint_name").(string), d.Get("constraint_crd_name").(string), constraintGVR.Group, applyonAPIGroups, applyonKinds, params)
	}
	result, err := client.Resource(constraintGVR).Create(context.TODO(), constraint, metav1.CreateOptions{})
	errExit(fmt.Sprintf("Failed to create NewConstraint %#v", constraint), err)
	log.Printf("Created NewConstraint %s", result)

	return resourceReadConstraint(d, m)
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

// NewConstraintWithParams ...
func NewConstraintWithParams(name, crd, gvrgroup string, applyOnApigroups []string, applyOnKinds []string, parms map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       crd,
			"apiVersion": gvrgroup + "/v1beta1",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"match": map[string]interface{}{
					"kinds": []map[string]interface{}{
						{
							"apiGroups": applyOnApigroups,
							"kinds":     applyOnKinds,
						},
					},
				},
				"parameters": parms,
			},
		},
	}
}

// NewConstraintWithoutParams ...
func NewConstraintWithoutParams(name, crd, gvrgroup string, applyOnApigroups []string, applyOnKinds []string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       crd,
			"apiVersion": gvrgroup + "/v1beta1",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"match": map[string]interface{}{
					"kinds": []map[string]interface{}{
						{
							"apiGroups": applyOnApigroups,
							"kinds":     applyOnKinds,
						},
					},
				},
			},
		},
	}
}
