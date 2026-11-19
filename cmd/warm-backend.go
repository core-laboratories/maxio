// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minio/madmin-go/v3"
	xhttp "github.com/core-laboratories/maxio/internal/http"
)

// WarmBackendGetOpts is used to express byte ranges within an object. The zero
// value represents the entire byte range of an object.
type WarmBackendGetOpts struct {
	startOffset int64
	length      int64
}

// WarmBackend provides interface to be implemented by remote tier backends
type WarmBackend interface {
	Put(ctx context.Context, object string, r io.Reader, length int64) (remoteVersionID, error)
	PutWithMeta(ctx context.Context, object string, r io.Reader, length int64, meta map[string]string) (remoteVersionID, error)
	Get(ctx context.Context, object string, rv remoteVersionID, opts WarmBackendGetOpts) (io.ReadCloser, error)
	Remove(ctx context.Context, object string, rv remoteVersionID) error
	InUse(ctx context.Context) (bool, error)
}

const probeObject = "probeobject"

// checkWarmBackend checks if tier config credentials have sufficient privileges
// to perform all operations defined in the WarmBackend interface.
func checkWarmBackend(ctx context.Context, w WarmBackend) error {
	remoteVersionID, err := w.Put(ctx, probeObject, strings.NewReader("MinIO"), 5)
	if err != nil {
		if _, ok := err.(BackendDown); ok {
			return err
		}
		return tierPermErr{
			Op:  tierPut,
			Err: err,
		}
	}

	r, err := w.Get(ctx, probeObject, "", WarmBackendGetOpts{})
	xhttp.DrainBody(r)
	if err != nil {
		if _, ok := err.(BackendDown); ok {
			return err
		}
		switch {
		case isErrBucketNotFound(err):
			return errTierBucketNotFound
		case isErrSignatureDoesNotMatch(err):
			return errTierInvalidCredentials
		default:
			return tierPermErr{
				Op:  tierGet,
				Err: err,
			}
		}
	}
	if err = w.Remove(ctx, probeObject, remoteVersionID); err != nil {
		if _, ok := err.(BackendDown); ok {
			return err
		}
		return tierPermErr{
			Op:  tierDelete,
			Err: err,
		}
	}
	return err
}

type tierOp uint8

const (
	_ tierOp = iota
	tierGet
	tierPut
	tierDelete
)

func (op tierOp) String() string {
	switch op {
	case tierGet:
		return "GET"
	case tierPut:
		return "PUT"
	case tierDelete:
		return "DELETE"
	}
	return "UNKNOWN"
}

type tierPermErr struct {
	Op  tierOp
	Err error
}

func (te tierPermErr) Error() string {
	return fmt.Sprintf("failed to perform %s: %v", te.Op, te.Err)
}

func errIsTierPermError(err error) bool {
	var tpErr tierPermErr
	return errors.As(err, &tpErr)
}

// remoteVersionID represents the version id of an object in the remote tier.
// Its usage is remote tier cloud implementation specific.
type remoteVersionID string

// newWarmBackend instantiates the tier type specific WarmBackend, runs
// checkWarmBackend on it.
func newWarmBackend(ctx context.Context, tier madmin.TierConfig, probe bool) (d WarmBackend, err error) {
	switch tier.Type {
	case madmin.S3:
		d, err = newWarmBackendS3(*tier.S3, tier.Name)
	case madmin.Azure:
		d, err = newWarmBackendAzure(*tier.Azure, tier.Name)
	case madmin.GCS:
		d, err = newWarmBackendGCS(*tier.GCS, tier.Name)
	case madmin.MinIO:
		d, err = newWarmBackendMinIO(*tier.MinIO, tier.Name)
	default:
		return nil, errTierTypeUnsupported
	}
	if err != nil {
		tierLogIf(ctx, err)
		return nil, errTierInvalidConfig
	}

	if probe {
		if err = checkWarmBackend(ctx, d); err != nil {
			return nil, err
		}
	}

	return d, nil
}
