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

The IPFS extension is disabled by default. Set `MAXIO_IPFS_ENABLED=true` in `environment` and configure its shared bucket and IPFS Cluster connection to enable it. Supply `MAXIO_IPFS_*` settings through the chart `environment` values and keep Cluster credentials in deployment secrets. See [IPFS configuration](../../docs/ipfs/configuration.md) and [security](../../docs/ipfs/security.md).

The chart can deploy one persistent Kubo node with `kubo.enabled=true`; the default is `false`. This is independent of `MAXIO_IPFS_ENABLED`, so setting either flag alone does not enable the other. The chart connects MaxIO to the optional node over HTTPS and injects `MAXIO_IPFS_KUBO_ENDPOINT`, `MAXIO_IPFS_KUBO_TOKEN`, and `MAXIO_IPFS_KUBO_CA_FILE`. Do not also set those keys in `environment` or `extraSecret`.

Before enabling Kubo, create two Kubernetes Secrets in the release namespace:

- `kubo.tls.existingSecret` must contain `tls.crt`, `tls.key`, and `ca.crt`. The server certificate must cover `<kubo-service>.<namespace>.svc`, where the service is `<release>-maxio-kubo` unless the chart name is overridden. `ca.crt` is the CA certificate trusted by MaxIO.
- `kubo.auth.existingSecret` must contain a bearer token under the key selected by `kubo.auth.tokenKey` (default `token`). Use a high-entropy token.

For example, after preparing certificate files and a token Secret:

```sh
helm install maxio ./helm/maxio \
  --set kubo.enabled=true \
  --set kubo.tls.existingSecret=maxio-kubo-tls \
  --set kubo.auth.existingSecret=maxio-kubo-auth
```

This deploys Kubo while leaving the MaxIO IPFS extension off. Enable the extension after supplying `MAXIO_IPFS_ENABLED=true`, `MAXIO_IPFS_TENANT_BUCKET`, the IPFS Cluster endpoint and credentials, and a suitable replication factor through `environment` and deployment Secrets. Kubo's admin RPC and built-in gateway bind only to loopback; the RPC Service exposes a TLS proxy that requires the bearer token. A NetworkPolicy limits RPC ingress to MaxIO pods in the same namespace and release; your CNI must enforce NetworkPolicy for that limit to apply. If you enable egress restrictions on MaxIO, allow TCP 8443 to the Kubo pods as well as traffic to IPFS Cluster. The Kubo repo uses a separate PVC (`kubo.persistence.size`, default `100Gi`), or set `kubo.persistence.existingClaim`. The swarm service defaults to ClusterIP; change `kubo.swarmService.type` and configure Kubo's announced swarm address if peers outside the cluster must reach it. The chart does not deploy IPFS Cluster. Configure its endpoint and credentials separately, and ensure it has at least as many pinning peers as `MAXIO_IPFS_REPLICATION`; a single Kubo node does not provide multiple IPFS replicas.
