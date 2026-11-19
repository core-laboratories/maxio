// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package grid

//go:generate msgp -unexported -file=$GOFILE -tests=false -o=grid_types_msgp_test.go

type testRequest struct {
	Num    int
	String string
}

type testResponse struct {
	OrgNum    int
	OrgString string
	Embedded  testRequest
}

func newTestRequest() *testRequest {
	return &testRequest{}
}

func newTestResponse() *testResponse {
	return &testResponse{}
}
