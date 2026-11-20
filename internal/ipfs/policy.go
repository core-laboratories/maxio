// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ipfs

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Access string

const (
	Private Access = "private"
	Public  Access = "public"
)

// Policy is a request. It is never an authorization grant.
type Policy struct {
	Pin       *bool
	Access    *Access
	Republish *bool
	CoreIDs   []string
}

func ParseMetadata(metadata map[string]string) (Policy, error) {
	fields := make(map[string]string)
	for key, value := range metadata {
		k := strings.ToLower(key)
		if strings.HasPrefix(k, "x-amz-meta-") {
			k = strings.TrimPrefix(k, "x-amz-meta-")
		}
		if strings.HasPrefix(k, "maxio-") {
			fields[k] = value
		}
	}
	return parsePolicy(fields, true)
}

func ParseTags(encoded string) (Policy, error) {
	values, err := url.ParseQuery(encoded)
	if err != nil {
		return Policy{}, err
	}
	fields := make(map[string]string)
	for key, vals := range values {
		if strings.HasPrefix(key, "maxio-") {
			if len(vals) != 1 {
				return Policy{}, fmt.Errorf("duplicate IPFS tag %q", key)
			}
			fields[key] = vals[0]
		}
	}
	return parsePolicy(fields, false)
}

func parsePolicy(fields map[string]string, metadata bool) (Policy, error) {
	var p Policy
	for key, value := range fields {
		switch key {
		case "maxio-pin":
			b, err := parseBool(value)
			if err != nil {
				return Policy{}, fmt.Errorf("%s: %w", key, err)
			}
			p.Pin = &b
		case "maxio-access":
			a := Access(value)
			if a != Public && a != Private {
				return Policy{}, fmt.Errorf("%s: invalid access", key)
			}
			p.Access = &a
		case "maxio-ipns-republish":
			if metadata {
				return Policy{}, errors.New("IPNS republishing is controlled by S3 tags")
			}
			b, err := parseBool(value)
			if err != nil {
				return Policy{}, fmt.Errorf("%s: %w", key, err)
			}
			if b {
				return Policy{}, errors.New("IPNS republishing is enabled at startup; only false is accepted on a record")
			}
			p.Republish = &b
		case "maxio-core-ids":
			if !metadata {
				return Policy{}, errors.New("CORE IDs are not supported in object tags")
			}
			if len(value) > 2048 {
				return Policy{}, errors.New("CORE ID list is too large")
			}
			for _, id := range strings.Split(value, ",") {
				id = strings.TrimSpace(id)
				if id == "" || len(id) > 128 {
					return Policy{}, errors.New("invalid CORE ID")
				}
				p.CoreIDs = append(p.CoreIDs, id)
			}
		default:
			return Policy{}, fmt.Errorf("unsupported IPFS control %q", key)
		}
	}
	return p, nil
}

func parseBool(v string) (bool, error) {
	switch v {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("expected true or false")
	}
}

func EffectivePin(serverDefault bool, policies ...Policy) bool {
	pin := serverDefault
	for _, policy := range policies {
		if policy.Pin != nil {
			pin = *policy.Pin
		}
	}
	return pin
}
