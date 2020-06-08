package k8spolicy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/icza/dyno"
	yamlParser "gopkg.in/yaml.v2"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	meta_v1_unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
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

// Read YAML template and Update based on terraform inputs
func readTemplateAndUpdateYaml(d *schema.ResourceData) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(dir + "/k8spolicy/constrainttemplate/template.yaml")
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

	//Final YAML output
	yaml, err := json.Marshal(unstruct)
	if err != nil {
		log.Fatal(err)
	}

	return string(yaml)
}

func resourceCreate(d *schema.ResourceData, m interface{}) error {

	//Get updated YAML
	yaml := readTemplateAndUpdateYaml(d)

	client, rawObj, err := getRestClientFromYaml(string(yaml), m.(KubeProvider))
	if err != nil {
		return fmt.Errorf("failed to create kubernetes rest client for resource: %+v", err)
	}
	// Create the resource in Kubernetes
	response, err := client.Create(rawObj, meta_v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create resource in kubernetes: %+v", err)
	}

	d.SetId(response.GetSelfLink())
	d.Set("uid", response.GetUID())
	d.Set("resource_version", response.GetResourceVersion())

	comparisonString, err := compareMaps(rawObj.UnstructuredContent(), response.UnstructuredContent())
	if err != nil {
		return err
	}
	log.Printf("[COMPAREOUT_CREATE] %+v\n", comparisonString)

	return resourceRead(d, m)
}

func resourceRead(d *schema.ResourceData, m interface{}) error {

	//Get updated YAML
	yaml := readTemplateAndUpdateYaml(d)

	client, rawObj, err := getRestClientFromYaml(string(yaml), m.(KubeProvider))
	if err != nil {
		return fmt.Errorf("failed to create kubernetes rest client for resource: %+v", err)
	}

	// Get the resource from Kubernetes
	metaObjLive, err := client.Get(rawObj.GetName(), meta_v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get resource '%s' from kubernetes: %+v", metaObjLive.GetSelfLink(), err)
	}

	if metaObjLive.GetUID() == "" {
		return fmt.Errorf("Failed to parse item and get UUID: %+v", metaObjLive)
	}

	// Capture the UID and Resource_version from the cluster at the current time
	d.Set("live_uid", metaObjLive.GetUID())
	d.Set("live_resource_version", metaObjLive.GetResourceVersion())

	comparisonOutput, err := compareMaps(rawObj.UnstructuredContent(), metaObjLive.UnstructuredContent())
	if err != nil {
		return err
	}
	log.Printf("[COMPAREOUT_READ] %+v\n", comparisonOutput)
	return nil
}

func resourceDelete(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceUpdate(d *schema.ResourceData, m interface{}) error {

	return nil
}

func getRestClientFromYaml(yaml string, provider KubeProvider) (dynamic.ResourceInterface, *meta_v1_unstruct.Unstructured, error) {
	// To make things play nice we need the JSON representation of the object as
	// the `RawObj`
	// 1. UnMarshal YAML into map
	// 2. Marshal map into JSON
	// 3. UnMarshal JSON into the Unstructured type so we get some K8s checking
	// 4. Marshal back into JSON ... now we know it's likely to play nice with k8s
	rawYamlParsed := &map[string]interface{}{}
	err := yamlParser.Unmarshal([]byte(yaml), rawYamlParsed)
	if err != nil {
		return nil, nil, err
	}

	rawJSON, err := json.Marshal(dyno.ConvertMapI2MapS(*rawYamlParsed))
	if err != nil {
		return nil, nil, err
	}

	unstrut := meta_v1_unstruct.Unstructured{}
	err = unstrut.UnmarshalJSON(rawJSON)
	if err != nil {
		return nil, nil, err
	}

	unstructContent := unstrut.UnstructuredContent()
	log.Printf("[UNSTRUCT]: %+v\n", unstructContent)

	// Use the k8s Discovery service to find all valid APIs for this cluster
	clientSet, config := provider()
	discoveryClient := clientSet.Discovery()
	resources, err := discoveryClient.ServerResources()
	// There is a partial failure mode here where not all groups are returned `GroupDiscoveryFailedError`
	// we'll try and continue in this condition as it's likely something we don't need
	// and if it is the `checkAPIResourceIsPresent` check will fail and stop the process
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return nil, nil, err
	}

	// Validate that the APIVersion provided in the YAML is valid for this cluster
	apiResource, exists := checkAPIResourceIsPresent(resources, unstrut)
	if !exists {
		return nil, nil, fmt.Errorf("resource provided in yaml isn't valid for cluster, check the APIVersion and Kind fields are valid")
	}

	resource := k8sschema.GroupVersionResource{Group: apiResource.Group, Version: apiResource.Version, Resource: apiResource.Name}
	// For core services (ServiceAccount, Service etc) the group is incorrectly parsed.
	// "v1" should be empty group and "v1" for verion
	if resource.Group == "v1" && resource.Version == "" {
		resource.Group = ""
		resource.Version = "v1"
	}
	client := dynamic.NewForConfigOrDie(&config).Resource(resource)

	if apiResource.Namespaced {
		namespace := unstrut.GetNamespace()
		if namespace == "" {
			namespace = "default"
		}
		return client.Namespace(namespace), &unstrut, nil
	}

	return client, &unstrut, nil
}

// checkAPIResourceIsPresent Loops through a list of available APIResources and
// checks there is a resource for the APIVersion and Kind defined in the 'resource'
// if found it returns true and the APIResource which matched
func checkAPIResourceIsPresent(available []*meta_v1.APIResourceList, resource meta_v1_unstruct.Unstructured) (*meta_v1.APIResource, bool) {
	for _, rList := range available {
		if rList == nil {
			continue
		}
		group := rList.GroupVersion
		for _, r := range rList.APIResources {
			if group == resource.GroupVersionKind().GroupVersion().String() && r.Kind == resource.GetKind() {
				r.Group = rList.GroupVersion
				r.Kind = rList.Kind
				return &r, true
			}
		}
	}
	return nil, false
}
