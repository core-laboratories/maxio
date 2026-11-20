# IPNS publishing with client-held keys

IPNS signing keys stay with the client. One Kubo node can serve many users because MaxIO forwards only already signed records to its Kubo RPC. Neither MaxIO nor Kubo generates, imports, stores, or receives a user's private key.

1. Upload the target file under your `users/{uuid}/` prefix with `maxio-pin=true,maxio-access=public`. Wait for an authorized S3 HEAD response to report `x-maxio-pin-status: pinned` and read `x-maxio-cid`.
2. In your client, use your existing IPNS private key to sign a raw IPNS record whose value is `/ipfs/{cid}`. Use the same key for subsequent updates. Keep the key on your device or in your own key manager. The client must set a suitable sequence number and validity period.
3. Upload the raw signed record bytes as a separate S3 object in the same shared bucket at `users/{uuid}/.maxio/ipns/{name}`, where `{name}` is the base36 IPNS name derived from that key. The object body is the serialized signed record, not JSON or a key file. Do not attach IPFS pin metadata to this control object.
4. MaxIO checks the record signature and name through Kubo `name/inspect`, checks that its CID belongs to a public pin in the same tenant and bucket with the configured replicas, then forwards the bytes to Kubo `name/put`. An authorized S3 HEAD of the record object reports `x-maxio-pin-status: pinned` and `x-maxio-ipns-name` after publication. With the optional public gateway, `GET /ipns/{name}` serves the corresponding S3 object.

```sh
aws --endpoint-url https://s3.example.com s3api put-object \
  --bucket shared --key users/UUID/photo.jpg --body photo.jpg \
  --metadata maxio-pin=true,maxio-access=public
# Read x-maxio-cid and x-maxio-pin-status from an authorized raw S3 HEAD.
# Sign /ipfs/<cid> locally, producing record.ipns and its base36 <name>.
aws --endpoint-url https://s3.example.com s3api put-object \
  --bucket shared --key users/UUID/.maxio/ipns/<name> --body record.ipns
# After the record status is pinned: https://gateway.example.com/ipns/<name>
# Optional: opt this record out when server-side repeat publication is enabled.
aws --endpoint-url https://s3.example.com s3api put-object-tagging \
  --bucket shared --key users/UUID/.maxio/ipns/<name> \
  --tagging 'TagSet=[{Key=maxio-ipns-republish,Value=false}]'
```

Server-side repeat publication is disabled at startup unless `MAXIO_IPFS_IPNS_REPUBLISH=true` is set. When enabled, it is on for every signed IPNS record by default. A record can opt out with the `maxio-ipns-republish=false` S3 tag; omitting the tag follows the startup setting. `maxio-ipns-republish=true` is rejected. `PutObjectTagging` replaces the full tag set, so preserve any unrelated tags when changing this option. When enabled for a record, MaxIO republishes the same signed bytes every 24 hours while the record remains valid and its target remains publicly pinned. This leaves time to retry before the DHT's 48-hour record retention window. The client chooses the signed expiration time; MaxIO stops republishing at that time and never extends it. A client must upload a newly signed record before the old record expires and whenever the target CID changes. The signed record object is limited to 10 KiB. Deleting the record or unpinning the target stops local publication and gateway serving, but a record already distributed to IPNS cannot be immediately revoked from other nodes. `maxio-ipns` and `maxio-ipns-auto` metadata and tags are rejected; they belonged to the previous server-key workflow. Kubo's `name/inspect` and `name/put` RPCs are experimental, so deploy Kubo 0.40 or newer. [Kubo RPC reference](https://docs.ipfs.tech/reference/kubo/rpc/), [IPNS record specification](https://specs.ipfs.tech/ipns/ipns-record/).
