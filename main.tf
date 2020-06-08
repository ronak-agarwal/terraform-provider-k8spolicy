provider k8spolicy {
  load_config_file = true
}

resource "k8spolicy_constraint_template" "my-policy" {

  constraint_name = "RonakLimit"
  constraint_crd_name = "RonakCRD"
  parameters = "paramas"

  template_defination = "Test Ronak"
}
