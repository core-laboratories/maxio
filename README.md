# MaxIO

MaxIO is an S3-compatible object storage server.

![MaxIO Icon](./icon.webp)

## Requirements

- Go 1.27.4 or newer
- Docker for container builds, or Helm for Kubernetes deployment

## Build and run

```sh
go build -mod=mod -o maxio .
MAXIO_IPFS_ENABLED=false ./maxio server /path/to/data --console-address :9001
```

The API listens on port 9000 and the console on port 9001. Set `MAXIO_ROOT_USER` and `MAXIO_ROOT_PASSWORD` to choose credentials.

## Container

```sh
docker build -t maxio:local .
docker run --rm -p 9000:9000 -p 9001:9001 -v maxio-data:/data \
  -e MAXIO_ROOT_USER=admin -e MAXIO_ROOT_PASSWORD=change-me \
  -e MAXIO_IPFS_ENABLED=false \
  maxio:local
```

The Dockerfile builds the server from this source tree and runs it as a non-root user on Alpine Linux. Rebuild the image regularly to pick up Alpine package fixes.

## Kubernetes

```sh
helm install maxio ./helm/maxio --set environment.MAXIO_IPFS_ENABLED=false
```

See [the chart README](helm/maxio/README.md) and [documentation](docs/README.md) for configuration and deployment details.

## Releases

The release workflow builds binaries for Linux, macOS, and Windows, publishes a multi-platform image to GHCR, and pushes the Helm chart as an OCI artifact. See [release instructions](docs/releases.md).

## Development

Run `go test -mod=mod ./...` for the Go test suite and `helm lint helm/maxio` for chart validation. Dependency checksums are generated locally and are not committed.

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidance and [SECURITY.md](SECURITY.md) for security reporting.

## Optional IPFS integration

MaxIO can optionally import public S3 objects to Kubo and request two replicas through IPFS Cluster. The extension is enabled by default and requires IPFS configuration before startup; set `MAXIO_IPFS_ENABLED=false` to disable it. The quickstart commands above explicitly disable IPFS for a standalone S3 deployment. See [IPFS architecture](docs/ipfs/architecture.md), [configuration](docs/ipfs/configuration.md), and [security limits](docs/ipfs/security.md) before deployment.

## License

S3 object storage server is licensed under the GNU AGPL v3.0. See [LICENSE](LICENSE) for the full terms.
