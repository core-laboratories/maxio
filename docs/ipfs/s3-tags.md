# S3 tag controls

Standard `PutObjectTagging`, `GetObjectTagging`, and `DeleteObjectTagging` remain the control interface. `maxio-pin=true|false` and `maxio-access=public|private` request a pin or visibility policy for that exact object version. Tags replace the previous tag set using ordinary S3 semantics; missing IPFS tags inherit metadata and server defaults. MaxIO's pin status is stored separately from tags and exposed through S3 HEAD.

`maxio-ipns` and `maxio-ipns-auto` tags are rejected. Upload a [client-signed record object](ipns.md) to publish IPNS. The startup setting `MAXIO_IPFS_IPNS_REPUBLISH` controls the default for all records; on a signed-record object, `maxio-ipns-republish=false` opts out. Omitting the tag follows the startup setting, and `true` is rejected. The ordinary S3 10-tag limit still applies. Other unrelated S3 tags are unaffected.
