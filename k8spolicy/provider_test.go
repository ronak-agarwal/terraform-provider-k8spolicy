package k8spolicy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	upstream "github.com/terraform-providers/terraform-provider-kubernetes/kubernetes"
	"k8s.io/apimachinery/pkg/api/errors"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"k8spolicy":  testAccProvider,
		"kubernetes": upstream.Provider(),
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccCheckK8sPolicyDestroy(s *terraform.State) error {
	return testAccCheckK8sPolicyStatus(s, false)
}

func testAccCheckK8sPolicyExists(s *terraform.State) error {
	return testAccCheckK8sPolicyStatus(s, true)
}

func testAccCheckK8sPolicyStatus(s *terraform.State, shouldExist bool) error {
	conn, _ := testAccProvider.Meta().(KubeProvider)()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "k8spolicy" {
			continue
		}

		content, err := conn.RESTClient().Get().AbsPath(rs.Primary.ID).DoRaw()
		if (errors.IsNotFound(err) || errors.IsGone(err)) && shouldExist {
			return fmt.Errorf("Failed to find resource, likely a failure to create occured: %+v %v", err, string(content))
		}

	}

	return nil
}
