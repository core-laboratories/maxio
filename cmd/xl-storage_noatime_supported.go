//go:build unix && !darwin && !freebsd

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"os"
	"syscall"
)

var (
	// Disallow updating access times
	// Add non-block to avoid syscall to attempt to set epoll on files.
	readMode = os.O_RDONLY | 0x40000 | syscall.O_NONBLOCK // O_NOATIME

	// Write with data sync only used only for `xl.meta` writes
	writeMode = 0x1000 | syscall.O_NONBLOCK // O_DSYNC
)
