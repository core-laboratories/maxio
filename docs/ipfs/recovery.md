# Recovery

Jobs are persisted in the internal metadata bucket and use revision checks under MaxIO namespace locks. Workers retry network failures with exponential backoff up to 12 attempts. A startup scan and a 15-minute scan of the shared tenant bucket discover missed writes and removed versions. The scan reconstructs a pin-disabled record for missed events; client metadata copied without a validated event is not trusted as a privileged grant. Operators must reapply an authorized tag update to pin these recovered objects.

Status is `pinned` only after Cluster reports two peers in the pinned state. Pinned jobs are periodically rechecked to detect a lost replica. Dead jobs require operator investigation and currently have no dedicated replay API. Preserve the internal metadata bucket during backup and restore.
