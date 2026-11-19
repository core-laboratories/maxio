//go:build !unix || darwin || freebsd

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"os"
)

var (
	// No special option for reads on windows
	readMode = os.O_RDONLY

	// Write with sync no buffering only used only for `xl.meta` writes
	writeMode = os.O_SYNC
)
