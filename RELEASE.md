# Building and Releasing terraform-provider-graphql

This repository uses [goreleaser](https://goreleaser.com/) to build and publish releases

### Building:
To create a local (non-release) build, just run `make build`

### Releasing:

#### Requirements:
1. Releases can only be published from master
2. Releases are only published when a tag with the following schema is added: `v.*.*.*`
3. All builds and tests must be passing prior to any merges to the master branch. 
4. After merging a pull request on master, add a tag for the new release. Once the tag is added, the `release` github action will automatically publish a release with the tag.


