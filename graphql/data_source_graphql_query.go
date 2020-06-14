package graphql

import (
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
			"variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: false,
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
		Read: dataSourceGraphqlQuery,
	}
}

func dataSourceGraphqlQuery(d *schema.ResourceData, m interface{}) error {
	queryResponseObj, err := QueryExecute(d, m, "query")
	if err != nil {
		return err
	}
	if err := d.Set("queryResponse", queryResponseObj); err != nil {
		return err
	}
	return nil
}
