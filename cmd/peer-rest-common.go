// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import "time"

const (
	peerRESTVersion       = "v39" // add more flags to speedtest API
	peerRESTVersionPrefix = SlashSeparator + peerRESTVersion
	peerRESTPrefix        = minioReservedBucketPath + "/peer"
	peerRESTPath          = peerRESTPrefix + peerRESTVersionPrefix
)

const (
	peerRESTMethodHealth                = "/health"
	peerRESTMethodVerifyBinary          = "/verifybinary"
	peerRESTMethodCommitBinary          = "/commitbinary"
	peerRESTMethodStartProfiling        = "/startprofiling"
	peerRESTMethodDownloadProfilingData = "/downloadprofilingdata"
	peerRESTMethodSpeedTest             = "/speedtest"
	peerRESTMethodDriveSpeedTest        = "/drivespeedtest"
	peerRESTMethodDevNull               = "/devnull"
	peerRESTMethodNetperf               = "/netperf"
	peerRESTMethodGetReplicationMRF     = "/getreplicationmrf"
)

const (
	peerRESTBucket          = "bucket"
	peerRESTBuckets         = "buckets"
	peerRESTUser            = "user"
	peerRESTGroup           = "group"
	peerRESTUserTemp        = "user-temp"
	peerRESTPolicy          = "policy"
	peerRESTUserOrGroup     = "user-or-group"
	peerRESTUserType        = "user-type"
	peerRESTIsGroup         = "is-group"
	peerRESTSignal          = "signal"
	peerRESTSubSys          = "sub-sys"
	peerRESTProfiler        = "profiler"
	peerRESTSize            = "size"
	peerRESTConcurrent      = "concurrent"
	peerRESTDuration        = "duration"
	peerRESTStorageClass    = "storage-class"
	peerRESTEnableSha256    = "enableSha256"
	peerRESTEnableMultipart = "enableMultipart"
	peerRESTAccessKey       = "access-key"
	peerRESTMetricsTypes    = "types"
	peerRESTDisk            = "disk"
	peerRESTHost            = "host"
	peerRESTJobID           = "job-id"
	peerRESTDepID           = "depID"
	peerRESTStartRebalance  = "start-rebalance"
	peerRESTMetrics         = "metrics"
	peerRESTDryRun          = "dry-run"
	peerRESTUploadID        = "up-id"

	peerRESTURL         = "url"
	peerRESTSha256Sum   = "sha256sum"
	peerRESTReleaseInfo = "releaseinfo"
	peerRESTExecAt      = "exec-at"

	peerRESTListenBucket = "bucket"
	peerRESTListenPrefix = "prefix"
	peerRESTListenSuffix = "suffix"
	peerRESTListenEvents = "events"
	peerRESTLogMask      = "log-mask"
)

const restartUpdateDelay = 250 * time.Millisecond
