package graphql

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGraphqlMutation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"read_query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_mutation": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_mutation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"update_mutation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mutation_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"read_query_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"delete_mutation_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"mutation_keys": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"computed_update_operation_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"query_response": {
				Type:        schema.TypeString,
				Description: "The raw body of the HTTP response from the last read of the object.",
				Computed:    true,
			},
			"existing_hash": {
				Type:        schema.TypeString,
				Description: "Represents the state of existence of a mutation in order to support intellegent updates.",
				Computed:    true,
			},
		},
		Create: resourceGraphqlMutationCreateUpdate,
		Update: resourceGraphqlMutationCreateUpdate,
		Read:   resourceGraphqlRead,
		Delete: resourceGraphqlMutationDelete,
	}
}

func resourceGraphqlMutationCreateUpdate(d *schema.ResourceData, m interface{}) error {
	mutationVariables := d.Get("mutation_variables").(map[string]interface{})
	var queryResponseObj []byte
	var err error
	mutationExistsHash := d.Get("existing_hash").(string)

	if mutationExistsHash == "" {
		queryResponseObj, err = QueryExecute(d, m, "create_mutation", "mutation_variables")
		if err != nil {
			return err
		}

		existingHash := hashString(queryResponseObj)
		if err := d.Set("existing_hash", string(existingHash)); err != nil {
			return err
		}

	} else {
		computedVariables := d.Get("computed_update_operation_variables").(map[string]interface{})

		for k, v := range mutationVariables {
			computedVariables[k] = v
		}

		if err := d.Set("computed_update_operation_variables", computedVariables); err != nil {
			return err
		}

		queryResponseObj, err = QueryExecute(d, m, "update_mutation", "computed_update_operation_variables")
		if err != nil {
			return err
		}
	}
	objID := hashString(queryResponseObj)
	d.SetId(string(objID))

	return resourceGraphqlRead(d, m)
}

func resourceGraphqlRead(d *schema.ResourceData, m interface{}) error {
	dataKeys := d.Get("mutation_keys").([]interface{})
	mutationVariables := d.Get("mutation_variables").(map[string]interface{})
	queryResponseBytes, err := QueryExecute(d, m, "read_query", "read_query_variables")
	if err != nil {
		return err
	}
	if err := d.Set("query_response", string(queryResponseBytes)); err != nil {
		return err
	}

	var robj = make(map[string]interface{})
	err = json.Unmarshal(queryResponseBytes, &robj)
	if err != nil {
		return err
	}

	rkas := buildResourceKeyArgs(dataKeys)

	// Set delete mutation variables
	mvks, err := computeMutationVariableKeys(rkas, robj)
	if err != nil {
		return err
	}
	if err := d.Set("delete_mutation_variables", mvks); err != nil {
		return err
	}

	// Combine computed update mutation variables with provided input variables
	for k, v := range mutationVariables {
		mvks[k] = v.(string)
	}

	if err := d.Set("computed_update_operation_variables", mvks); err != nil {
		return err
	}

	return nil
}

func resourceGraphqlMutationDelete(d *schema.ResourceData, m interface{}) error {
	_, err := QueryExecute(d, m, "delete_mutation", "delete_mutation_variables")
	if err != nil {
		return err
	}
	return nil
}
