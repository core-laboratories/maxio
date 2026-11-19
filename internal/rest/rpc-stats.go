// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package rest

import (
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"
)

var globalStats = struct {
	errs uint64

	tcpDialErrs                 uint64
	tcpDialCount                uint64
	tcpDialTotalDur             uint64
	tcpTimeForFirstByteTotalDur uint64
}{}

// RPCStats holds information about the DHCP/TCP metrics and errors
type RPCStats struct {
	Errs uint64

	DialAvgDuration uint64
	TTFBAvgDuration uint64
	DialErrs        uint64
}

// GetRPCStats returns RPC stats, include calls errors and dhcp/tcp metrics
func GetRPCStats() RPCStats {
	s := RPCStats{
		Errs:     atomic.LoadUint64(&globalStats.errs),
		DialErrs: atomic.LoadUint64(&globalStats.tcpDialErrs),
	}
	if v := atomic.LoadUint64(&globalStats.tcpDialCount); v > 0 {
		s.DialAvgDuration = atomic.LoadUint64(&globalStats.tcpDialTotalDur) / v
		s.TTFBAvgDuration = atomic.LoadUint64(&globalStats.tcpTimeForFirstByteTotalDur) / v
	}
	return s
}

// Return a function which update the global stats related to tcp connections
func setupReqStatsUpdate(req *http.Request) (*http.Request, func()) {
	var dialStart, dialEnd int64
	start := time.Now()
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			atomic.AddUint64(&globalStats.tcpTimeForFirstByteTotalDur, uint64(time.Since(start)))
		},
		ConnectStart: func(network, addr string) {
			atomic.StoreInt64(&dialStart, time.Now().UnixNano())
		},
		ConnectDone: func(network, addr string, err error) {
			if err == nil {
				atomic.StoreInt64(&dialEnd, time.Now().UnixNano())
			}
		},
	}

	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace)), func() {
		if ds := atomic.LoadInt64(&dialStart); ds > 0 {
			if de := atomic.LoadInt64(&dialEnd); de == 0 {
				atomic.AddUint64(&globalStats.tcpDialErrs, 1)
			} else if de >= ds {
				atomic.AddUint64(&globalStats.tcpDialCount, 1)
				atomic.AddUint64(&globalStats.tcpDialTotalDur, uint64(dialEnd-dialStart))
			}
		}
	}
}
