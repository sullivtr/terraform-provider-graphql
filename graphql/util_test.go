package graphql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeMutationVariableKeys(t *testing.T) {
	body := `{"data": {"todo": {"id": "computed_id", "otherComputedValue": "computed_value_two"}}}`

	var robj = make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &robj)
	if err != nil {
		t.Fatalf("Unable to unmarshal json response body %v", err)
	}

	compute_keys := make(map[string]interface{})
	compute_keys["id_key"] = "todo.id"
	compute_keys["other_computed_value"] = "todo.otherComputedValue"

	m, err := computeMutationVariableKeys(compute_keys, robj)
	if err != nil {
		t.Fatalf("Unable to compute mutation keys from response object %v", err)
	}

	idKey := m["id_key"]
	otherComputedKey := m["other_computed_value"]

	assert.Equal(t, idKey, "computed_id")
	assert.Equal(t, otherComputedKey, "computed_value_two")
}
