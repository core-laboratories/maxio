# Automatic pinning

After successful `PutObject`, `CopyObject`, multipart completion, tagging, and applicable restore notifications, MaxIO records a durable job for the exact object version. The worker streams content to Kubo, checks the CID with a hash-only import, requests a two-peer Cluster pin, and verifies two pinned peers. Failures retry with exponential backoff and enter a dead state after 12 attempts. S3 writes succeed even when the IPFS worker or Cluster is unavailable.

Automatic pinning is off by default. It requires public access in this implementation. Object metadata and tags can override the server setting. A tag update to `maxio-pin=false` releases the automatic reference. The status header can report `pending`, `pinning`, `pinned`, `failed`, `dead`, or `unpinned`.
