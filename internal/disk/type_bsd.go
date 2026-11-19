//go:build darwin || freebsd || dragonfly || openbsd || solaris
// +build darwin freebsd dragonfly openbsd solaris

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package disk

// getFSType returns the filesystem type of the underlying mounted filesystem
func getFSType(fstype []int8) string {
	b := make([]byte, len(fstype))
	for i, v := range fstype {
		b[i] = byte(v)
	}
	return string(b)
}
