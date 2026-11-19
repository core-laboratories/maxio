// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ipfs

import (
	"strings"

	"github.com/google/uuid"
)

// TenantID is a stable server-assigned identifier derived from the IAM principal.
// It is an access-control identifier, not a secret or authentication token.
func TenantID(principal string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte("maxio-tenant:"+principal)).String()
}

func TenantPrefix(principal string) string {
	return "users/" + TenantID(principal) + "/"
}

func AuthorizedKey(principal, key string) bool {
	return principal != "" && strings.HasPrefix(key, TenantPrefix(principal))
}
