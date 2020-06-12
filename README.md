# terraform-provider-k8spolicy

STATUS - In Development

You need to have Gatekeeper OPA v3 running on your kubernetes - https://github.com/open-policy-agent/gatekeeper

Supported list of policy libraries - https://github.com/open-policy-agent/gatekeeper/tree/master/library

## Usage

#### Create ConstraintTemplate
Sample Policy Template used for below example - https://github.com/open-policy-agent/gatekeeper/blob/master/library/pod-security-policy/apparmor/template.yaml

Below is main.tf sample for creating ConstraintTemplate

```hcl
resource "k8spolicy_constraint_template" "my-policy" {

  constraint_name = "k8spspapparmor"
  constraint_crd_name = "K8sPSPAppArmor"
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


#### Create Constraint

TODO
