// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package dns

import "testing"

func TestDNSJoin(t *testing.T) {
	tests := []struct {
		in  []string
		out string
	}{
		{[]string{"bla", "bliep", "example", "org"}, "bla.bliep.example.org."},
		{[]string{"example", "."}, "example."},
		{[]string{"example", "org."}, "example.org."}, // technically we should not be called like this.
		{[]string{"."}, "."},
	}

	for i, tc := range tests {
		if x := dnsJoin(tc.in...); x != tc.out {
			t.Errorf("Test %d, expected %s, got %s", i, tc.out, x)
		}
	}
}

func TestPath(t *testing.T) {
	for _, path := range []string{"mydns", "skydns"} {
		result := msgPath("service.staging.skydns.local.", path)
		if result != etcdPathSeparator+path+"/local/skydns/staging/service" {
			t.Errorf("Failure to get domain's path with prefix: %s", result)
		}
	}
}

func TestUnPath(t *testing.T) {
	result1 := msgUnPath("/skydns/local/cluster/staging/service/")
	if result1 != "service.staging.cluster.local.skydns" {
		t.Errorf("Failure to get domain from etcd key (with a trailing '/'), expect: 'service.staging.cluster.local.', actually get: '%s'", result1)
	}

	result2 := msgUnPath("/skydns/local/cluster/staging/service")
	if result2 != "service.staging.cluster.local.skydns" {
		t.Errorf("Failure to get domain from etcd key (without trailing '/'), expect: 'service.staging.cluster.local.' actually get: '%s'", result2)
	}

	result3 := msgUnPath("/singleleveldomain/")
	if result3 != "singleleveldomain" {
		t.Errorf("Failure to get domain from etcd key (with leading and trailing '/'), expect: 'singleleveldomain.' actually get: '%s'", result3)
	}

	result4 := msgUnPath("/singleleveldomain")
	if result4 != "singleleveldomain" {
		t.Errorf("Failure to get domain from etcd key (without trailing '/'), expect: 'singleleveldomain.' actually get: '%s'", result4)
	}

	result5 := msgUnPath("singleleveldomain")
	if result5 != "singleleveldomain" {
		t.Errorf("Failure to get domain from etcd key (without leading and trailing '/'), expect: 'singleleveldomain.' actually get: '%s'", result5)
	}
}
