package graphql

import (
	"encoding/json"
	"fmt"
	"log"

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
				Optional: true,
			},
			"compute_mutation_keys": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"compute_from_create": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"computed_update_operation_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"computed_delete_operation_variables": {
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
				Description: "Represents the state of existence of a mutation in order to support intelligent updates.",
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
		queryResponseObj, err = queryExecute(d, m, "create_mutation", "mutation_variables")
		if err != nil {
			return err
		}

		existingHash := hash(queryResponseObj)
		if err := d.Set("existing_hash", fmt.Sprint(existingHash)); err != nil {
			return err
		}

		computeFromCreate := d.Get("compute_from_create").(bool)
		if computeFromCreate {
			computeMutationVariables(queryResponseObj, d)
		}

	} else {
		computedVariables := d.Get("computed_update_operation_variables").(map[string]interface{})

		for k, v := range mutationVariables {
			computedVariables[k] = v
		}

		if err := d.Set("computed_update_operation_variables", computedVariables); err != nil {
			return err
		}

		queryResponseObj, err = queryExecute(d, m, "update_mutation", "computed_update_operation_variables")
		if err != nil {
			return err
		}
	}
	objID := hash(queryResponseObj)
	d.SetId(fmt.Sprint(objID))

	return resourceGraphqlRead(d, m)
}

func resourceGraphqlRead(d *schema.ResourceData, m interface{}) error {
	queryResponseBytes, err := queryExecute(d, m, "read_query", "read_query_variables")
	if err != nil {
		return err
	}
	if err := d.Set("query_response", string(queryResponseBytes)); err != nil {
		return err
	}

	computeFromCreate := d.Get("compute_from_create").(bool)
	if !computeFromCreate {
		computeMutationVariables(queryResponseBytes, d)
	}

	return nil
}

func resourceGraphqlMutationDelete(d *schema.ResourceData, m interface{}) error {
	_, err := queryExecute(d, m, "delete_mutation", "computed_delete_operation_variables")
	if err != nil {
		return err
	}
	return nil
}

func computeMutationVariables(queryResponseBytes []byte, d *schema.ResourceData) error {
	dataKeys := d.Get("compute_mutation_keys").(map[string]interface{})
	mutationVariables := d.Get("mutation_variables").(map[string]interface{})
	deleteMutationVariables := d.Get("delete_mutation_variables").(map[string]interface{})

	var robj = make(map[string]interface{})
	err := json.Unmarshal(queryResponseBytes, &robj)
	if err != nil {
		return err
	}

	// Set delete mutation variables
	mvks, err := computeMutationVariableKeys(dataKeys, robj)
	if err != nil {
		log.Printf("[ERROR] Unable to compute mutation variable keys: %s ", err)
	} else {
		// Combine computed delete mutation variables with provided input variables
		dvks := make(map[string]string)
		for k, v := range deleteMutationVariables {
			dvks[k] = v.(string)
		}
		for k, v := range mvks {
			dvks[k] = v
		}
		if err := d.Set("computed_delete_operation_variables", dvks); err != nil {
			return err
		}

		// Combine computed update mutation variables with provided input variables
		for k, v := range mutationVariables {
			mvks[k] = v.(string)
		}
		if err := d.Set("computed_update_operation_variables", mvks); err != nil {
			return err
		}
	}

	return nil
}
