package graphql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
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
			"oauth2_login_query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_login_query_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"oauth2_login_query_value_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"graphql_mutation": resourceGraphqlMutation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"graphql_query": dataSourceGraphql(),
		},
		ConfigureContextFunc: graphqlConfigure,
	}
}

func graphqlConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := &graphqlProviderConfig{
		GQLServerUrl:   d.Get("url").(string),
		RequestHeaders: d.Get("headers").(map[string]interface{}),
	}

	if oauth2LoginQueryValueAttribute := d.Get("oauth2_login_query_value_attribute"); d.Get("oauth2_login_query") != "" && oauth2LoginQueryValueAttribute != "" {
		queryResponse, resBytes, err := queryExecute(ctx, d, config, "oauth2_login_query", "oauth2_login_query_variables")
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("unable to execute oauth2_login_query: %w", err))
		}

		if queryErrors := queryResponse.ProcessErrors(); queryErrors.HasError() {
			return nil, *queryErrors
		}

		var queryResponseData map[string]interface{}
		if err := json.Unmarshal(resBytes, &queryResponseData); err != nil {
			return nil, diag.FromErr(err)
		}

		var value string
		if value, err = getOAuth2LoginQueryAttributeValue(oauth2LoginQueryValueAttribute.(string), queryResponseData); err != nil {
			return nil, diag.FromErr(err)
		}

		config.RequestAuthorizationHeaders = map[string]interface{}{
			"Authorization": fmt.Sprintf("Bearer %s", value),
		}
	}

	return config, diag.Diagnostics{}
}

type graphqlProviderConfig struct {
	GQLServerUrl   string
	RequestHeaders map[string]interface{}

	RequestAuthorizationHeaders map[string]interface{}
}

func getOAuth2LoginQueryAttributeValue(attribute string, data map[string]interface{}) (string, error) {
	resourceKeyArgs := buildResourceKeyArgs(attribute)
	value, err := getResourceKey(data, resourceKeyArgs...)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}
