// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"encoding/base64"
	"net/url"
	"strconv"
)

// Parse bucket url queries
func getListObjectsV1Args(values url.Values) (prefix, marker, delimiter string, maxkeys int, encodingType string, errCode APIErrorCode) {
	errCode = ErrNone

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = ErrInvalidMaxKeys
			return prefix, marker, delimiter, maxkeys, encodingType, errCode
		}
	} else {
		maxkeys = maxObjectList
	}

	prefix = values.Get("prefix")
	marker = values.Get("marker")
	delimiter = values.Get("delimiter")
	encodingType = values.Get("encoding-type")
	return prefix, marker, delimiter, maxkeys, encodingType, errCode
}

func getListBucketObjectVersionsArgs(values url.Values) (prefix, marker, delimiter string, maxkeys int, encodingType, versionIDMarker string, errCode APIErrorCode) {
	errCode = ErrNone

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = ErrInvalidMaxKeys
			return prefix, marker, delimiter, maxkeys, encodingType, versionIDMarker, errCode
		}
	} else {
		maxkeys = maxObjectList
	}

	prefix = values.Get("prefix")
	marker = values.Get("key-marker")
	delimiter = values.Get("delimiter")
	encodingType = values.Get("encoding-type")
	versionIDMarker = values.Get("version-id-marker")
	return prefix, marker, delimiter, maxkeys, encodingType, versionIDMarker, errCode
}

// Parse bucket url queries for ListObjects V2.
func getListObjectsV2Args(values url.Values) (prefix, token, startAfter, delimiter string, fetchOwner bool, maxkeys int, encodingType string, errCode APIErrorCode) {
	errCode = ErrNone

	// The continuation-token cannot be empty.
	if val, ok := values["continuation-token"]; ok {
		if len(val[0]) == 0 {
			errCode = ErrIncorrectContinuationToken
			return prefix, token, startAfter, delimiter, fetchOwner, maxkeys, encodingType, errCode
		}
	}

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = ErrInvalidMaxKeys
			return prefix, token, startAfter, delimiter, fetchOwner, maxkeys, encodingType, errCode
		}
	} else {
		maxkeys = maxObjectList
	}

	prefix = values.Get("prefix")
	startAfter = values.Get("start-after")
	delimiter = values.Get("delimiter")
	fetchOwner = values.Get("fetch-owner") == "true"
	encodingType = values.Get("encoding-type")

	if token = values.Get("continuation-token"); token != "" {
		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			errCode = ErrIncorrectContinuationToken
			return prefix, token, startAfter, delimiter, fetchOwner, maxkeys, encodingType, errCode
		}
		token = string(decodedToken)
	}
	return prefix, token, startAfter, delimiter, fetchOwner, maxkeys, encodingType, errCode
}

// Parse bucket url queries for ?uploads
func getBucketMultipartResources(values url.Values) (prefix, keyMarker, uploadIDMarker, delimiter string, maxUploads int, encodingType string, errCode APIErrorCode) {
	errCode = ErrNone

	if values.Get("max-uploads") != "" {
		var err error
		if maxUploads, err = strconv.Atoi(values.Get("max-uploads")); err != nil {
			errCode = ErrInvalidMaxUploads
			return prefix, keyMarker, uploadIDMarker, delimiter, maxUploads, encodingType, errCode
		}
	} else {
		maxUploads = maxUploadsList
	}

	prefix = values.Get("prefix")
	keyMarker = values.Get("key-marker")
	uploadIDMarker = values.Get("upload-id-marker")
	delimiter = values.Get("delimiter")
	encodingType = values.Get("encoding-type")
	return prefix, keyMarker, uploadIDMarker, delimiter, maxUploads, encodingType, errCode
}

// Parse object url queries
func getObjectResources(values url.Values) (uploadID string, partNumberMarker, maxParts int, encodingType string, errCode APIErrorCode) {
	var err error
	errCode = ErrNone

	if values.Get("max-parts") != "" {
		if maxParts, err = strconv.Atoi(values.Get("max-parts")); err != nil {
			errCode = ErrInvalidMaxParts
			return uploadID, partNumberMarker, maxParts, encodingType, errCode
		}
	} else {
		maxParts = maxPartsList
	}

	if values.Get("part-number-marker") != "" {
		if partNumberMarker, err = strconv.Atoi(values.Get("part-number-marker")); err != nil {
			errCode = ErrInvalidPartNumberMarker
			return uploadID, partNumberMarker, maxParts, encodingType, errCode
		}
	}

	uploadID = values.Get("uploadId")
	encodingType = values.Get("encoding-type")
	return uploadID, partNumberMarker, maxParts, encodingType, errCode
}
