# Security boundaries

IPFS CIDs are public identifiers, not secrets. Public pinning is an explicit policy choice. IPFS Cluster and Kubo RPC are administrative interfaces: protect them with network isolation, TLS or loopback transport, and authentication. Do not expose them through the public gateway. The gateway serves only public jobs that meet the configured replica count and reads bytes from MaxIO; it cannot proxy arbitrary CIDs.

Private plaintext publishing and CORE ID controls are rejected while their security model is incomplete. IPNS keys stay with clients. MaxIO accepts signed records only for a public pin owned by the same tenant; public IPNS publications cannot be immediately revoked from other nodes. Authoritative job and tenant records are stored internally, separate from client-editable S3 metadata and tags. A copied object's metadata is not automatically treated as authorization. Existing S3 IAM and Object Lock enforcement remain in the normal S3 path.

This implementation does not yet provide per-tenant quotas, audit records dedicated to IPFS operations, encryption for private distribution, or complete race-free reference accounting. Those are deployment blockers for a public multi-tenant private-content service.
