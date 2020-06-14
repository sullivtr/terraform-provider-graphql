package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/tyler-technologies/terraform-provider-gitfile/graphql"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: graphql.Provider,
	})
}
