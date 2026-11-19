# S3 metadata controls

Standard S3 `x-amz-meta-*` headers carry requested IPFS policy. `x-amz-meta-maxio-pin` accepts `true` or `false`; `x-amz-meta-maxio-access` accepts `public` or `private`. A public object can be pinned through the extension. Private plaintext pinning is rejected because a CID and public-network pin would expose the bytes beyond MaxIO's HTTP gateway.

`x-amz-meta-maxio-ipns`, `x-amz-meta-maxio-ipns-auto`, and `x-amz-meta-maxio-core-ids` are reserved and currently rejected. The CORE ID field is removed from S3 GET/HEAD metadata responses when IPFS is enabled. Metadata is only requested policy, never an authorization grant. Copying metadata from another object does not grant IPFS controls. S3 ETags remain ordinary S3 ETags, not CIDs.

A user must still have the normal S3 write and tag permissions. The tenant prefix restriction is enforced through IAM request authorization for the configured shared bucket.
