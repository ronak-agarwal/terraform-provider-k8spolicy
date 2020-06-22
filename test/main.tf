provider k8spolicy {
  load_config_file = true
}

resource "k8spolicy_constraint_template" "my-policy" {

  constraint_crd_name = "Ronaklimit"
  parameters = "${file("${path.module}/parameters.json")}"
  rego_defination = "${file("${path.module}/rego.yml")}"
}

resource "k8spolicy_constraint" "my-constraint" {
  depends_on = [k8spolicy_constraint_template.my-policy]
  constraint_name = "ronak-constraint"
  constraint_crd_name = "Ronaklimit"
  applyon_apigroups = [""]
  applyon_kinds = ["Pod"]
  parameters_values = "${file("${path.module}/parameters_values.json")}"
}
