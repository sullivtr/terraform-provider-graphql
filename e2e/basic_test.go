package e2e

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestBasicCreateUpdateMutations(t *testing.T) {

	initialTextOutput := "\\\"text\\\":\\\"Here is something todo\\\""
	updatedTextOutput := "\\\"text\\\":\\\"Todo has been updated\\\""

	varFileCreate := []string{"./variable_initial_create.tfvars"}
	varFileUpdate := []string{"./variable_update.tfvars"}
	terraformOptionsCreate := &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "./test_basic",
		VarFiles:     varFileCreate,
		Logger:       logger.Discard,
	}

	terraformOptionsUpdate := &terraform.Options{
		TerraformDir: "./test_basic",
		VarFiles:     varFileUpdate,
		Logger:       logger.Discard,
	}

	terraform.InitAndApply(t, terraformOptionsCreate)
	// Validate creation
	assert.FileExists(t, "./gql-server/test.json")
	output, _ := terraform.OutputE(t, terraformOptionsCreate, "mutation_output")
	assert.Contains(t, output, initialTextOutput)

	// Validate data source output
	dataSourceOutput, _ := terraform.OutputE(t, terraformOptionsCreate, "query_output")
	assert.Contains(t, dataSourceOutput, initialTextOutput)

	// Validate computed delete variables
	deleteVariableOutput, _ := terraform.OutputE(t, terraformOptionsCreate, "computed_delete_variables")
	assert.Contains(t, deleteVariableOutput, "\"id\" =")
	assert.Contains(t, deleteVariableOutput, "\"testvar1\" =")

	// Run update & validate changes
	terraform.InitAndApply(t, terraformOptionsUpdate)
	output, _ = terraform.OutputE(t, terraformOptionsUpdate, "mutation_output")
	assert.Contains(t, output, updatedTextOutput)
	assert.NotContains(t, output, initialTextOutput)

	terraform.Destroy(t, terraformOptionsUpdate)
	assert.NoFileExists(t, "./gql-server/test.json")
}
