# MaxIO documentation

MaxIO is an S3-compatible object storage server. Build it from the repository root with `go build -mod=mod -o maxio .` and start an S3-only instance with `MAXIO_IPFS_ENABLED=false ./maxio server /data`. Configure IPFS before starting with the default enabled setting.

The [root README](../README.md) covers the main build, Docker, and Helm commands. See [releases](releases.md) for publishing and [IPFS](ipfs/README.md) for the extension. The topic directories here contain examples for integrations and operations. Some examples use upstream S3 clients or legacy environment variable names that remain supported by the server.
