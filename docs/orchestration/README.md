# Deploy MaxIO with Kubernetes

The repository includes a [Helm chart](../../helm/maxio/README.md) for Kubernetes deployments. From the repository root, install a standalone S3 deployment with IPFS disabled:

```sh
helm install maxio ./helm/maxio --set environment.MAXIO_IPFS_ENABLED=false
```

IPFS is enabled by default. For an IPFS deployment, configure Kubo and IPFS Cluster first and provide the settings described in the [IPFS configuration guide](../ipfs/configuration.md).

Tagged releases publish the chart both as a GitHub Release asset and as an [OCI package](../releases.md).
