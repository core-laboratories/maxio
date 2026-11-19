// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"strings"

	"github.com/core-laboratories/maxio/internal/bucket/versioning"
	"github.com/core-laboratories/maxio/internal/logger"
)

// BucketVersioningSys - policy subsystem.
type BucketVersioningSys struct{}

// Enabled enabled versioning?
func (sys *BucketVersioningSys) Enabled(bucket string) bool {
	vc, err := sys.Get(bucket)
	if err != nil {
		logger.CriticalIf(GlobalContext, err)
	}
	return vc.Enabled()
}

// PrefixEnabled returns true is versioning is enabled at bucket level and if
// the given prefix doesn't match any excluded prefixes pattern. This is
// part of a MinIO versioning configuration extension.
func (sys *BucketVersioningSys) PrefixEnabled(bucket, prefix string) bool {
	vc, err := sys.Get(bucket)
	if err != nil {
		logger.CriticalIf(GlobalContext, err)
	}
	return vc.PrefixEnabled(prefix)
}

// Suspended suspended versioning?
func (sys *BucketVersioningSys) Suspended(bucket string) bool {
	vc, err := sys.Get(bucket)
	if err != nil {
		logger.CriticalIf(GlobalContext, err)
	}
	return vc.Suspended()
}

// PrefixSuspended returns true if the given prefix matches an excluded prefix
// pattern. This is part of a MinIO versioning configuration extension.
func (sys *BucketVersioningSys) PrefixSuspended(bucket, prefix string) bool {
	vc, err := sys.Get(bucket)
	if err != nil {
		logger.CriticalIf(GlobalContext, err)
	}

	return vc.PrefixSuspended(prefix)
}

// Get returns stored bucket policy
func (sys *BucketVersioningSys) Get(bucket string) (*versioning.Versioning, error) {
	if bucket == minioMetaBucket || strings.HasPrefix(bucket, minioMetaBucket) {
		return &versioning.Versioning{XMLNS: "http://s3.amazonaws.com/doc/2006-03-01/"}, nil
	}

	vcfg, _, err := globalBucketMetadataSys.GetVersioningConfig(bucket)
	return vcfg, err
}

// NewBucketVersioningSys - creates new versioning system.
func NewBucketVersioningSys() *BucketVersioningSys {
	return &BucketVersioningSys{}
}
