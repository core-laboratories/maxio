// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package loggertypes

// TargetType indicates type of the target e.g. console, http, kafka
type TargetType uint8

//go:generate stringer -type=TargetType -trimprefix=Target $GOFILE

// Constants for target types
const (
	_ TargetType = iota
	TargetConsole
	TargetHTTP
	TargetKafka
)

// TargetStats contains statistics for a target.
type TargetStats struct {
	// QueueLength is the queue length if any messages are queued.
	QueueLength int

	// TotalMessages is the total number of messages sent in the lifetime of the target
	TotalMessages int64

	// FailedMessages should log message count that failed to send.
	FailedMessages int64
}
