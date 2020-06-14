package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sullivtr/terraform-provider-graphql/graphql"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: graphql.Provider,
	})
}
