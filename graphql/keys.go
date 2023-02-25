package graphql

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
)

func buildResourceKeyArgs(key string) []string {
	if key == "" {
		return nil
	}
	split := strings.Split(key, ".")
	frk := split[0]
	if frk != "data" {
		// Append placeholder for the "data" parent in GraphQL response
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
		return nil, fmt.Errorf("no path provided for mutation key")
	}

	// If path contains an array index, find the key from the items object at index
	if strings.Contains(ks[0], "[") {
		var items []interface{}

		// Parse out (array) key name and its index
		str := ks[0]
		str = strings.Replace(str, "[", ".", -1)
		str = strings.Trim(str, "]")
		split := strings.Split(str, ".")
		index, _ := strconv.ParseInt(split[1], 10, 32)

		// Ensure the object exists on m
		if val, ok = m[split[0]]; !ok {
			return nil, fmt.Errorf("mutation_key not found in nested object: %s", split[0])
		} else if items, ok = val.([]interface{}); !ok {
			// The object was not an array, and therefore cannot be indexed
			return nil, fmt.Errorf("structure at key [%s] must be an array", split[0])
		} else if index > int64(len(items))-1 {
			return nil, fmt.Errorf("provided index, %d, out of range for items in object %s with length of %d", index, split[0], len(items))
		}

		// Parse object at provided index
		var obj map[string]interface{}
		for i := range items {
			if index == int64(i) {
				if strValue, ok := items[i].(string); ok {
					return strValue, nil
				}
				if obj, ok = items[i].(map[string]interface{}); !ok {
					return nil, fmt.Errorf("malformed structure at provided index: %d", i)
				}
				break
			}
		}
		return getResourceKey(obj, ks[1:]...)
	}

	// The expected key is not nested inside of an array, so grab its value or continue walking the object tree
	if val, ok = m[ks[0]]; !ok { // Check that the item exists
		return nil, fmt.Errorf("mutation_key not found")
	} else if len(ks) == 1 {
		return val, nil // use value of root
	} else if m, ok = val.(map[string]interface{}); !ok { // If value is non an object, return error
		return nil, fmt.Errorf("malformed structure at %#v", val)
	} else {
		// parse value for the next property
		return getResourceKey(m, ks[1:]...)
	}
}

func mapQueryResponseInputKey(m interface{}, value, prev string, parentKeys []string) (key string, ok bool) {
	if mapObj, isMap := m.(map[string]interface{}); isMap {
		for k, v := range mapObj {
			var jsonV string
			// Check if v is a string, and if not, marshal it as a json string for comparison
			if str, isString := v.(string); !isString {
				bytes, err := json.Marshal(v)
				if err != nil {
					return
				}
				jsonV = string(bytes)
			} else {
				jsonV = str
			}

			if !strings.Contains(jsonV, value) {
				continue
			}

			// If jsonV and input value are the same, we have found our match.
			if strings.Compare(jsonV, value) == 0 || jsonV == value {
				if len(parentKeys) != 0 && !strings.Contains(parentKeys[len(parentKeys)-1], "[") {
					newPrev := parentKeys[len(parentKeys)-1]
					if newPrev != prev {
						parentKeys = parentKeys[:len(parentKeys)-1]
					}
				}

				parentKeys = append(parentKeys, k)
				key = strings.Join(parentKeys, ".")
				ok = true
				return
			} else if slice, isSlice := v.([]interface{}); isSlice { // If v is a slice, we need to loop over items in the slice
				// key, ok, parentKeys = handleSlice(slice, value, parentKeys, prev)

				for i, iv := range slice {
					// if the items in the slice are strings, we can simply compare its values
					if innerStr, valueIsString := iv.(string); valueIsString {
						if innerStr == value {
							if len(parentKeys) != 0 && !strings.Contains(parentKeys[len(parentKeys)-1], "[") {
								newPrev := parentKeys[len(parentKeys)-1]
								if newPrev != prev {
									parentKeys = parentKeys[:len(parentKeys)-1]
								}
							}
							parentKeys = append(parentKeys, fmt.Sprintf("%s[%d]", k, i))
							key = strings.Join(parentKeys, ".")
							ok = true
							return
						}
					} else if innerObj, isObject := iv.(map[string]interface{}); isObject { // If v is an object, traverse the object futher
						parentKeys = append(parentKeys, fmt.Sprintf("%s[%d]", k, i))
						key, ok = mapQueryResponseInputKey(innerObj, value, k, parentKeys)
						if !ok {
							parentKeys = parentKeys[:len(parentKeys)-1]
							continue
						}
						return
					} else {
						parentKeys = parentKeys[:len(parentKeys)-1]
						continue
					}
				}
			} else if nv, isMap := v.(map[string]interface{}); isMap {
				parentKeys = append(parentKeys, k)
				key, ok = mapQueryResponseInputKey(nv, value, k, parentKeys)
				if !ok {
					continue
				}
			}
		}
	}
	return
}

func hash(v []byte) int {
	queryResponseObj := make(map[string]interface{})
	_ = json.Unmarshal(v, &queryResponseObj)
	out, err := json.Marshal(queryResponseObj)
	if err != nil {
		panic(err)
	}
	return hashCodeString(string(out))
}

// hashCodeString hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func hashCodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
