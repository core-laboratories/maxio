# IPFS architecture

MaxIO remains the authoritative S3 store. The optional IPFS worker records a durable job for each S3 object version in the internal metadata bucket, streams the S3 bytes to Kubo using CIDv1 UnixFS, verifies the CID with a second hash-only import, asks IPFS Cluster to pin with replication minimum and maximum 2, and reports `pinned` only after Cluster reports two pinned peers. The S3 ETag and object format are unchanged.

The integration is enabled by default and requires administrator configuration before the server starts. Jobs, tenant records, and revisions live in MaxIO's internal metadata bucket; namespace locks coordinate writes across servers. A periodic scan of the configured shared bucket recovers missed object events. The public gateway serves verified public jobs from MaxIO storage and never proxies arbitrary Kubo requests.

Current implementation supports public file objects only. Directory DAGs, IPNS, private publishing, CORE ID access, explicit pins, per-tenant quotas, and a full metrics surface are not implemented. Their control fields are rejected so they cannot silently weaken access. See [security.md](security.md).

[Kubo RPC](https://docs.ipfs.tech/reference/kubo/rpc/) and [IPFS Cluster REST](https://ipfscluster.io/documentation/reference/api/) are internal administrator services. Deploy two independent Kubo nodes and an IPFS Cluster with two peers. Do not expose either administrative API to customers.
