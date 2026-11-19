// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"net/http"

	"github.com/minio/mux"
)

const (
	healthCheckPath            = "/health"
	healthCheckLivenessPath    = "/live"
	healthCheckReadinessPath   = "/ready"
	healthCheckClusterPath     = "/cluster"
	healthCheckClusterReadPath = "/cluster/read"
	healthCheckPathPrefix      = minioReservedBucketPath + healthCheckPath
)

// registerHealthCheckRouter - add handler functions for liveness and readiness routes.
func registerHealthCheckRouter(router *mux.Router) {
	// Healthcheck router
	healthRouter := router.PathPrefix(healthCheckPathPrefix).Subrouter()

	// Cluster check handler to verify cluster is active
	healthRouter.Methods(http.MethodGet).Path(healthCheckClusterPath).HandlerFunc(httpTraceAll(ClusterCheckHandler))
	healthRouter.Methods(http.MethodHead).Path(healthCheckClusterPath).HandlerFunc(httpTraceAll(ClusterCheckHandler))
	healthRouter.Methods(http.MethodGet).Path(healthCheckClusterReadPath).HandlerFunc(httpTraceAll(ClusterReadCheckHandler))
	healthRouter.Methods(http.MethodHead).Path(healthCheckClusterReadPath).HandlerFunc(httpTraceAll(ClusterReadCheckHandler))

	// Liveness handler
	healthRouter.Methods(http.MethodGet).Path(healthCheckLivenessPath).HandlerFunc(httpTraceAll(LivenessCheckHandler))
	healthRouter.Methods(http.MethodHead).Path(healthCheckLivenessPath).HandlerFunc(httpTraceAll(LivenessCheckHandler))

	// Readiness handler
	healthRouter.Methods(http.MethodGet).Path(healthCheckReadinessPath).HandlerFunc(httpTraceAll(ReadinessCheckHandler))
	healthRouter.Methods(http.MethodHead).Path(healthCheckReadinessPath).HandlerFunc(httpTraceAll(ReadinessCheckHandler))
}
