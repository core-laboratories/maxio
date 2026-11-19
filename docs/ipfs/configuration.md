# Configuration

The IPFS extension uses administrator environment variables. It is enabled by default. A configured tenant bucket and authenticated Kubo/Cluster connections are required at startup; set `MAXIO_IPFS_ENABLED=false` to disable it.

| Variable | Default | Purpose |
| --- | --- | --- |
| `MAXIO_IPFS_ENABLED` | `true` | Start worker and optional gateway. Set `false` to disable. |
| `MAXIO_IPFS_TENANT_BUCKET` | unset | Shared S3 bucket containing `users/{uuid}/` keys. Required when enabled. |
| `MAXIO_IPFS_AUTO_PIN` | `false` | Pin newly created objects unless overridden. |
| `MAXIO_IPFS_DEFAULT_ACCESS` | `private` | Default requested visibility. Public auto pin requires `public`. |
| `MAXIO_IPFS_CLUSTER_ENDPOINT` | unset | IPFS Cluster REST URL. |
| `MAXIO_IPFS_CLUSTER_USER` | unset | Cluster basic-auth username. |
| `MAXIO_IPFS_CLUSTER_PASSWORD` | unset | Cluster basic-auth password. |
| `MAXIO_IPFS_KUBO_ENDPOINT` | unset | Kubo RPC URL for streaming imports. |
| `MAXIO_IPFS_KUBO_TOKEN` | unset | Bearer token required for non-loopback Kubo RPC. |
| `MAXIO_IPFS_REPLICATION` | `2` | Must be 2. |
| `MAXIO_IPFS_IMPORT_CHUNKER` | `size-262144` | Stable UnixFS chunker; other values currently rejected. |
| `MAXIO_IPFS_STATUS_HEADERS` | `true` | Add status fields to authorized S3 HEAD responses. |
| `MAXIO_IPFS_PUBLIC_GATEWAY` | `false` | Start the public gateway. |
| `MAXIO_IPFS_GATEWAY_LISTEN` | `127.0.0.1:8090` | Loopback listener; place a TLS reverse proxy in front. |
| `MAXIO_IPFS_GATEWAY_DOMAIN` | unset | Optional CID subdomain base. |
| `MAXIO_IPFS_GATEWAY_BYTES_PER_SECOND` | `10485760` | Per-response bandwidth cap. |

Cluster and Kubo RPC must use HTTPS, or HTTP on loopback. Cluster authentication is required; non-loopback Kubo RPC also requires a bearer token enforced by its trusted reverse proxy. Provide secrets through deployment secret storage, not a checked-in values file. In Helm, pass nonsecret settings under `environment` and mount or inject credentials from a Kubernetes Secret using the deployment mechanism available in your chart configuration.

To enable automatic public pinning, set the bucket, two RPC URLs, Cluster credentials, `MAXIO_IPFS_AUTO_PIN=true`, and `MAXIO_IPFS_DEFAULT_ACCESS=public`. Configure the gateway separately.
