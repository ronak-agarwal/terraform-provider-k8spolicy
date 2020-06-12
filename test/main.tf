provider k8spolicy {
  load_config_file = true
}

resource "k8spolicy_constraint_template" "my-policy" {

  constraint_name = "ronaklimit"
  constraint_crd_name = "Ronaklimit"
  parameters = "${file("${path.module}/parameters.json")}"

  rego_defination = "${file("${path.module}/rego.yml")}"
}
