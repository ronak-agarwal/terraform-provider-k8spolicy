package k8spolicy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceK8SPolicy_constraint_template(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders, Steps: []resource.TestStep{
			{
				Config: testAccResourceK8sPolicyConfigbasicInitial(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("k8spolicy_constraint_template.test", "manifest", "k8scontainerlimits"),
				),
			},
		},
	})
}

func testAccResourceK8sPolicyConfigbasicInitial() string {
	return fmt.Sprintf(`
  resource "k8spolicy_constraint_template" "test" {
  	manifest = {
  "apiVersion": "templates.gatekeeper.sh/v1beta1",
  "kind": "ConstraintTemplate",
  "metadata": {
    "name": "k8scontainerlimits"
  },
  "spec": {
    "crd": {
      "spec": {
        "names": {
          "kind": "K8sContainerLimits"
        },
        "validation": {
          "openAPIV3Schema": {
            "properties": {
              "cpu": {
                "type": "string"
              },
              "memory": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "targets": [
      {
        "target": "admission.k8s.gatekeeper.sh",
        "rego": "package k8scontainerlimits\nmissing(obj, field) = true {\n  not obj[field]\n}\n"
      }
    ]
  }
}
  }
  `)
}
