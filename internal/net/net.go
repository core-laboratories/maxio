// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package net

import (
	"fmt"

	"github.com/prometheus/procfs"
)

// GetInterfaceNetStats - get procfs.NetDevLine by interfaceName
func GetInterfaceNetStats(interf string) (procfs.NetDevLine, error) {
	proc, err := procfs.Self()
	if err != nil {
		return procfs.NetDevLine{}, err
	}
	netDev, err := proc.NetDev()
	if err != nil {
		return procfs.NetDevLine{}, err
	}
	ndl, ok := netDev[interf]
	if !ok {
		return procfs.NetDevLine{}, fmt.Errorf("%v interface not found", interf)
	}
	return ndl, nil
}
