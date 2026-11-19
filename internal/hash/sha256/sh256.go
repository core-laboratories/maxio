// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package sha256

import (
	"crypto/sha256"
	"hash"
)

// New initializes a new sha256.New()
func New() hash.Hash { return sha256.New() }

// Sum256 returns the SHA256 checksum of the data.
func Sum256(data []byte) [sha256.Size]byte { return sha256.Sum256(data) }

// Size is the size of a SHA256 checksum in bytes.
const Size = sha256.Size
