# terraform-provider-graphql

## Synopsis

A [Terraform](http://terraform.io) plugin to manage graphql queries and mutations.

## Example:
```
provider "graphql" {
  url = "https://my-graphql-service-url.io"
  headers = {
    "x-api-key": "4324nsfkdsanj32k!!4FakeApiKey8873"
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
  }
}
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

