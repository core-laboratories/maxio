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
	webhookQueueLength    = "queue_length"
	webhookTotalMessages  = "total_messages"
	webhookFailedMessages = "failed_messages"
	nameL                 = "name"
	endpointL             = "endpoint"
)

var (
	allWebhookLabels        = []string{nameL, endpointL}
	webhookFailedMessagesMD = NewCounterMD(webhookFailedMessages,
		"Number of messages that failed to send",
		allWebhookLabels...)
	webhookQueueLengthMD = NewGaugeMD(webhookQueueLength,
		"Webhook queue length",
		allWebhookLabels...)
	webhookTotalMessagesMD = NewCounterMD(webhookTotalMessages,
		"Total number of messages sent to this target",
		allWebhookLabels...)
)

// loadLoggerWebhookMetrics - `MetricsLoaderFn` for logger webhook
// such as failed messages and total messages.
func loadLoggerWebhookMetrics(ctx context.Context, m MetricValues, c *metricsCache) error {
	tgts := append(logger.SystemTargets(), logger.AuditTargets()...)
	for _, t := range tgts {
		labels := []string{nameL, t.String(), endpointL, t.Endpoint()}
		m.Set(webhookFailedMessages, float64(t.Stats().FailedMessages), labels...)
		m.Set(webhookQueueLength, float64(t.Stats().QueueLength), labels...)
		m.Set(webhookTotalMessages, float64(t.Stats().TotalMessages), labels...)
	}

	return nil
}
