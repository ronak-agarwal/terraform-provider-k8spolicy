# terraform-provider-k8spolicy

You need to have Gatekeeper OPA v3 running on your kubernetes - https://github.com/open-policy-agent/gatekeeper

Supported list of policy libraries - https://github.com/open-policy-agent/gatekeeper/tree/master/library

## Usage

```hcl
resource "k8spolicy_constraint_template" "my-policy" {

  constraint_name = ""
  constraint_crd_name = ""
  parameters = << PARAMETERS

  PARAMETERS

  template_defination = << DEFINATION

  DEFINATION
}
```
