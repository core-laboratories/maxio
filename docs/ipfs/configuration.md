# Configuration

The IPFS extension uses administrator environment variables. It is disabled by default. Set `MAXIO_IPFS_ENABLED=true` to enable it; a configured tenant bucket and authenticated Kubo/Cluster connections are then required at startup.

| Variable | Default | Purpose |
| --- | --- | --- |
| `MAXIO_IPFS_ENABLED` | `false` | Set `true` to start the worker and optional gateway. |
| `MAXIO_IPFS_TENANT_BUCKET` | unset | Shared S3 bucket containing `users/{uuid}/` keys. Required when enabled. |
| `MAXIO_IPFS_AUTO_PIN` | `false` | Pin newly created objects unless overridden. |
| `MAXIO_IPFS_IPNS_REPUBLISH` | `false` | Enable repeat publication of signed IPNS records at startup. Individual records may opt out. |
| `MAXIO_IPFS_DEFAULT_ACCESS` | `private` | Default requested visibility. Public auto pin requires `public`. |
| `MAXIO_IPFS_CLUSTER_ENDPOINT` | unset | IPFS Cluster REST URL. |
| `MAXIO_IPFS_CLUSTER_USER` | unset | Cluster basic-auth username. |
| `MAXIO_IPFS_CLUSTER_PASSWORD` | unset | Cluster basic-auth password. |
| `MAXIO_IPFS_KUBO_ENDPOINT` | unset | Kubo RPC URL for streaming imports. |
| `MAXIO_IPFS_KUBO_TOKEN` | unset | Bearer token required for non-loopback Kubo RPC. |
| `MAXIO_IPFS_KUBO_CA_FILE` | unset | Optional PEM CA file for a private HTTPS Kubo RPC certificate. The Helm Kubo option sets this automatically. |
| `MAXIO_IPFS_REPLICATION` | `2` | Positive number of Cluster peers to pin each CID. Must not exceed available pinning peers. |
| `MAXIO_IPFS_IMPORT_CHUNKER` | `size-262144` | Stable UnixFS chunker; other values currently rejected. |
| `MAXIO_IPFS_STATUS_HEADERS` | `true` | Add status fields to authorized S3 HEAD responses. |
| `MAXIO_IPFS_PUBLIC_GATEWAY` | `false` | Start the public gateway. |
| `MAXIO_IPFS_GATEWAY_LISTEN` | `127.0.0.1:8090` | Loopback listener; place a TLS reverse proxy in front. |
| `MAXIO_IPFS_GATEWAY_DOMAIN` | unset | Optional CID subdomain base. |
| `MAXIO_IPFS_GATEWAY_BYTES_PER_SECOND` | `10485760` | Per-response bandwidth cap. |

Cluster and Kubo RPC must use HTTPS, or HTTP on loopback. Cluster authentication is required; non-loopback Kubo RPC also requires a bearer token enforced by its trusted reverse proxy. Provide secrets through deployment secret storage, not a checked-in values file. In Helm, pass nonsecret settings under `environment` and mount or inject credentials from a Kubernetes Secret using the deployment mechanism available in your chart configuration.

For [IPNS publishing](ipns.md), Kubo RPC must support `name/inspect` and `name/put`. Clients hold their own signing keys and upload signed records through S3. Kubo and MaxIO receive no private IPNS keys.

Set `MAXIO_IPFS_REPLICATION=1` for one IPFS copy, or a larger positive integer for more copies. This sets both the minimum and maximum Cluster replication factors. MaxIO reports a job as pinned only after that many peers report a pinned copy. High availability of MaxIO's S3 storage does not add IPFS copies; with one IPFS replica, a failure of its pinning peer can make the CID unavailable from IPFS until recovery.

To enable automatic public pinning, set `MAXIO_IPFS_ENABLED=true`, the bucket, two RPC URLs, Cluster credentials, `MAXIO_IPFS_AUTO_PIN=true`, and `MAXIO_IPFS_DEFAULT_ACCESS=public`. Configure the gateway separately.
