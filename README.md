<p align="center">
  <img src="assets/icon.png" alt="graphql provider" width="400"/>

  <h3 align="center">Terraform GraphQL Provider</h3>

  <p align="center">
    <a href="https://github.com/sullivtr/terraform-provider-graphql/actions?query=workflow%3Abuild"><img alt="Build" src="https://img.shields.io/github/workflow/status/sullivtr/terraform-provider-graphql/build?label=build&logo=github&style=flat-square"></a>
    <a href="https://github.com/sullivtr/terraform-provider-graphql/actions?query=workflow%3A%22acceptance+tests%22"><img alt="Acceptance Tests" src="https://img.shields.io/github/workflow/status/sullivtr/terraform-provider-graphql/acceptance%20tests?label=acceptance%20tests&logo=github&style=flat-square"></a>
    <a href="https://github.com/sullivtr/terraform-provider-graphql/actions?query=workflow%3A%22e2e+tests%22"><img alt="E2E Test" src="https://img.shields.io/github/workflow/status/sullivtr/terraform-provider-graphql/e2e%20tests?label=e2e%20test&logo=github&style=flat-square"></a>
    <a href="https://codecov.io/gh/sullivtr/terraform-provider-graphql">
      <img src="https://img.shields.io/codecov/c/github/sullivtr/terraform-provider-graphql?logo=codecov&style=flat-square" />
    </a>
    <a href="https://github.com/sullivtr/terraform-provider-graphql/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/sullivtr/terraform-provider-graphql?logo=github&style=flat-square"></a>
    <a href="https://github.com/sullivtr/terraform-provider-graphql/releases/latest"><img alt="Release" src="https://goreportcard.com/badge/github.com/sullivtr/terraform-provider-graphql?style=flat-square"></a>
  </p>
</p>

---

## Synopsis

A [Terraform](https://terraform.io) [GraphQL](https://graphql.org/) provider to manage the full lifecyle of graphql based api resources. 

It also supports optional client-side rate limiting via the `rate_limit_per_second` and `rate_limit_burst` provider arguments, for use against GraphQL APIs that enforce a request-rate limit. See the [provider documentation](https://sullivtr.github.io/terraform-provider-graphql) for details.
  
## Docs:

#### Click [HERE](https://sullivtr.github.io/terraform-provider-graphql) to view the documentation site.

## Testing this provider:
> See [Local Development](CONTRIBUTING.md#local-development) for setup details.

- To run just the unit & acceptance tests (the same command CI uses), run `make unittest`.
  This is equivalent to `go test -v ./...` on the graphql package.
  > NOTE: The acceptance tests utilize a mocked http response to simulate a graphql api. Therefore, TF_ACC is set automatically by the test init() functions.

- In the root of this project, run `make fulltest`
  This will build the plugin, copy the binaries to the basic_test/terraform.d/* folder, and run both the unit and E2E tests.

- To run just the E2E test without a full build, run `make e2etest`.
   
# License

MPL-2 - See the included LICENSE file for more details.


