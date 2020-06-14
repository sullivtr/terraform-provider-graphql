package graphql

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGraphqlMutation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: false,
				ForceNew: true,
			},
		},
		Create: resourceGraphqlMutationCreateUpdate,

		Update: resourceGraphqlMutationUpdate,
		Read:   resourceGraphqlRead,
		Delete: resourceGraphqlMutationDelete,
	}
}

func resourceGraphqlMutationCreate(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "createMutation")
	if err != nil {
		return err
	}
	objID := hashString(queryResponseObj)
	d.SetId(string(objID))
	return nil
}

func resourceGraphqlRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGraphqlMutationUpdate(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "updateMutation")
	if err != nil {
		return err
	}
	objID := hashString(queryResponseObj)
	d.SetId(string(objID))
	return nil
}

func resourceGraphqlMutationDelete(d *schema.ResourceData, m interface{}) error {
	_, err := QueryExecute(d, m, "deleteMutation")
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
