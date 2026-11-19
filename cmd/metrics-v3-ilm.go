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
	expiryPendingTasks             = "expiry_pending_tasks"
	transitionActiveTasks          = "transition_active_tasks"
	transitionPendingTasks         = "transition_pending_tasks"
	transitionMissedImmediateTasks = "transition_missed_immediate_tasks"
	versionsScanned                = "versions_scanned"
)

var (
	ilmExpiryPendingTasksMD             = NewGaugeMD(expiryPendingTasks, "Number of pending ILM expiry tasks in the queue")
	ilmTransitionActiveTasksMD          = NewGaugeMD(transitionActiveTasks, "Number of active ILM transition tasks")
	ilmTransitionPendingTasksMD         = NewGaugeMD(transitionPendingTasks, "Number of pending ILM transition tasks in the queue")
	ilmTransitionMissedImmediateTasksMD = NewCounterMD(transitionMissedImmediateTasks, "Number of missed immediate ILM transition tasks")
	ilmVersionsScannedMD                = NewCounterMD(versionsScanned, "Total number of object versions checked for ILM actions since server start")
)

// loadILMMetrics - `MetricsLoaderFn` for ILM metrics.
func loadILMMetrics(_ context.Context, m MetricValues, _ *metricsCache) error {
	if globalExpiryState != nil {
		m.Set(expiryPendingTasks, float64(globalExpiryState.PendingTasks()))
	}
	if globalTransitionState != nil {
		m.Set(transitionActiveTasks, float64(globalTransitionState.ActiveTasks()))
		m.Set(transitionPendingTasks, float64(globalTransitionState.PendingTasks()))
		m.Set(transitionMissedImmediateTasks, float64(globalTransitionState.MissedImmediateTasks()))
	}
	m.Set(versionsScanned, float64(globalScannerMetrics.lifetime(scannerMetricILM)))

	return nil
}
