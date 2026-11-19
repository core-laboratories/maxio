// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"sync"

	"github.com/minio/madmin-go/v3"
	"github.com/core-laboratories/maxio/internal/pubsub"
)

const bootstrapTraceLimit = 4 << 10

type bootstrapTracer struct {
	mu   sync.RWMutex
	info []madmin.TraceInfo
}

var globalBootstrapTracer = &bootstrapTracer{}

func (bs *bootstrapTracer) Record(info madmin.TraceInfo) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if len(bs.info) > bootstrapTraceLimit {
		return
	}
	bs.info = append(bs.info, info)
}

func (bs *bootstrapTracer) Events() []madmin.TraceInfo {
	traceInfo := make([]madmin.TraceInfo, 0, bootstrapTraceLimit)

	bs.mu.RLock()
	traceInfo = append(traceInfo, bs.info...)
	bs.mu.RUnlock()

	return traceInfo
}

func (bs *bootstrapTracer) Publish(ctx context.Context, trace *pubsub.PubSub[madmin.TraceInfo, madmin.TraceType]) {
	for _, bsEvent := range bs.Events() {
		if bsEvent.Message != "" {
			select {
			case <-ctx.Done():
			default:
				trace.Publish(bsEvent)
			}
		}
	}
}
