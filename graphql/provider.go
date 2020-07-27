package graphql

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_GRAPHQL_URL", nil),
			},
			"headers": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"graphql_mutation": resourceGraphqlMutation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"graphql_query": dataSourceGraphql(),
		},
		ConfigureFunc: graphqlConfigure,
	}
}

func graphqlConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &graphqlProviderConfig{
		GQLServerUrl:   d.Get("url").(string),
		RequestHeaders: d.Get("headers").(map[string]interface{}),
	}
	return config, nil
}

type graphqlProviderConfig struct {
	GQLServerUrl   string
	RequestHeaders map[string]interface{}
}
