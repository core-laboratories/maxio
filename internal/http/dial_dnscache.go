// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package http

import (
	"context"
	"net"
	"time"
)

// LookupHost is a function to make custom lookupHost for optional cached DNS requests
type LookupHost func(ctx context.Context, host string) (addrs []string, err error)

// DialContextWithLookupHost is a helper function which returns `net.DialContext` function.
// It randomly fetches an IP via custom LookupHost function and dials it by the given dial
// function. LookupHost may implement an internal DNS caching implementation, lookupHost
// input if nil then net.DefaultResolver.LookupHost is used.
//
// It dials one by one and returns first connected `net.Conn`.
// If it fails to dial all IPs from cache it returns first error. If no baseDialFunc
// is given, it sets default dial function.
//
// You can use returned dial function for `http.Transport.DialContext`.
//
// In this function, it uses functions from `rand` package. To make it really random,
// you MUST call `rand.Seed` and change the value from the default in your application
func DialContextWithLookupHost(lookupHost LookupHost, baseDialCtx DialContext) DialContext {
	if lookupHost == nil {
		lookupHost = net.DefaultResolver.LookupHost
	}

	if baseDialCtx == nil {
		// This is same as which `http.DefaultTransport` uses.
		baseDialCtx = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	return func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		if net.ParseIP(host) != nil {
			// For IP only setups there is no need for DNS lookups.
			return baseDialCtx(ctx, "tcp", addr)
		}

		ips, err := lookupHost(ctx, host)
		if err != nil {
			return nil, err
		}

		for _, ip := range ips {
			conn, err = baseDialCtx(ctx, "tcp", net.JoinHostPort(ip, port))
			if err == nil {
				break
			}
		}

		return conn, err
	}
}
