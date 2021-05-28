package graphql

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeMutationVariableKeys(t *testing.T) {
	cases := []struct {
		body             string
		computeKeys      map[string]interface{}
		expectedValues   map[string]interface{}
		expectedErrorMsg string
	}{
		{
			body:           `{"data": {"todo": {"id": "computed_id", "otherComputedValue": "computed_value_two"}}}`,
			computeKeys:    map[string]interface{}{"id_key": "todo.id", "other_computed_value": "todo.otherComputedValue"},
			expectedValues: map[string]interface{}{"id_key": "computed_id", "other_computed_value": "computed_value_two"},
		},
		{
			body:           `{"data": {"todos": [{"id": "computed_id"}, {"id": "second_id"}]}}`,
			computeKeys:    map[string]interface{}{"id_key": "todos[1].id"},
			expectedValues: map[string]interface{}{"id_key": "second_id"},
		},
		{
			body:             `{"data": {"todos": [{"id": "computed_id"}, {"id": "second_id"}]}}`,
			computeKeys:      map[string]interface{}{"id_key": "todos[3].id"},
			expectedErrorMsg: "provided index, 3, out of range for items in object todos with length of 2",
		},
		{
			body:             `{"data": {"todos": [{"id": "computed_id", "items": []}]}}`,
			computeKeys:      map[string]interface{}{"id_key": "todos[0].notreal[1]"},
			expectedErrorMsg: "mutation_key not found in nested object: notreal",
		},
		{
			body:             `{"data": {"todo": {"id": "computed_id"}}}`,
			computeKeys:      map[string]interface{}{"id_key": "todo[0].notreal"},
			expectedErrorMsg: "structure at key [todo] must be an array",
		},
		{
			body:             `{"data": {"todos": ["notanobject"]}}`,
			computeKeys:      map[string]interface{}{"id_key": "todos[0].id"},
			expectedErrorMsg: "malformed structure at provided index: 0",
		},
		{
			body:             `{"data": {"todos": [{"id": "computed_id"}, {"id": "second_id"}]}}`,
			computeKeys:      map[string]interface{}{"id_key": "todos.id"},
			expectedErrorMsg: "malformed structure at []interface {}{map[string]interface {}{\"id\":\"computed_id\"}, map[string]interface {}{\"id\":\"second_id\"}}",
		},
		{
			body:             `{"data": {"todo": {"id": "computed_id"}}}`,
			computeKeys:      map[string]interface{}{"id_key": "notreal.id"},
			expectedErrorMsg: "mutation_key not found",
		},
		{
			body:             `{"data": {"todo": {"id": "computed_id"}}}`,
			computeKeys:      map[string]interface{}{"id_key": ""},
			expectedErrorMsg: "no path provided for mutation key",
		},
	}

	for i, c := range cases {
		var robj = make(map[string]interface{})
		err := json.Unmarshal([]byte(c.body), &robj)
		if err != nil {
			t.Fatalf("Unable to unmarshal json response body %v", err)
		}

		m, err := computeMutationVariableKeys(c.computeKeys, robj)
		if c.expectedErrorMsg != "" {
			assert.Error(t, err, fmt.Sprintf("test case: %d", i))
			assert.EqualError(t, err, c.expectedErrorMsg, fmt.Sprintf("test case: %d", i))
		} else {
			for k, v := range c.expectedValues {
				assert.Equal(t, m[k], v, fmt.Sprintf("test case: %d", i))
			}
		}
	}
}
