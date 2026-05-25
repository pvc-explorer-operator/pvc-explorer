# Releases

## Versioning model

- Local builds use `git describe --tags --always --dirty`
- Tagged releases use the semver tag itself, such as `v1.2.0`
- The release workflow pushes image tags to GHCR and attaches the generated install bundle to the GitHub Release

## Build locally

```bash
make build
./bin/manager
```

To override the version explicitly:

```bash
make build VERSION=v0.2.0-rc1
```

## Build for kind

Use the local kind scripts to rebuild the image, load it, and restart the controller pod.

```bash
./kind/setup.sh
./kind/rebuild.sh
```

## Publish a release

Push a semver tag and let GitHub Actions build and publish the release artifacts:

```bash
git tag v1.2.0
git push origin v1.2.0
```
