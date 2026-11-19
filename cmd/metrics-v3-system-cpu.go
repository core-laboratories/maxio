// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"math"
)

const (
	sysCPUAvgIdle   = "avg_idle"
	sysCPUAvgIOWait = "avg_iowait"
	sysCPULoad      = "load"
	sysCPULoadPerc  = "load_perc"
	sysCPUNice      = "nice"
	sysCPUSteal     = "steal"
	sysCPUSystem    = "system"
	sysCPUUser      = "user"
)

var (
	sysCPUAvgIdleMD   = NewGaugeMD(sysCPUAvgIdle, "Average CPU idle time")
	sysCPUAvgIOWaitMD = NewGaugeMD(sysCPUAvgIOWait, "Average CPU IOWait time")
	sysCPULoadMD      = NewGaugeMD(sysCPULoad, "CPU load average 1min")
	sysCPULoadPercMD  = NewGaugeMD(sysCPULoadPerc, "CPU load average 1min (percentage)")
	sysCPUNiceMD      = NewGaugeMD(sysCPUNice, "CPU nice time")
	sysCPUStealMD     = NewGaugeMD(sysCPUSteal, "CPU steal time")
	sysCPUSystemMD    = NewGaugeMD(sysCPUSystem, "CPU system time")
	sysCPUUserMD      = NewGaugeMD(sysCPUUser, "CPU user time")
)

// loadCPUMetrics - `MetricsLoaderFn` for system CPU metrics.
func loadCPUMetrics(ctx context.Context, m MetricValues, c *metricsCache) error {
	cpuMetrics, _ := c.cpuMetrics.Get()

	if cpuMetrics.LoadStat != nil {
		m.Set(sysCPULoad, cpuMetrics.LoadStat.Load1)
		perc := cpuMetrics.LoadStat.Load1 * 100 / float64(cpuMetrics.CPUCount)
		m.Set(sysCPULoadPerc, math.Round(perc*100)/100)
	}

	ts := cpuMetrics.TimesStat
	if ts != nil {
		tot := ts.User + ts.System + ts.Idle + ts.Iowait + ts.Nice + ts.Steal
		cpuUserVal := math.Round(ts.User/tot*100*100) / 100
		m.Set(sysCPUUser, cpuUserVal)
		cpuSystemVal := math.Round(ts.System/tot*100*100) / 100
		m.Set(sysCPUSystem, cpuSystemVal)
		cpuNiceVal := math.Round(ts.Nice/tot*100*100) / 100
		m.Set(sysCPUNice, cpuNiceVal)
		cpuStealVal := math.Round(ts.Steal/tot*100*100) / 100
		m.Set(sysCPUSteal, cpuStealVal)
	}

	// metrics-resource.go runs a job to collect resource metrics including their Avg values and
	// stores them in resourceMetricsMap. We can use it to get the Avg values of CPU idle and IOWait.
	cpuResourceMetrics, found := resourceMetricsMap[cpuSubsystem]
	if found {
		if cpuIdleMetric, ok := cpuResourceMetrics[getResourceKey(cpuIdle, nil)]; ok {
			avgVal := math.Round(cpuIdleMetric.Avg*100) / 100
			m.Set(sysCPUAvgIdle, avgVal)
		}
		if cpuIOWaitMetric, ok := cpuResourceMetrics[getResourceKey(cpuIOWait, nil)]; ok {
			avgVal := math.Round(cpuIOWaitMetric.Avg*100) / 100
			m.Set(sysCPUAvgIOWait, avgVal)
		}
	}
	return nil
}
