# terraform-provider-k8spolicy

STATUS - Testing of Alpha Release

## Prerequisite

1. You need to have Gatekeeper OPA v3 running on your kubernetes - https://github.com/open-policy-agent/gatekeeper

2. Identify list of Gatekeeper policies from this library which you need to apply on your cluster - https://github.com/open-policy-agent/gatekeeper/tree/master/library


## Usage

1. You can download the plugin binary from latest release - https://github.com/ronak-agarwal/terraform-provider-k8spolicy/releases/download/v1.0-alpha/terraform-provider-k8spolicy
2. Copy binary ~/.terraform.d/plugins/darwin_amd64/

#### (A) Create ConstraintTemplate
-- In Testing --

Sample Policy Template used for below example - https://github.com/open-policy-agent/gatekeeper/blob/master/library/pod-security-policy/apparmor/template.yaml

Below is main.tf sample for creating ConstraintTemplate

```hcl
resource "k8spolicy_constraint_template" "my-policy" {

  constraint_crd_name = "K8sPSPAppArmor"  //Template Name is always lower case of this value
  parameters = "${file("${path.module}/parameters.json")}"
  rego_defination = "${file("${path.module}/rego.yml")}"
}
```
There are two inputs -

a) Below is sample rego.yml for creating ConstraintTemplate (Required).

You can find this rego content from Gatekeeper policies library, below is sample template for AppArmor Rego

```hcl
package k8spspapparmor
violation[{"msg": msg, "details": {}}] {
    metadata := input.review.object.metadata
    container := input_containers[_]
    not input_apparmor_allowed(container, metadata)
    msg := sprintf("AppArmor profile is not allowed, pod: %v, container: %v. Allowed profiles: %v", [input.review.object.metadata.name, container.name, input.parameters.allowedProfiles])
}
input_apparmor_allowed(container, metadata) {
    metadata.annotations[key] == input.parameters.allowedProfiles[_]
    key == sprintf("container.apparmor.security.beta.kubernetes.io/%v", [container.name])
}
input_containers[c] {
    c := input.review.object.spec.containers[_]
}
input_containers[c] {
    c := input.review.object.spec.initContainers[_]
}
```

b) Below is sample params.json for creating ConstraintTemplate (Optional depends on Rego difination).

For creating such json you can copy the yaml content from any OPA Gatekeeper policy templates and paste in any online yaml to json convertor (Eg. https://onlineyamltools.com/convert-yaml-to-json).

You should pick yaml, anything between properties and targets [copy till targets and not targets keys].

```hcl
{
  "allowedProfiles": {
    "type": "array",
    "items": {
      "type": "string"
    }
  }
}
```


#### (B) Create Constraint

-- In Testing --


```hcl
resource "k8spolicy_constraint" "my-constraint" {
  depends_on = [k8spolicy_constraint_template.my-policy] //make sure constraint is created once its template CRD available
  constraint_name = "constraint-pod"
  constraint_crd_name = "K8sPSPAppArmor"
  applyon_apigroups = [""]
  applyon_kinds = ["Pod"]
  parameters_values = "${file("${path.module}/parameters_values.json")}"
}

```

Sample parameters_values.json picked from same apparmor constraint.yaml -
https://github.com/open-policy-agent/gatekeeper/blob/master/library/pod-security-policy/apparmor/constraint.yaml

```hcl
{
  "allowedProfiles": [
    "runtime/default"
  ]
}
```
