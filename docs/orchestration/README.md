# Deploy MaxIO with Kubernetes

The repository includes a [Helm chart](../../helm/maxio/README.md) for Kubernetes deployments. From the repository root, install a standalone S3 deployment:

```sh
helm install maxio ./helm/maxio
```

IPFS is disabled by default. For an IPFS deployment, set `environment.MAXIO_IPFS_ENABLED=true`, configure Kubo and IPFS Cluster, and provide the settings described in the [IPFS configuration guide](../ipfs/configuration.md). The Helm chart can deploy one Kubo node with `kubo.enabled=true` (also off by default); see the [chart README](../../helm/maxio/README.md). IPFS Cluster remains separate.

Tagged releases publish the chart both as a GitHub Release asset and as an [OCI package](../releases.md).
