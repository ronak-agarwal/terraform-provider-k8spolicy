package k8spolicy

import (
	"strings"
)

func expandStringList(configured []interface{}) []string {
	if configured[0] != nil && len(configured) == 1 && strings.Contains(configured[0].(string), ",") {
		return strings.Split(configured[0].(string), ",")
	}

	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		if v != nil {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
