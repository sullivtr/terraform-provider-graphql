<p align="center">
  <img src="assets/icon.png" alt="graphql provider" width="400"/>

  <h3 align="center">Terraform GraphQL Provider</h3>

  <p align="center">
    <a href="https://github.com/sullivtr/terraform-provider-graphql/actions?query=workflow%3Abuild"><img alt="Build" src="https://github.com/sullivtr/terraform-provider-graphql/workflows/build/badge.svg"></a>
    <a href="https://github.com/sullivtr/terraform-provider-graphql/actions?query=workflow%3A%22acceptance+tests%22"><img alt="Acceptance Tests" src="https://github.com/sullivtr/terraform-provider-graphql/workflows/acceptance%20tests/badge.svg"></a>
    <a href="https://github.com/sullivtr/terraform-provider-graphql/actions?query=workflow%3A%22e2e+tests%22"><img alt="E2E Test" src="https://github.com/sullivtr/terraform-provider-graphql/workflows/e2e%20tests/badge.svg"></a>
    <a href="https://github.com/sullivtr/terraform-provider-graphql/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/sullivtr/terraform-provider-graphql"></a>
    <a href="https://img.shields.io/github/downloads/sullivtr/terraform-provider-graphql/total?color=orange"><img alt="GitHub release" src="https://img.shields.io/github/downloads/sullivtr/terraform-provider-graphql/total?color=orange"></a>
    <a href="https://img.shields.io/github/last-commit/sullivtr/terraform-provider-graphql?color=ff69b4"><img alt="GitHub release" src="https://img.shields.io/github/last-commit/sullivtr/terraform-provider-graphql?color=ff69b4"></a>
  </p>
</p>

---

## Synopsis

A [Terraform](https://terraform.io) [GraphQL](https://graphql.org/) provider to manage the full lifecyle of graphql based api resources. 
  
## Docs:

#### Click [HERE](https://sullivtr.github.io/terraform-provider-graphql) to view the documentation site.

## Testing this provider:
- In the root of this project, run `make fulltest`
  This will build the plugin, and copy the binaries to the basic_test/terraform.d/* folder

- To run a test without a build, simply run `make test`
   
# License

Apache2 - See the included LICENSE file for more details.


