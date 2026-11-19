// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package event

import "net/http"

// Identity represents access key who caused the event.
type Identity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	AccessKeyID string `json:"accessKeyId"`
}

// UserRequest user request headers
type UserRequest struct {
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
}

// GetObjectContext provides the necessary details to perform
// download of the object, and return back the processed response
// to the server.
type GetObjectContext struct {
	OutputRoute string `json:"outputRoute"`
	OutputToken string `json:"outputToken"`
	InputS3URL  string `json:"inputS3Url"`
}

// Event represents lambda function event, this is undocumented in AWS S3. This
// structure bases itself on this structure but there is no binding.
//
//	{
//	  "xAmzRequestId": "a2871150-1df5-4dc9-ad9f-3da283ca1bf3",
//	  "getObjectContext": {
//	    "outputRoute": "...",
//	    "outputToken": "...",
//	    "inputS3Url": "<presignedURL>"
//	  },
//	  "configuration": { // not useful in MinIO
//	    "accessPointArn": "...",
//	    "supportingAccessPointArn": "...",
//	    "payload": ""
//	  },
//	  "userRequest": {
//	    "url": "...",
//	    "headers": {
//	      "Host": "...",
//	      "X-Amz-Content-SHA256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
//	    }
//	  },
//	  "userIdentity": {
//	    "type": "IAMUser",
//	    "principalId": "AIDAJF5MO57RFXQCE5ZNC",
//	    "arn": "...",
//	    "accountId": "...",
//	    "accessKeyId": "AKIA3WNQJCXE2DYPAU7R"
//	  },
//	  "protocolVersion": "1.00"
//	}
type Event struct {
	ProtocolVersion  string            `json:"protocolVersion"`
	GetObjectContext *GetObjectContext `json:"getObjectContext"`
	UserIdentity     Identity          `json:"userIdentity"`
	UserRequest      UserRequest       `json:"userRequest"`
}
