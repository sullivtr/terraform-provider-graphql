package e2e

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestOAuth2CreateUpdateMutations(t *testing.T) {

	varFileCreate := []string{"./variable_initial_create.tfvars"}
	varFileUpdate := []string{"./variable_update.tfvars"}
	terraformOptionsCreate := &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "./test_oauth2",
		VarFiles:     varFileCreate,
		Logger:       logger.Discard,
	}

	terraformOptionsUpdate := &terraform.Options{
		TerraformDir: "./test_oauth2",
		VarFiles:     varFileUpdate,
		Logger:       logger.Discard,
	}

	// Ensure workspace is clean
	assert.NoFileExists(t, "./gql-server/test.json")
	assert.NoFileExists(t, "./gql-server/loginAPI.json")

	terraform.InitAndApply(t, terraformOptionsCreate)

	// Validate creation
	assert.FileExists(t, "./gql-server/test.json")
	assert.FileExists(t, "./gql-server/loginAPI.json")

	// Validate data source output
	dataSourceOutput, _ := terraform.OutputJsonE(t, terraformOptionsCreate, "query_output")
	assert.Contains(t, dataSourceOutput, initialTextOutput)

	// Validate computed delete variables
	deleteVariableOutput, _ := terraform.OutputJsonE(t, terraformOptionsCreate, "computed_delete_variables")
	assert.Contains(t, deleteVariableOutput, idOutPut)
	assert.Contains(t, deleteVariableOutput, testVarComputed)

	// Run update & validate changes
	terraform.InitAndApply(t, terraformOptionsUpdate)
	updatedOutput, _ := terraform.OutputJsonE(t, terraformOptionsUpdate, "query_output")
	assert.Contains(t, updatedOutput, updatedTextOutput)
	assert.NotContains(t, updatedOutput, initialTextOutput)

	terraform.Destroy(t, terraformOptionsUpdate)
	assert.NoFileExists(t, "./gql-server/test.json")
	assert.FileExists(t, "./gql-server/loginAPI.json")
	os.Remove("./gql-server/loginAPI.json")
}

func TestOAuth2ForceReplace(t *testing.T) {

	varFileCreate := []string{"./variable_initial_create.tfvars"}
	varFileUpdate := []string{"./variable_force_replace_update.tfvars"}
	terraformOptionsCreate := &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "./test_oauth2",
		VarFiles:     varFileCreate,
		Logger:       logger.Discard,
	}

	terraformOptionsUpdate := &terraform.Options{
		TerraformDir: "./test_oauth2",
		VarFiles:     varFileUpdate,
		Logger:       logger.Discard,
	}

	// Ensure workspace is clean
	assert.NoFileExists(t, "./gql-server/test.json")
	assert.NoFileExists(t, "./gql-server/loginAPI.json")

	terraform.InitAndApply(t, terraformOptionsCreate)

	// Validate creation
	assert.FileExists(t, "./gql-server/test.json")
	assert.FileExists(t, "./gql-server/loginAPI.json")

	// Validate data source output
	dataSourceOutput, _ := terraform.OutputJsonE(t, terraformOptionsCreate, "query_output")
	assert.Contains(t, dataSourceOutput, initialTextOutput)

	// Validate computed delete variables
	deleteVariableOutput, _ := terraform.OutputJsonE(t, terraformOptionsCreate, "computed_delete_variables")
	assert.Contains(t, deleteVariableOutput, idOutPut)
	assert.Contains(t, deleteVariableOutput, testVarComputed)

	// Run update & validate changes
	terraform.InitAndApply(t, terraformOptionsUpdate)
	updatedOutput, _ := terraform.OutputJsonE(t, terraformOptionsUpdate, "query_output")
	assert.Contains(t, updatedOutput, updatedTextOutputForceReplace)
	assert.NotContains(t, updatedOutput, initialTextOutput)

	terraform.Destroy(t, terraformOptionsUpdate)
	assert.NoFileExists(t, "./gql-server/test.json")
	assert.FileExists(t, "./gql-server/loginAPI.json")
	os.Remove("./gql-server/loginAPI.json")
}

func TestOAuth2ValidateComputeMutationKeysFromCreate(t *testing.T) {

	varFileComputeFromCreate := []string{"./variable_compute_from_create.tfvars"}

	terraformOptionsComputeFromCreate := &terraform.Options{
		TerraformDir: "./test_oauth2",
		VarFiles:     varFileComputeFromCreate,
		Logger:       logger.Discard,
	}

	// Ensure workspace is clean
	assert.NoFileExists(t, "./gql-server/test.json")
	assert.NoFileExists(t, "./gql-server/loginAPI.json")

	terraform.InitAndApply(t, terraformOptionsComputeFromCreate)

	// Validate compute mutation keys from create
	assert.FileExists(t, "./gql-server/test.json")
	assert.FileExists(t, "./gql-server/loginAPI.json")

	// Validate data source output
	dataSourceOutput, _ := terraform.OutputJsonE(t, terraformOptionsComputeFromCreate, "query_output")
	assert.Contains(t, dataSourceOutput, initialTextOutput)

	// Validate computed delete variables
	readVariableUpdate, _ := terraform.OutputJsonE(t, terraformOptionsComputeFromCreate, "computed_read_variables")
	assert.Contains(t, readVariableUpdate, idOutPut)
	assert.Contains(t, readVariableUpdate, testVarComputed)

	// Validate computed delete variables
	deleteVariableOutput, _ := terraform.OutputJsonE(t, terraformOptionsComputeFromCreate, "computed_delete_variables")
	assert.Contains(t, deleteVariableOutput, idOutPut)
	assert.Contains(t, deleteVariableOutput, testVarComputed)

	terraform.Destroy(t, terraformOptionsComputeFromCreate)
	assert.NoFileExists(t, "./gql-server/test.json")
	assert.FileExists(t, "./gql-server/loginAPI.json")
	os.Remove("./gql-server/loginAPI.json")
}
