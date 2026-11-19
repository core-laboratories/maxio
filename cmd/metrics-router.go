// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"net/http"
	"strings"

	"github.com/minio/mux"
	"github.com/minio/pkg/v3/env"
)

const (
	prometheusMetricsPathLegacy     = "/prometheus/metrics"
	prometheusMetricsV2ClusterPath  = "/v2/metrics/cluster"
	prometheusMetricsV2BucketPath   = "/v2/metrics/bucket"
	prometheusMetricsV2NodePath     = "/v2/metrics/node"
	prometheusMetricsV2ResourcePath = "/v2/metrics/resource"

	// Metrics v3 endpoints
	metricsV3Path = "/metrics/v3"
)

// Standard env prometheus auth type
const (
	EnvPrometheusAuthType    = "MINIO_PROMETHEUS_AUTH_TYPE"
	EnvPrometheusOpenMetrics = "MINIO_PROMETHEUS_OPEN_METRICS"
)

type prometheusAuthType string

const (
	prometheusJWT    prometheusAuthType = "jwt"
	prometheusPublic prometheusAuthType = "public"
)

// registerMetricsRouter - add handler functions for metrics.
func registerMetricsRouter(router *mux.Router) {
	// metrics router
	metricsRouter := router.NewRoute().PathPrefix(minioReservedBucketPath).Subrouter()
	authType := prometheusAuthType(strings.ToLower(env.Get(EnvPrometheusAuthType, string(prometheusJWT))))

	auth := AuthMiddleware
	if authType == prometheusPublic {
		auth = NoAuthMiddleware
	}

	metricsRouter.Handle(prometheusMetricsPathLegacy, auth(metricsHandler()))
	metricsRouter.Handle(prometheusMetricsV2ClusterPath, auth(metricsServerHandler()))
	metricsRouter.Handle(prometheusMetricsV2BucketPath, auth(metricsBucketHandler()))
	metricsRouter.Handle(prometheusMetricsV2NodePath, auth(metricsNodeHandler()))
	metricsRouter.Handle(prometheusMetricsV2ResourcePath, auth(metricsResourceHandler()))

	// Metrics v3
	metricsV3Server := newMetricsV3Server(auth)

	// Register metrics v3 handler. It also accepts an optional query
	// parameter `?list` - see handler for details.
	metricsRouter.Methods(http.MethodGet).Path(metricsV3Path + "{pathComps:.*}").Handler(metricsV3Server)
}
