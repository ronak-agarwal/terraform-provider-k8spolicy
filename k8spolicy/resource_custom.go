package k8spolicy

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
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

func resourceCreate(d *schema.ResourceData, m interface{}) error {

	f, err := os.Open("./constrainttemplate/template.yaml")
	if err != nil {
		log.Fatal(err)
	}
	decode := yaml.NewYAMLOrJSONDecoder(f, 4096)
	ext := runtime.RawExtension{}
	if err := decode.Decode(&ext); err != nil {
		if err == io.EOF {
			log.Fatal(err)
		}
		log.Fatal(err)
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = make(map[string]interface{})
	var blob interface{}
	if err := json.Unmarshal(ext.Raw, &blob); err != nil {
		log.Fatal(err)
	}

	unstruct.Object = blob.(map[string]interface{})

	constraintNameMap := map[string]string{}
	constraintNameMap["name"] = d.Get("constraint_name").(string)
	unstructured.SetNestedStringMap(unstruct.Object, constraintNameMap, "metadata")

	crdNameMap := map[string]string{}
	crdNameMap["kind"] = d.Get("constraint_crd_name").(string)
	unstructured.SetNestedStringMap(unstruct.Object, crdNameMap, "spec", "crd", "spec", "names")

	unstruct.Object["spec"].(map[string]interface{})["targets"].([]interface{})[0].(map[string]interface{})["rego"] = d.Get("template_defination").(string)

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
