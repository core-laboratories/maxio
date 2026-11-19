# MaxIO Helm chart

This chart deploys the MaxIO S3-compatible object storage server. The chart is kept in this repository and is installed directly from `helm/maxio`. Packaged chart releases are not checked in.

When `image.tag` is empty, the chart uses `appVersion`. Release packages are published to `oci://ghcr.io/core-laboratories/charts/maxio` and use the matching `ghcr.io/core-laboratories/maxio` image version.

## Requirements

- Kubernetes cluster with persistent storage for production use
- Helm 4
- A MaxIO image published to a registry accessible from the cluster

## Install

```sh
helm lint helm/maxio
helm install maxio ./helm/maxio \
  --set image.repository=ghcr.io/core-laboratories/maxio \
  --set image.tag=0.1.0 \
  --set rootUser=admin \
  --set rootPassword=change-me
```

The default mode is distributed. For a local test cluster, add `--set mode=standalone --set persistence.enabled=false`.

## Upgrade

```sh
helm upgrade maxio ./helm/maxio -f my-values.yaml
```

See [values.yaml](values.yaml) for image, storage, networking, and access settings. Configuration environment variables beginning with `MINIO_` remain the server's supported keys.

## Optional IPFS integration

The IPFS extension is enabled by default. Configure its shared bucket, Kubo RPC, and Cluster connection before starting the chart, or set `MAXIO_IPFS_ENABLED=false` in `environment` to disable it. Supply `MAXIO_IPFS_*` settings through the chart `environment` values and keep Cluster credentials in deployment secrets. See [IPFS configuration](../../docs/ipfs/configuration.md) and [security](../../docs/ipfs/security.md). Deploy Kubo and IPFS Cluster separately; this chart does not install them.
