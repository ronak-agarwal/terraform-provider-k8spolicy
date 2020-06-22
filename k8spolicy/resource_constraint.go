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

func resourceConstraint() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateConstraint,
		Read:   resourceReadConstraint,
		Exists: resourceExistsConstraint,
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
		constraint = NewConstraintWithoutParams(d.Get("constraint_name").(string), d.Get("constraint_crd_name").(string), constraintGVR.Group, applyonAPIGroups, applyonKinds)
	} else {
		var params map[string]interface{}
		json.Unmarshal([]byte(d.Get("parameters_values").(string)), &params)
		constraint = NewConstraintWithParams(d.Get("constraint_name").(string), d.Get("constraint_crd_name").(string), constraintGVR.Group, applyonAPIGroups, applyonKinds, params)
	}
	result, err := client.Resource(constraintGVR).Create(context.TODO(), constraint, metav1.CreateOptions{})
	errExit(fmt.Sprintf("Failed to create NewConstraint %#v", constraint), err)
	log.Printf("Created NewConstraint %s", result)

	// Capture the UID at time of creation
	id := string(result.GetUID())
	d.SetId(id)
	// Add wait time of creation
	time.Sleep(10 * time.Second)

	return resourceReadConstraint(d, m)
}

func resourceReadConstraint(d *schema.ResourceData, m interface{}) error {
	constraintGVR := k8sschema.GroupVersionResource{
		Group:    "constraints.gatekeeper.sh",
		Version:  "v1beta1",
		Resource: strings.ToLower(d.Get("constraint_crd_name").(string)),
	}
	client := m.(*Config).Client
	result, err := client.Resource(constraintGVR).Get(context.TODO(), d.Get("constraint_name").(string), metav1.GetOptions{})
	errExit(fmt.Sprintf("Failed to read Constraint %#v", d.Get("constraint_name").(string)), err)
	// Capture the UID at time of creation
	id := string(result.GetUID())
	d.SetId(id)
	return nil
}

func resourceExistsConstraint(d *schema.ResourceData, m interface{}) (bool, error) {
	constraintGVR := k8sschema.GroupVersionResource{
		Group:    "constraints.gatekeeper.sh",
		Version:  "v1beta1",
		Resource: strings.ToLower(d.Get("constraint_crd_name").(string)),
	}
	client := m.(*Config).Client
	_, err := client.Resource(constraintGVR).Get(context.TODO(), d.Get("constraint_name").(string), metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		errExit(fmt.Sprintf("Failed to read Constraint Exists %#v", d.Get("constraint_name").(string)), err)
	}

	return true, err
}

func resourceDeleteConstraint(d *schema.ResourceData, m interface{}) error {
	constraintGVR := k8sschema.GroupVersionResource{
		Group:    "constraints.gatekeeper.sh",
		Version:  "v1beta1",
		Resource: strings.ToLower(d.Get("constraint_crd_name").(string)),
	}
	client := m.(*Config).Client
	err := client.Resource(constraintGVR).Delete(context.TODO(), d.Get("constraint_name").(string), metav1.DeleteOptions{})
	errExit(fmt.Sprintf("Failed to delete Constraint %#v", d.Get("constraint_name").(string)), err)
	// Success remove it from state
	d.SetId("")
	return nil
}

func resourceUpdateConstraint(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("constraint_name") {
		log.Printf("Cannot Update Existing Constraint Name !!")
		oldV, _ := d.GetChange("constraint_name")
		//log.Printf("new ConstraintTemplate CRD %s", newV)
		d.Set("constraint_name", strings.ToLower(oldV.(string)))
		return resourceReadConstraint(d, m)
	}

	constraintGVR := k8sschema.GroupVersionResource{
		Group:    "constraints.gatekeeper.sh",
		Version:  "v1beta1",
		Resource: strings.ToLower(d.Get("constraint_crd_name").(string)),
	}

	client := m.(*Config).Client

	getObj, err := client.Resource(constraintGVR).Get(context.TODO(), d.Get("constraint_name").(string), metav1.GetOptions{})
	errExit(fmt.Sprintf("Failed to get Constraint %#v", getObj), err)

	log.Printf("Updating Constraint")
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
		constraint = NewConstraintWithoutParams(d.Get("constraint_name").(string), d.Get("constraint_crd_name").(string), constraintGVR.Group, applyonAPIGroups, applyonKinds)
	} else {
		var params map[string]interface{}
		json.Unmarshal([]byte(d.Get("parameters_values").(string)), &params)
		constraint = NewConstraintWithParams(d.Get("constraint_name").(string), d.Get("constraint_crd_name").(string), constraintGVR.Group, applyonAPIGroups, applyonKinds, params)
	}
	constraint.SetResourceVersion(getObj.GetResourceVersion())
	result, err := client.Resource(constraintGVR).Update(context.TODO(), constraint, metav1.UpdateOptions{})
	errExit(fmt.Sprintf("Failed to Update Constraint %#v", constraint), err)
	log.Printf("Updated Constraint %s", result)

	// Capture the UID at time of creation
	id := string(result.GetUID())
	d.SetId(id)
	// Add wait time of creation
	time.Sleep(10 * time.Second)

	return resourceReadConstraint(d, m)
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
