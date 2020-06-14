package graphql

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGraphqlMutation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"readQuery": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"createMutation": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deleteMutation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"updateMutation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"createMutationVariables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},
			"updateMutationVariables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"readQueryVariables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"deleteMutationVariables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"queryResponse": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
		Create: resourceGraphqlMutationCreate,
		Update: resourceGraphqlMutationUpdate,
		Read:   resourceGraphqlRead,
		Delete: resourceGraphqlMutationDelete,
	}
}

func resourceGraphqlMutationCreate(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "createMutation", "createMutationVariables")
	if err != nil {
		return err
	}
	objID := hashString(queryResponseObj)
	d.SetId(string(objID))
	return resourceGraphqlRead(d, m)
}

func resourceGraphqlRead(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "readQuery", "readQueryVariables")
	if err != nil {
		return err
	}
	if err := d.Set("queryResponse", queryResponseObj); err != nil {
		return err
	}
	return nil
}

func resourceGraphqlMutationUpdate(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "updateMutation", "updateMutationVariables")
	if err != nil {
		return err
	}
	objID := hashString(queryResponseObj)
	d.SetId(string(objID))
	return nil
}

func resourceGraphqlMutationDelete(d *schema.ResourceData, m interface{}) error {
	_, err := QueryExecute(d, m, "deleteMutation", "deleteMutationVariables")
	if err != nil {
		return err
	}
	return nil
}

func hashString(v map[string]interface{}) int {
	out, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return hashcode.String(string(out))
}
