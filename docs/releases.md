# Releases

MaxIO releases use stable tags in the form `vMAJOR.MINOR.PATCH`, beginning with `v0.1.0`. The release workflow runs only when a matching tag is pushed, and verifies that the tag commit is on `master` and that `helm/maxio/Chart.yaml` has the matching `version` and `appVersion` (without the `v`). Update both chart fields before tagging.

For example, after the changes are committed and the checks pass:

```sh
git tag -a v0.1.0 -m 'MaxIO v0.1.0'
git push origin master
git push origin v0.1.0
```

You can also create a GitHub Release for `v0.1.0` in the GitHub UI. If the tag does not exist yet, GitHub creates it when you publish the release. Select the commit on `master` containing the matching chart version. The workflow uploads its build assets to that release when it already exists; otherwise it creates the release. Re-running the publish job replaces assets with the same names.

The workflow runs Go tests, govulncheck, and Helm lint, cross-builds Linux, macOS, and Windows binaries for amd64 and arm64, creates archives, and publishes a GitHub Release. The release assets include the packaged Helm chart (`maxio-0.1.0.tgz`), the binary archives, and Docker image archives for Linux amd64 and arm64 (`maxio-0.1.0-linux-<arch>-docker.tar.gz`). The workflow also pushes the multi-architecture image as `ghcr.io/core-laboratories/maxio:0.1.0` and the Helm chart as `oci://ghcr.io/core-laboratories/charts/maxio` version `0.1.0`. The source chart keeps `image.tag` empty so its `appVersion` selects the matching image. No chart package, checksum file, or repository index is committed.

Pull the image with `docker pull ghcr.io/core-laboratories/maxio:0.1.0`, or load a downloaded release archive with `gzip -dc maxio-0.1.0-linux-amd64-docker.tar.gz | docker load` (use the archive for your CPU architecture).

The release binary's `Version` remains an RFC 3339 timestamp for compatibility with the existing server code. Its displayed `ReleaseTag` is the Git tag, such as `v0.1.0`. The Docker build accepts the release linker flags through `MAXIO_LDFLAGS`; ordinary local builds keep the development version.

Go 1.27.4 or newer is required; the workflows use the latest 1.27 patch and the Docker image follows the 1.27 Alpine 3.24 patch line.

The workflow uses the repository's `GITHUB_TOKEN`. Allow GitHub Actions and the named actions in repository or organization settings. The publishing jobs request `contents: write` for the GitHub Release and `packages: write` for GHCR; no personal access token is required for publishing to this repository's packages. Set the container and chart package visibility to public if anonymous installs are desired. Configure branch and tag protections in GitHub settings as appropriate for your release process.

The container smoke workflow runs with IPFS disabled by default. Deployments that enable IPFS must provide the configuration described in [IPFS configuration](ipfs/configuration.md). Publishing a release does not provision Kubo or IPFS Cluster.

The inherited in-place update checker targets the original timestamp release feed and is incompatible with these SemVer releases. MaxIO now disables that checker and its administrative in-place update operation. Download release assets from GitHub Releases until a signed MaxIO update channel is implemented.
