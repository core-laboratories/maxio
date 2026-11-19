// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/core-laboratories/maxio/internal/bucket/lifecycle"
	xhttp "github.com/core-laboratories/maxio/internal/http"
)

// TestParseRestoreObjStatus tests parseRestoreObjStatus
func TestParseRestoreObjStatus(t *testing.T) {
	testCases := []struct {
		restoreHdr     string
		expectedStatus restoreObjStatus
		expectedErr    error
	}{
		{
			// valid: represents a restored object, 'pending' expiry.
			restoreHdr: `ongoing-request="false", expiry-date="Fri, 21 Dec 2012 00:00:00 GMT"`,
			expectedStatus: restoreObjStatus{
				ongoing: false,
				expiry:  time.Date(2012, 12, 21, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: nil,
		},
		{
			// valid: represents an ongoing restore object request.
			restoreHdr: `ongoing-request="true"`,
			expectedStatus: restoreObjStatus{
				ongoing: true,
			},
			expectedErr: nil,
		},
		{
			// invalid; ongoing restore object request can't have expiry set on it.
			restoreHdr:     `ongoing-request="true", expiry-date="Fri, 21 Dec 2012 00:00:00 GMT"`,
			expectedStatus: restoreObjStatus{},
			expectedErr:    errRestoreHDRMalformed,
		},
		{
			// invalid; completed restore object request must have expiry set on it.
			restoreHdr:     `ongoing-request="false"`,
			expectedStatus: restoreObjStatus{},
			expectedErr:    errRestoreHDRMalformed,
		},
	}
	for i, tc := range testCases {
		actual, err := parseRestoreObjStatus(tc.restoreHdr)
		if err != tc.expectedErr {
			t.Fatalf("Test %d: got %v expected %v", i+1, err, tc.expectedErr)
		}
		if actual != tc.expectedStatus {
			t.Fatalf("Test %d: got %v expected %v", i+1, actual, tc.expectedStatus)
		}
	}
}

// TestRestoreObjStatusRoundTrip restoreObjStatus roundtrip
func TestRestoreObjStatusRoundTrip(t *testing.T) {
	testCases := []restoreObjStatus{
		ongoingRestoreObj(),
		completedRestoreObj(time.Now().UTC()),
	}
	for i, tc := range testCases {
		actual, err := parseRestoreObjStatus(tc.String())
		if err != nil {
			t.Fatalf("Test %d: parse restore object failed: %v", i+1, err)
		}
		if actual.ongoing != tc.ongoing || actual.expiry.Format(http.TimeFormat) != tc.expiry.Format(http.TimeFormat) {
			t.Fatalf("Test %d: got %v expected %v", i+1, actual, tc)
		}
	}
}

// TestRestoreObjOnDisk tests restoreObjStatus' OnDisk method
func TestRestoreObjOnDisk(t *testing.T) {
	testCases := []struct {
		restoreStatus restoreObjStatus
		ondisk        bool
	}{
		{
			// restore in progress
			restoreStatus: ongoingRestoreObj(),
			ondisk:        false,
		},
		{
			// restore completed but expired
			restoreStatus: completedRestoreObj(time.Now().Add(-time.Hour)),
			ondisk:        false,
		},
		{
			// restore completed
			restoreStatus: completedRestoreObj(time.Now().Add(time.Hour)),
			ondisk:        true,
		},
	}

	for i, tc := range testCases {
		if actual := tc.restoreStatus.OnDisk(); actual != tc.ondisk {
			t.Fatalf("Test %d: expected %v but got %v", i+1, tc.ondisk, actual)
		}
	}
}

// TestIsRestoredObjectOnDisk tests isRestoredObjectOnDisk helper function
func TestIsRestoredObjectOnDisk(t *testing.T) {
	testCases := []struct {
		meta   map[string]string
		ondisk bool
	}{
		{
			// restore in progress
			meta: map[string]string{
				xhttp.AmzRestore: ongoingRestoreObj().String(),
			},
			ondisk: false,
		},
		{
			// restore completed
			meta: map[string]string{
				xhttp.AmzRestore: completedRestoreObj(time.Now().Add(time.Hour)).String(),
			},
			ondisk: true,
		},
		{
			// restore completed but expired
			meta: map[string]string{
				xhttp.AmzRestore: completedRestoreObj(time.Now().Add(-time.Hour)).String(),
			},
			ondisk: false,
		},
	}

	for i, tc := range testCases {
		if actual := isRestoredObjectOnDisk(tc.meta); actual != tc.ondisk {
			t.Fatalf("Test %d: expected %v but got %v for %v", i+1, tc.ondisk, actual, tc.meta)
		}
	}
}

func TestObjectIsRemote(t *testing.T) {
	fi := newFileInfo("object", 8, 8)
	fi.Erasure.Index = 1
	if !fi.IsValid() {
		t.Fatalf("unable to get xl meta")
	}

	testCases := []struct {
		meta   map[string]string
		remote bool
	}{
		{
			// restore in progress
			meta: map[string]string{
				xhttp.AmzRestore: ongoingRestoreObj().String(),
			},
			remote: true,
		},
		{
			// restore completed
			meta: map[string]string{
				xhttp.AmzRestore: completedRestoreObj(time.Now().Add(time.Hour)).String(),
			},
			remote: false,
		},
		{
			// restore completed but expired
			meta: map[string]string{
				xhttp.AmzRestore: completedRestoreObj(time.Now().Add(-time.Hour)).String(),
			},
			remote: true,
		},
		{
			// restore never initiated
			meta:   map[string]string{},
			remote: true,
		},
	}
	for i, tc := range testCases {
		// Set transition status to complete
		fi.TransitionStatus = lifecycle.TransitionComplete
		fi.Metadata = tc.meta
		if got := fi.IsRemote(); got != tc.remote {
			t.Fatalf("Test %d.a: expected %v got %v", i+1, tc.remote, got)
		}
		oi := fi.ToObjectInfo("bucket", "object", false)
		if got := oi.IsRemote(); got != tc.remote {
			t.Fatalf("Test %d.b: expected %v got %v", i+1, tc.remote, got)
		}
	}
	// Reset transition status; An object that's not transitioned is not remote.
	fi.TransitionStatus = ""
	fi.Metadata = nil
	if got := fi.IsRemote(); got != false {
		t.Fatalf("Expected object not to be remote but got %v", got)
	}
}

func TestValidateTransitionTier(t *testing.T) {
	globalTierConfigMgr = NewTierConfigMgr()
	testCases := []struct {
		xml         []byte
		expectedErr error
	}{
		{
			// non-existent storage-class
			xml:         []byte(`<LifecycleConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Rule><ID>rule</ID><Prefix /><Status>Enabled</Status><Transition><Days>1</Days><StorageClass>"NONEXISTENT"</StorageClass></Transition></Rule></LifecycleConfiguration>`),
			expectedErr: errInvalidStorageClass,
		},
		{
			// no transition rule
			xml:         []byte(`<LifecycleConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Rule><ID>rule</ID><Prefix /><Status>Enabled</Status><Expiration><Days>1</Days></Expiration></Rule></LifecycleConfiguration>`),
			expectedErr: nil,
		},
	}
	for i, tc := range testCases {
		lc, err := lifecycle.ParseLifecycleConfig(bytes.NewReader(tc.xml))
		if err != nil {
			t.Fatalf("Test %d: Failed to parse lifecycle config %v", i+1, err)
		}

		err = validateTransitionTier(lc)
		if err != tc.expectedErr {
			t.Fatalf("Test %d: Expected %v but got %v", i+1, tc.expectedErr, err)
		}
	}
}
