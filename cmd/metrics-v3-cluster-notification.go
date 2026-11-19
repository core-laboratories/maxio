// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
)

const (
	notificationCurrentSendInProgress = "current_send_in_progress"
	notificationEventsErrorsTotal     = "events_errors_total"
	notificationEventsSentTotal       = "events_sent_total"
	notificationEventsSkippedTotal    = "events_skipped_total"
)

var (
	notificationCurrentSendInProgressMD = NewCounterMD(notificationCurrentSendInProgress, "Number of concurrent async Send calls active to all targets")
	notificationEventsErrorsTotalMD     = NewCounterMD(notificationEventsErrorsTotal, "Events that were failed to be sent to the targets")
	notificationEventsSentTotalMD       = NewCounterMD(notificationEventsSentTotal, "Total number of events sent to the targets")
	notificationEventsSkippedTotalMD    = NewCounterMD(notificationEventsSkippedTotal, "Events that were skipped to be sent to the targets due to the in-memory queue being full")
)

// loadClusterNotificationMetrics - `MetricsLoaderFn` for cluster notification metrics.
func loadClusterNotificationMetrics(_ context.Context, m MetricValues, _ *metricsCache) error {
	if globalEventNotifier == nil {
		return nil
	}

	nstats := globalEventNotifier.targetList.Stats()
	m.Set(notificationCurrentSendInProgress, float64(nstats.CurrentSendCalls))
	m.Set(notificationEventsErrorsTotal, float64(nstats.EventsErrorsTotal))
	m.Set(notificationEventsSentTotal, float64(nstats.TotalEvents))
	m.Set(notificationEventsSkippedTotal, float64(nstats.EventsSkipped))

	return nil
}
