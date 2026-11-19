// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import "errors"

// errErasureReadQuorum - did not meet read quorum.
var errErasureReadQuorum = errors.New("Read failed. Insufficient number of drives online")

// errErasureWriteQuorum - did not meet write quorum.
var errErasureWriteQuorum = errors.New("Write failed. Insufficient number of drives online")

// errNoHealRequired - returned when healing is attempted on a previously healed disks.
var errNoHealRequired = errors.New("No healing is required")
