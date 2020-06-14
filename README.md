# terraform-provider-graphql

## Synopsis

A [Terraform](http://terraform.io) plugin to manage graphql queries and mutations.

## Example:
```
data "graphql_query" "queryexample" {
  variables = {
    ID = "nfddksajf3948290dsa!f"
  }

  query = <<EOF
{
  secretByName(secretType: ENCRYPTED_TEXT, name:"secret-name") {
    name,
    id
  }
}
EOF          
}

resource "graphql_mutation" "mutationexample" {
  variables = {
    "secretType" = "secret_type"
    "secret" = {
      "thing1": "thing2"
    }
  }

  mutation = <<EOF
{
mutation($secret: CreateSecretInput!){
    createSecret(input: $secret){
          secret{
            id,
            name
            ... on EncryptedText{
              name
              secretManagerId
              id
            }
          usageScope{
              appEnvScopes{
                application{
                  filterType
                  appId
                }
              environment{
                  filterType
                  envId
                }
            }
          }
        }
      }
    }
}
EOF  
}
```
## Data Sources

### graphql_query

TODO: Write Docs :) 

## Resources

### graphql_mutation

TODO: Write Docs :) 


# License

Apache2 - See the included LICENSE file for more details.

