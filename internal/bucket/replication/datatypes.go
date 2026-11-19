// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package replication

//go:generate msgp -file=$GOFILE

// StatusType of Replication for x-amz-replication-status header
type StatusType string

// Type - replication type enum
type Type int

const (
	// Pending - replication is pending.
	Pending StatusType = "PENDING"

	// Completed - replication completed ok.
	Completed StatusType = "COMPLETED"

	// CompletedLegacy was called "COMPLETE" incorrectly.
	CompletedLegacy StatusType = "COMPLETE"

	// Failed - replication failed.
	Failed StatusType = "FAILED"

	// Replica - this is a replica.
	Replica StatusType = "REPLICA"
)

// String returns string representation of status
func (s StatusType) String() string {
	return string(s)
}

// Empty returns true if this status is not set
func (s StatusType) Empty() bool {
	return string(s) == ""
}

// VersionPurgeStatusType represents status of a versioned delete or permanent delete w.r.t bucket replication
type VersionPurgeStatusType string

const (
	// VersionPurgePending - versioned delete replication is pending.
	VersionPurgePending VersionPurgeStatusType = "PENDING"

	// VersionPurgeComplete - versioned delete replication is now complete, erase version on disk.
	VersionPurgeComplete VersionPurgeStatusType = "COMPLETE"

	// VersionPurgeFailed - versioned delete replication failed.
	VersionPurgeFailed VersionPurgeStatusType = "FAILED"
)

// Empty returns true if purge status was not set.
func (v VersionPurgeStatusType) Empty() bool {
	return string(v) == ""
}

// Pending returns true if the version is pending purge.
func (v VersionPurgeStatusType) Pending() bool {
	return v == VersionPurgePending || v == VersionPurgeFailed
}
