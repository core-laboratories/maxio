// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import "context"

const (
	configRRSParity      = "rrs_parity"
	configStandardParity = "standard_parity"
)

var (
	configRRSParityMD = NewGaugeMD(configRRSParity,
		"Reduced redundancy storage class parity")
	configStandardParityMD = NewGaugeMD(configStandardParity,
		"Standard storage class parity")
)

// loadClusterConfigMetrics - `MetricsLoaderFn` for cluster config
// such as standard and RRS parity.
func loadClusterConfigMetrics(ctx context.Context, m MetricValues, c *metricsCache) error {
	clusterDriveMetrics, err := c.clusterDriveMetrics.Get()
	if err != nil {
		metricsLogIf(ctx, err)
	} else {
		m.Set(configStandardParity, float64(clusterDriveMetrics.storageInfo.Backend.StandardSCParity))
		m.Set(configRRSParity, float64(clusterDriveMetrics.storageInfo.Backend.RRSCParity))
	}

	return nil
}
