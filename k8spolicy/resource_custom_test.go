package k8spolicy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceK8SPolicy_constraint_template(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders, Steps: []resource.TestStep{
			{
				Config: testAccResourceK8sPolicyConfigbasicInitial(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("k8spolicy_constraint_template.test", "constraint_name", "k8scontainerlimits"),
					resource.TestCheckResourceAttr("k8spolicy_constraint_template.test", "constraint_crd_name", "k8scontainerCRD"),
					resource.TestCheckResourceAttr("k8spolicy_constraint_template.test", "parameters", "params"),
					resource.TestCheckResourceAttr("k8spolicy_constraint_template.test", "template_defination", "defination"),
				),
			},
		},
	})
}

func testAccResourceK8sPolicyConfigbasicInitial() string {

	return ""
}
