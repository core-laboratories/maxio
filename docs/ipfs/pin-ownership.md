# Tenants and pin ownership

For the configured shared bucket, authenticated S3 principals are mapped to a stable server-derived UUID and restricted to `users/{uuid}/`. The mapping is persisted in MaxIO internal metadata and is never accepted from client-supplied input. S3 IAM policies still apply. An administrator with bucket-owner credentials retains administrative access.

Each object version has its own job and automatic pin reference. A global Cluster unpin is attempted only when no other active object job references the CID. The same CID may be referenced by multiple tenants. Explicit pin references, per-tenant usage quotas, and tenant management endpoints are not implemented; operators should avoid treating this as a complete quota or billing boundary.
