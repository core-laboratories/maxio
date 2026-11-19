# Public gateway

Set `MAXIO_IPFS_PUBLIC_GATEWAY=true` to start an HTTP listener on loopback. Place a TLS reverse proxy at the public edge. `GET` and `HEAD /ipfs/{cid}` serve only jobs that request public access and have two verified replicas. Single byte ranges are supported. An optional `MAXIO_IPFS_GATEWAY_DOMAIN` enables `{cid}.ipfs.{domain}` lookup.

The gateway returns generic binary downloads with `nosniff` and a sandbox content policy. It does not fetch arbitrary CIDs from Kubo, support `/ipns/`, or serve directories. It caps concurrent responses at 16 and limits bandwidth per response. It reads the matching S3 object version from MaxIO storage.
