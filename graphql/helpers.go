package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
)

func buildResourceKeyArgs(slice []interface{}) [][]string {
	var ss = [][]string{}

	for _, v := range slice {
		split := strings.Split(v.(string), ".")
		frk := split[0]
		if frk != "data" {
			split = append(split, "_")
			copy(split[1:], split[0:])
			split[0] = "data"
		}
		ss = append(ss, split)
	}

	return ss
}

func computeMutationVariableKeys(keyMaps [][]string, responseObject map[string]interface{}) (map[string]string, error) {
	var mvks = make(map[string]string)
	for _, v := range keyMaps {
		k, v, err := getResourceKey(responseObject, v...)
		if err != nil {
			return nil, err
		}
		mvks[k] = v.(string)
	}
	return mvks, nil
}

func getResourceKey(m map[string]interface{}, ks ...string) (key string, val interface{}, err error) {
	var ok bool

	if len(ks) == 0 {
		return "", nil, fmt.Errorf("Query response object is empty")
	}
	if val, ok = m[ks[0]]; !ok {
		return "", nil, fmt.Errorf("query_response_key not found")
	} else if len(ks) == 1 {
		return ks[0], val, nil
	} else if m, ok = val.(map[string]interface{}); !ok {
		return "", nil, fmt.Errorf("malformed structure at %#v", val)
	} else {
		return getResourceKey(m, ks[1:]...)
	}
}

func hashString(v []byte) int {
	queryResponseObj := make(map[string]interface{})
	_ = json.Unmarshal(v, &queryResponseObj)
	out, err := json.Marshal(queryResponseObj)
	if err != nil {
		panic(err)
	}
	return hashcode.String(string(out))
}
