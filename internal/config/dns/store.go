// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package dns

// Error - DNS related errors error.
type Error struct {
	Bucket string
	Err    error
}

// ErrInvalidBucketName for buckets with invalid name
type ErrInvalidBucketName Error

func (e ErrInvalidBucketName) Error() string {
	return e.Bucket + " invalid bucket name error: " + e.Err.Error()
}

func (e Error) Error() string {
	return "dns related error: " + e.Err.Error()
}

// ErrBucketConflict for buckets that already exist
type ErrBucketConflict Error

func (e ErrBucketConflict) Error() string {
	return e.Bucket + " bucket conflict error: " + e.Err.Error()
}

// Store dns record store
type Store interface {
	Put(bucket string) error
	Get(bucket string) ([]SrvRecord, error)
	Delete(bucket string) error
	List() (map[string][]SrvRecord, error)
	DeleteRecord(record SrvRecord) error
	Close() error
	String() string
}
