// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package init

import (
	"os"

	"github.com/klauspost/cpuid/v2"
)

func init() {
	// All MinIO operations must be under UTC.
	os.Setenv("TZ", "UTC")

	// Temporary workaround for
	// https://github.com/golang/go/issues/49233
	// Keep until upstream has been fixed.
	cpuid.CPU.Disable(cpuid.AVX512F, cpuid.AVX512BW, cpuid.AVX512CD, cpuid.AVX512DQ,
		cpuid.AVX512ER, cpuid.AVX512FP16, cpuid.AVX512IFMA, cpuid.AVX512PF, cpuid.AVX512VBMI,
		cpuid.AVX512VBMI2, cpuid.AVX512VL, cpuid.AVX512VNNI, cpuid.AVX512VP2INTERSECT, cpuid.AVX512VPOPCNTDQ)
}
