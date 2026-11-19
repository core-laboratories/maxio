# Private gateway

No private gateway is exposed. The implementation rejects private pin requests, since publishing plaintext through ordinary Kubo and IPFS Cluster would make content available independently of a MaxIO gateway check. A future private mode needs encryption or a truly isolated swarm plus authenticated authorization tied to an object/share record.
