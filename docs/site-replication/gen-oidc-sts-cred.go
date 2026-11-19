//go:build ignore
// +build ignore

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package main

// This programs mocks user interaction against Dex IDP and generates STS
// credentials. It is for MinIO testing purposes only.
//
// Run like:
//
// $ MINIO_ENDPOINT=http://localhost:9000 go run gen-oidc-sts-cred.go

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	cr "github.com/minio/minio-go/v7/pkg/credentials"
	cmd "github.com/core-laboratories/maxio/cmd"
)

func main() {
	ctx := context.Background()

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		log.Fatalf("Please specify a MinIO server endpoint environment variable like:\n\n\texport MINIO_ENDPOINT=http://localhost:9000")
	}

	appParams := cmd.OpenIDClientAppParams{
		ClientID:     "minio-client-app",
		ClientSecret: "minio-client-app-secret",
		ProviderURL:  "http://127.0.0.1:5556/dex",
		RedirectURL:  "http://127.0.0.1:10000/oauth_callback",
	}

	oidcToken, err := cmd.MockOpenIDTestUserInteraction(ctx, appParams, "dillon@example.io", "dillon")
	if err != nil {
		log.Fatalf("Failed to generate OIDC token: %v", err)
	}

	roleARN := os.Getenv("ROLE_ARN")
	webID := cr.STSWebIdentity{
		Client:      &http.Client{},
		STSEndpoint: endpoint,
		GetWebIDTokenExpiry: func() (*cr.WebIdentityToken, error) {
			return &cr.WebIdentityToken{
				Token: oidcToken,
			}, nil
		},
		RoleARN: roleARN,
	}

	value, err := webID.Retrieve()
	if err != nil {
		log.Fatalf("Expected to generate credentials: %v", err)
	}

	// Print credentials separated by colons:
	fmt.Printf("%s:%s:%s\n", value.AccessKeyID, value.SecretAccessKey, value.SessionToken)
}
