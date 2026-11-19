// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package dsync

// Dsync represents dsync client object which is initialized with
// authenticated clients, used to initiate lock REST calls.
type Dsync struct {
	// List of rest client objects, one per lock server.
	GetLockers func() ([]NetLocker, string)

	// Timeouts to apply.
	Timeouts Timeouts
}
