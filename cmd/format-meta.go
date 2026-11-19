// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

// Format related consts
const (
	// Format config file carries backend format specific details.
	formatConfigFile = "format.json"
)

const (
	// Version of the formatMetaV1
	formatMetaVersionV1 = "1"
)

// format.json currently has the format:
// {
//   "version": "1",
//   "format": "XXXXX",
//   "XXXXX": {
//
//   }
// }
// Here "XXXXX" depends on the backend, currently we have "fs" and "xl" implementations.
// formatMetaV1 should be inherited by backend format structs. Please look at format-fs.go
// and format-xl.go for details.

// Ideally we will never have a situation where we will have to change the
// fields of this struct and deal with related migration.
type formatMetaV1 struct {
	// Version of the format config.
	Version string `json:"version"`
	// Format indicates the backend format type, supports two values 'xl' and 'fs'.
	Format string `json:"format"`
	// ID is the identifier for the minio deployment
	ID string `json:"id"`
}
