// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/minio/kms-go/kes"
	"github.com/minio/madmin-go/v3"
	"github.com/core-laboratories/maxio/internal/kms"
	"github.com/core-laboratories/maxio/internal/logger"
	"github.com/minio/mux"
	"github.com/minio/pkg/v3/policy"
)

const (
	// Bucket Encryption configuration file name.
	bucketSSEConfig = "bucket-encryption.xml"
)

// PutBucketEncryptionHandler - Stores given bucket encryption configuration
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketEncryption.html
func (api objectAPIHandlers) PutBucketEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r, w, "PutBucketEncryption")

	defer logger.AuditLog(ctx, w, r, mustGetClaimsFromToken(r))

	objAPI := api.ObjectAPI()
	if objAPI == nil {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrServerNotInitialized), r.URL)
		return
	}

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	if s3Error := checkRequestAuthType(ctx, r, policy.PutBucketEncryptionAction, bucket, ""); s3Error != ErrNone {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(s3Error), r.URL)
		return
	}

	// Check if bucket exists.
	if _, err := objAPI.GetBucketInfo(ctx, bucket, BucketOptions{}); err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	// Parse bucket encryption xml
	encConfig, err := validateBucketSSEConfig(io.LimitReader(r.Body, maxBucketSSEConfigSize))
	if err != nil {
		apiErr := APIError{
			Code:           "MalformedXML",
			Description:    fmt.Sprintf("%s (%s)", errorCodes[ErrMalformedXML].Description, err),
			HTTPStatusCode: errorCodes[ErrMalformedXML].HTTPStatusCode,
		}
		writeErrorResponse(ctx, w, apiErr, r.URL)
		return
	}

	// Return error if KMS is not initialized
	if GlobalKMS == nil {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrKMSNotConfigured), r.URL)
		return
	}
	kmsKey := encConfig.KeyID()
	if kmsKey != "" {
		kmsContext := kms.Context{"MinIO admin API": "ServerInfoHandler"} // Context for a test key operation
		_, err := GlobalKMS.GenerateKey(ctx, &kms.GenerateKeyRequest{Name: kmsKey, AssociatedData: kmsContext})
		if err != nil {
			if errors.Is(err, kes.ErrKeyNotFound) {
				writeErrorResponse(ctx, w, toAPIError(ctx, errKMSKeyNotFound), r.URL)
				return
			}
			writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
			return
		}
	}

	configData, err := xml.Marshal(encConfig)
	if err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	// Store the bucket encryption configuration in the object layer
	updatedAt, err := globalBucketMetadataSys.Update(ctx, bucket, bucketSSEConfig, configData)
	if err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	// Call site replication hook.
	//
	// We encode the xml bytes as base64 to ensure there are no encoding
	// errors.
	cfgStr := base64.StdEncoding.EncodeToString(configData)
	replLogIf(ctx, globalSiteReplicationSys.BucketMetaHook(ctx, madmin.SRBucketMeta{
		Type:      madmin.SRBucketMetaTypeSSEConfig,
		Bucket:    bucket,
		SSEConfig: &cfgStr,
		UpdatedAt: updatedAt,
	}))

	writeSuccessResponseHeadersOnly(w)
}

// GetBucketEncryptionHandler - Returns bucket policy configuration
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketEncryption.html
func (api objectAPIHandlers) GetBucketEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r, w, "GetBucketEncryption")

	defer logger.AuditLog(ctx, w, r, mustGetClaimsFromToken(r))

	objAPI := api.ObjectAPI()
	if objAPI == nil {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrServerNotInitialized), r.URL)
		return
	}

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	if s3Error := checkRequestAuthType(ctx, r, policy.GetBucketEncryptionAction, bucket, ""); s3Error != ErrNone {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(s3Error), r.URL)
		return
	}

	// Check if bucket exists
	var err error
	if _, err = objAPI.GetBucketInfo(ctx, bucket, BucketOptions{}); err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	config, _, err := globalBucketMetadataSys.GetSSEConfig(bucket)
	if err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	configData, err := xml.Marshal(config)
	if err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	// Write bucket encryption configuration to client
	writeSuccessResponseXML(w, configData)
}

// DeleteBucketEncryptionHandler - Removes bucket encryption configuration
func (api objectAPIHandlers) DeleteBucketEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r, w, "DeleteBucketEncryption")

	defer logger.AuditLog(ctx, w, r, mustGetClaimsFromToken(r))

	objAPI := api.ObjectAPI()
	if objAPI == nil {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrServerNotInitialized), r.URL)
		return
	}

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	if s3Error := checkRequestAuthType(ctx, r, policy.PutBucketEncryptionAction, bucket, ""); s3Error != ErrNone {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(s3Error), r.URL)
		return
	}

	// Check if bucket exists
	var err error
	if _, err = objAPI.GetBucketInfo(ctx, bucket, BucketOptions{}); err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	// Delete bucket encryption config from object layer
	updatedAt, err := globalBucketMetadataSys.Delete(ctx, bucket, bucketSSEConfig)
	if err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL)
		return
	}

	// Call site replication hook.
	replLogIf(ctx, globalSiteReplicationSys.BucketMetaHook(ctx, madmin.SRBucketMeta{
		Type:      madmin.SRBucketMetaTypeSSEConfig,
		Bucket:    bucket,
		SSEConfig: nil,
		UpdatedAt: updatedAt,
	}))

	writeSuccessNoContent(w)
}
