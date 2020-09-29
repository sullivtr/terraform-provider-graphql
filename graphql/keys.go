package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
)

func buildResourceKeyArgs(key string) []string {
	split := strings.Split(key, ".")
	frk := split[0]
	if frk != "data" {
		split = append(split, "_")
		copy(split[1:], split[0:])
		split[0] = "data"
	}
	return split
}

func computeMutationVariableKeys(keyMaps map[string]interface{}, responseObject map[string]interface{}) (map[string]string, error) {
	var mvks = make(map[string]string)
	for k, v := range keyMaps {
		resourceKeyArgs := buildResourceKeyArgs(v.(string))
		key, err := getResourceKey(responseObject, resourceKeyArgs...)
		if err != nil {
			return nil, err
		}
		mvks[k] = key.(string)
	}
	return mvks, nil
}

func getResourceKey(m map[string]interface{}, ks ...string) (val interface{}, err error) {
	var ok bool

	if len(ks) == 0 {
		return nil, fmt.Errorf("Query response object is empty")
	}
	if val, ok = m[ks[0]]; !ok {
		return nil, fmt.Errorf("mutation_key not found")
	} else if len(ks) == 1 {
		return val, nil
	} else if m, ok = val.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("malformed structure at %#v", val)
	} else {
		return getResourceKey(m, ks[1:]...)
	}
}

func hash(v []byte) int {
	queryResponseObj := make(map[string]interface{})
	_ = json.Unmarshal(v, &queryResponseObj)
	out, err := json.Marshal(queryResponseObj)
	if err != nil {
		panic(err)
	}
	return hashcode.String(string(out))
}
