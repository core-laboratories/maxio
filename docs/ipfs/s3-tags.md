# S3 tag controls

Standard `PutObjectTagging`, `GetObjectTagging`, and `DeleteObjectTagging` remain the control interface. `maxio-pin=true|false` and `maxio-access=public|private` request a pin or visibility policy for that exact object version. Tags replace the previous tag set using ordinary S3 semantics; missing IPFS tags inherit metadata and server defaults. MaxIO's pin status is stored separately from tags and exposed through S3 HEAD.

`maxio-ipns` and `maxio-ipns-auto` are reserved and rejected until managed IPNS is implemented. The ordinary S3 10-tag limit still applies. Other unrelated S3 tags are unaffected.
