# Contributing to MaxIO

Thank you for improving MaxIO.

## Development setup

1. Install Go 1.27.4 or newer.
2. Clone `github.com/core-laboratories/maxio`.
3. Run `go build -mod=mod -o maxio .` and `go test -mod=mod ./...`.
4. Run `helm lint helm/maxio` when changing the chart.

Open a pull request with a short description of the change, its motivation, and the tests you ran. Keep changes focused and update documentation for behavior changes. Go module checksum files and temporary files are generated locally and ignored by Git.

For security reports, follow [SECURITY.md](SECURITY.md).
