# terraform-provider-k8spolicy


resource "k8spolicy_constraint_template" "my-server" {

  constraint_name = ""
  constraint_crd_name = ""
  parameters = << PARAMETERS

  PARAMETERS

  template_defination = << DEFINATION

  DEFINATION
}
