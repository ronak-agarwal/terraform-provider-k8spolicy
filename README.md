# terraform-provider-k8spolicy

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
