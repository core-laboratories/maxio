# IPFS architecture

MaxIO remains the authoritative S3 store. The optional IPFS worker records a durable job for each S3 object version in the internal metadata bucket, streams the S3 bytes to Kubo using CIDv1 UnixFS, verifies the CID with a second hash-only import, asks IPFS Cluster to pin with minimum and maximum replication set by `MAXIO_IPFS_REPLICATION`, and reports `pinned` only after that many peers report pinned copies. The S3 ETag and object format are unchanged.

The integration is disabled by default. Set `MAXIO_IPFS_ENABLED=true` and provide administrator configuration before the server starts to enable it. Jobs, tenant records, and revisions live in MaxIO's internal metadata bucket; namespace locks coordinate writes across servers. A periodic scan of the configured shared bucket recovers missed object events. The public gateway serves verified public jobs from MaxIO storage and never proxies arbitrary Kubo requests.

Current implementation supports public file objects and client-signed IPNS record objects. Directory DAGs, private publishing, CORE ID access, explicit pins, per-tenant quotas, and a full metrics surface are not implemented. Their control fields are rejected so they cannot silently weaken access. See [IPNS](ipns.md) and [security.md](security.md).

[Kubo RPC](https://docs.ipfs.tech/reference/kubo/rpc/) and [IPFS Cluster REST](https://ipfscluster.io/documentation/reference/api/) are internal administrator services. Deploy at least as many pinning Kubo nodes and Cluster peers as the requested replication factor. Do not expose either administrative API to customers.
