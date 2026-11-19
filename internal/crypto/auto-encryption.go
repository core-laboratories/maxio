// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package crypto

import (
	"github.com/core-laboratories/maxio/internal/config"
	"github.com/minio/pkg/v3/env"
)

const (
	// EnvKMSAutoEncryption is the environment variable used to en/disable
	// SSE-S3 auto-encryption. SSE-S3 auto-encryption, if enabled,
	// requires a valid KMS configuration and turns any non-SSE-C
	// request into an SSE-S3 request.
	// If present EnvAutoEncryption must be either "on" or "off".
	EnvKMSAutoEncryption = "MINIO_KMS_AUTO_ENCRYPTION"
)

// LookupAutoEncryption returns true if and only if
// the MINIO_KMS_AUTO_ENCRYPTION env. variable is
// set to "on".
func LookupAutoEncryption() bool {
	auto, _ := config.ParseBool(env.Get(EnvKMSAutoEncryption, config.EnableOff))
	return auto
}
