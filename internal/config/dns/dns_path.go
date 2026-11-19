// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package dns

import (
	"path"
	"strings"

	"github.com/miekg/dns"
)

// msgPath converts a domainname to an etcd path. If s looks like service.staging.skydns.local.,
// the resulting key will be /skydns/local/skydns/staging/service .
func msgPath(s, prefix string) string {
	l := dns.SplitDomainName(s)
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	return path.Join(append([]string{etcdPathSeparator + prefix + etcdPathSeparator}, l...)...)
}

// dnsJoin joins labels to form a fully qualified domain name. If the last label is
// the root label it is ignored. Not other syntax checks are performed.
func dnsJoin(labels ...string) string {
	if len(labels) == 0 {
		return ""
	}
	ll := len(labels)
	if labels[ll-1] == "." {
		return strings.Join(labels[:ll-1], ".") + "."
	}
	return dns.Fqdn(strings.Join(labels, "."))
}

// msgUnPath converts a etcd path to domainName.
func msgUnPath(s string) string {
	ks := strings.Split(strings.Trim(s, etcdPathSeparator), etcdPathSeparator)
	for i, j := 0, len(ks)-1; i < j; i, j = i+1, j-1 {
		ks[i], ks[j] = ks[j], ks[i]
	}
	return strings.Join(ks, ".")
}
