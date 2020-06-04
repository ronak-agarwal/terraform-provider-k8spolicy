package k8spolicy

import (
	k8sunstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

func parseJSON(json string) (ur *k8sunstructured.Unstructured, err error) {
	body := []byte(json)
	u, err := k8sruntime.Decode(k8sunstructured.UnstructuredJSONScheme, body)
	if err != nil {
		return ur, err
	}

	ur = u.(*k8sunstructured.Unstructured)

	return ur, nil
}
