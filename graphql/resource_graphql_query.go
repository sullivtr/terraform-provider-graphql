package graphql

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGraphqlMutation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mutation": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
		Read:   resourceGraphqlRead,
		Delete: resourceGraphqlMutationDelete,
	}
}

func resourceGraphqlMutationCreateUpdate(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "mutation")
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

func resourceGraphqlMutationDelete(d *schema.ResourceData, m interface{}) error {
	_, err := QueryExecute(d, m, "mutation")
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
