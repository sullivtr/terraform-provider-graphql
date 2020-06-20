package graphql

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/hashcode"
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
			"create_mutation_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},
			"update_mutation_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
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
			"query_response": {
				Type:        schema.TypeString,
				Description: "The raw body of the HTTP response from the last read of the object.",
				Computed:    true,
			},
		},
		Create: resourceGraphqlMutationCreate,
		Update: resourceGraphqlMutationUpdate,
		Read:   resourceGraphqlRead,
		Delete: resourceGraphqlMutationDelete,
	}
}

func resourceGraphqlMutationCreate(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "create_mutation", "create_mutation_variables")
	if err != nil {
		return err
	}
	objID := hashString(queryResponseObj)
	d.SetId(string(objID))
	return resourceGraphqlRead(d, m)
}

func resourceGraphqlRead(d *schema.ResourceData, m interface{}) error {
	queryResponseBytes, err := QueryExecute(d, m, "read_query", "read_query_variables")
	if err != nil {
		return err
	}
	if err := d.Set("query_response", string(queryResponseBytes)); err != nil {
		return err
	}
	return nil
}

func resourceGraphqlMutationUpdate(d *schema.ResourceData, m interface{}) error {
	queryResponseBytes, err := QueryExecute(d, m, "update_mutation", "update_mutation_variables")
	if err != nil {
		return err
	}
	objID := hashString(queryResponseBytes)
	d.SetId(string(objID))
	return nil
}

func resourceGraphqlMutationDelete(d *schema.ResourceData, m interface{}) error {
	_, err := QueryExecute(d, m, "delete_mutation", "delete_mutation_variables")
	if err != nil {
		return err
	}
	return nil
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
