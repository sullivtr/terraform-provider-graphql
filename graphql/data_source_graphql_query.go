package graphql

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGraphql() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"query_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},
			"query_response": {
				Type:        schema.TypeString,
				Description: "The raw body of the HTTP response from the last read of the object.",
				Computed:    true,
			},
		},
		Read: dataSourceGraphqlQuery,
	}
}

func dataSourceGraphqlQuery(d *schema.ResourceData, m interface{}) error {
	queryResponseBytes, err := queryExecute(d, m, "query", "query_variables")
	if err != nil {
		return err
	}
	objID := hash(queryResponseBytes)
	d.SetId(fmt.Sprint(objID))
	if err := d.Set("query_response", string(queryResponseBytes)); err != nil {
		return err
	}
	return nil
}
