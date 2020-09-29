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

	computeKeys := make(map[string]interface{})
	computeKeys["id_key"] = "todo.id"
	computeKeys["other_computed_value"] = "todo.otherComputedValue"

	m, err := computeMutationVariableKeys(computeKeys, robj)
	if err != nil {
		t.Fatalf("Unable to compute mutation keys from response object %v", err)
	}

	idKey := m["id_key"]
	otherComputedKey := m["other_computed_value"]

	assert.Equal(t, idKey, "computed_id")
	assert.Equal(t, otherComputedKey, "computed_value_two")
}
