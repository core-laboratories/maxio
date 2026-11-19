// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"

	"github.com/core-laboratories/maxio/internal/logger"
)

const (
	auditFailedMessages    = "failed_messages"
	auditTargetQueueLength = "target_queue_length"
	auditTotalMessages     = "total_messages"
	targetID               = "target_id"
)

var (
	auditFailedMessagesMD = NewCounterMD(auditFailedMessages,
		"Total number of messages that failed to send since start",
		targetID)
	auditTargetQueueLengthMD = NewGaugeMD(auditTargetQueueLength,
		"Number of unsent messages in queue for target",
		targetID)
	auditTotalMessagesMD = NewCounterMD(auditTotalMessages,
		"Total number of messages sent since start",
		targetID)
)

// loadAuditMetrics - `MetricsLoaderFn` for audit
// such as failed messages and total messages.
func loadAuditMetrics(_ context.Context, m MetricValues, c *metricsCache) error {
	audit := logger.CurrentStats()
	for id, st := range audit {
		labels := []string{targetID, id}
		m.Set(auditFailedMessages, float64(st.FailedMessages), labels...)
		m.Set(auditTargetQueueLength, float64(st.QueueLength), labels...)
		m.Set(auditTotalMessages, float64(st.TotalMessages), labels...)
	}

	return nil
}
