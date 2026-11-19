// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package dsync

//go:generate msgp -file $GOFILE

// LockArgs is minimal required values for any dsync compatible lock operation.
type LockArgs struct {
	// Unique ID of lock/unlock request.
	UID string

	// Resources contains single or multiple entries to be locked/unlocked.
	Resources []string

	// Owner represents unique ID for this instance, an owner who originally requested
	// the locked resource, useful primarily in figuring out stale locks.
	Owner string

	// Source contains the line number, function and file name of the code
	// on the client node that requested the lock.
	Source string `msgp:"omitempty"`

	// Quorum represents the expected quorum for this lock type.
	Quorum *int `msgp:"omitempty"`
}

// ResponseCode is the response code for a locking request.
type ResponseCode uint8

// Response codes for a locking request.
const (
	RespOK ResponseCode = iota
	RespLockConflict
	RespLockNotInitialized
	RespLockNotFound
	RespErr
)

// LockResp is a locking request response.
type LockResp struct {
	Code ResponseCode
	Err  string
}
