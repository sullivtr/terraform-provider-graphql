# Building and Releasing terraform-provider-gitfile

This repository uses [goreleaser](https://goreleaser.com/) to build and publish releases

### Building:
To create a local (non-release) build, just run `make build`

### Releasing:

#### Requirements:
1. A GITHUB_TOKEN environment variable must be set ([see here](https://github.com/settings/tokens))
1. A new tag will be created with the provided commit message automatically.

```bash
export GITHUB_TOKEN=some-token-value
make publish VERSION=1.0 MESSAGE="Initial Release"
```

