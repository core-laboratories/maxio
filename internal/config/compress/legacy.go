// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package compress

import (
	"strings"

	"github.com/core-laboratories/maxio/internal/config"
)

// Legacy envs.
const (
	EnvCompress                 = "MINIO_COMPRESS"
	EnvCompressMimeTypesLegacy1 = "MINIO_COMPRESS_MIMETYPES"

	// These envs were wrong but we supported them for a long time
	// so keep them here to support existing deployments.
	EnvCompressEnableLegacy          = "MINIO_COMPRESS_ENABLE"
	EnvCompressAllowEncryptionLegacy = "MINIO_COMPRESS_ALLOW_ENCRYPTION"
	EnvCompressExtensionsLegacy      = "MINIO_COMPRESS_EXTENSIONS"
	EnvCompressMimeTypesLegacy2      = "MINIO_COMPRESS_MIME_TYPES"
)

// SetCompressionConfig - One time migration code needed, for migrating from older config to new for Compression.
func SetCompressionConfig(s config.Config, cfg Config) {
	if !cfg.Enabled {
		// No need to save disabled settings in new config.
		return
	}
	s[config.CompressionSubSys][config.Default] = config.KVS{
		config.KV{
			Key:   config.Enable,
			Value: config.EnableOn,
		},
		config.KV{
			Key:   Extensions,
			Value: strings.Join(cfg.Extensions, config.ValueSeparator),
		},
		config.KV{
			Key:   MimeTypes,
			Value: strings.Join(cfg.MimeTypes, config.ValueSeparator),
		},
	}
}
