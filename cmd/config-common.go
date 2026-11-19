// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/core-laboratories/maxio/internal/hash"
)

var errConfigNotFound = errors.New("config file not found")

func readConfigWithMetadata(ctx context.Context, store objectIO, configFile string, opts ObjectOptions) ([]byte, ObjectInfo, error) {
	r, err := store.GetObjectNInfo(ctx, minioMetaBucket, configFile, nil, http.Header{}, opts)
	if err != nil {
		if isErrObjectNotFound(err) {
			return nil, ObjectInfo{}, errConfigNotFound
		}

		return nil, ObjectInfo{}, err
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, ObjectInfo{}, err
	}
	if len(buf) == 0 {
		return nil, ObjectInfo{}, errConfigNotFound
	}
	return buf, r.ObjInfo, nil
}

func readConfig(ctx context.Context, store objectIO, configFile string) ([]byte, error) {
	buf, _, err := readConfigWithMetadata(ctx, store, configFile, ObjectOptions{})
	return buf, err
}

type objectDeleter interface {
	DeleteObject(ctx context.Context, bucket, object string, opts ObjectOptions) (ObjectInfo, error)
}

func deleteConfig(ctx context.Context, objAPI objectDeleter, configFile string) error {
	_, err := objAPI.DeleteObject(ctx, minioMetaBucket, configFile, ObjectOptions{
		DeletePrefix:       true,
		DeletePrefixObject: true, // use prefix delete on exact object (this is an optimization to avoid fan-out calls)
	})
	if err != nil && isErrObjectNotFound(err) {
		return errConfigNotFound
	}
	return err
}

func saveConfigWithOpts(ctx context.Context, store objectIO, configFile string, data []byte, opts ObjectOptions) error {
	hashReader, err := hash.NewReader(ctx, bytes.NewReader(data), int64(len(data)), "", getSHA256Hash(data), int64(len(data)))
	if err != nil {
		return err
	}

	_, err = store.PutObject(ctx, minioMetaBucket, configFile, NewPutObjReader(hashReader), opts)
	return err
}

func saveConfig(ctx context.Context, store objectIO, configFile string, data []byte) error {
	return saveConfigWithOpts(ctx, store, configFile, data, ObjectOptions{MaxParity: true})
}

func checkConfig(ctx context.Context, objAPI ObjectLayer, configFile string) error {
	if _, err := objAPI.GetObjectInfo(ctx, minioMetaBucket, configFile, ObjectOptions{}); err != nil {
		// Treat object not found as config not found.
		if isErrObjectNotFound(err) {
			return errConfigNotFound
		}

		return err
	}
	return nil
}
