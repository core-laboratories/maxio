# Recovery

Jobs are persisted in the internal metadata bucket and use revision checks under MaxIO namespace locks. Workers retry network failures with exponential backoff up to 12 attempts. A startup scan and a 15-minute scan of the shared tenant bucket discover missed writes and removed versions. The scan reconstructs a pin-disabled record for missed events; client metadata copied without a validated event is not trusted as a privileged grant. Operators must reapply an authorized tag update to pin these recovered objects.

Status is `pinned` only after Cluster reports the configured number of peers in the pinned state. Pinned jobs are periodically rechecked to detect a lost replica. Dead jobs require operator investigation and currently have no dedicated replay API. Preserve the internal metadata bucket during backup and restore.

The scan recovers client-signed IPNS record objects and retries publication. Repeat publication is disabled by default at startup; set `MAXIO_IPFS_IPNS_REPUBLISH=true` to enable it for records without an opt-out tag. The worker then republishes valid signed bytes every 24 hours without a private key. Clients must sign a new record before expiry and back up their own private keys; MaxIO cannot extend a signed record's validity. See [IPNS publishing](ipns.md).
