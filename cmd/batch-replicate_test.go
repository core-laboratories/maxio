// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"slices"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseBatchJobReplicate(t *testing.T) {
	replicateYaml := `
replicate:
  apiVersion: v1
  # source of the objects to be replicated
  source:
    type: minio # valid values are "s3" or "minio"
    bucket: mytest
    prefix: object-prefix1 # 'PREFIX' is optional
    # If your source is the 'local' alias specified to 'mc batch start', then the 'endpoint' and 'credentials' fields are optional and can be omitted
    # Either the 'source' or 'remote' *must* be the "local" deployment
#    endpoint: "http://127.0.0.1:9000"
#    # path: "on|off|auto" # "on" enables path-style bucket lookup. "off" enables virtual host (DNS)-style bucket lookup. Defaults to "auto"
#    credentials:
#      accessKey: minioadmin # Required
#      secretKey: minioadmin # Required
#    # sessionToken: SESSION-TOKEN # Optional only available when rotating credentials are used
    snowball: # automatically activated if the source is local
      disable: true # optionally turn-off snowball archive transfer
#      batch: 100 # upto this many objects per archive
#      inmemory: true # indicates if the archive must be staged locally or in-memory
#      compress: false # S2/Snappy compressed archive
#      smallerThan: 5MiB # create archive for all objects smaller than 5MiB
#      skipErrs: false # skips any source side read() errors

  # target where the objects must be replicated
  target:
    type: minio # valid values are "s3" or "minio"
    bucket: mytest
    prefix: stage # 'PREFIX' is optional
    # If your source is the 'local' alias specified to 'mc batch start', then the 'endpoint' and 'credentials' fields are optional and can be omitted

    # Either the 'source' or 'remote' *must* be the "local" deployment
    endpoint: "http://127.0.0.1:9001"
    # path: "on|off|auto" # "on" enables path-style bucket lookup. "off" enables virtual host (DNS)-style bucket lookup. Defaults to "auto"
    credentials:
      accessKey: minioadmin
      secretKey: minioadmin
    # sessionToken: SESSION-TOKEN # Optional only available when rotating credentials are used

  # NOTE: All flags are optional
  # - filtering criteria only applies for all source objects match the criteria
  # - configurable notification endpoints
  # - configurable retries for the job (each retry skips successfully previously replaced objects)
  flags:
    filter:
      newerThan: "7d10h31s" # match objects newer than this value (e.g. 7d10h31s)
      olderThan: "7d" # match objects older than this value (e.g. 7d10h31s)
#      createdAfter: "date" # match objects created after "date"
#      createdBefore: "date" # match objects created before "date"

      ## NOTE: tags are not supported when "source" is remote.
      tags:
         - key: "name"
           value: "pick*" # match objects with tag 'name', with all values starting with 'pick'

      metadata:
         - key: "content-type"
           value: "image/*" # match objects with 'content-type', with all values starting with 'image/'

#    notify:
#      endpoint: "https://notify.endpoint" # notification endpoint to receive job status events
#      token: "Bearer xxxxx" # optional authentication token for the notification endpoint
#
#    retry:
#      attempts: 10 # number of retries for the job before giving up
#      delay: "500ms" # least amount of delay between each retry

`
	var job BatchJobRequest
	err := yaml.Unmarshal([]byte(replicateYaml), &job)
	if err != nil {
		t.Fatal("Failed to parse batch-job-replicate yaml", err)
	}
	if !slices.Equal(job.Replicate.Source.Prefix.F(), []string{"object-prefix1"}) {
		t.Fatal("Failed to parse batch-job-replicate yaml", err)
	}
	multiPrefixReplicateYaml := `
replicate:
  apiVersion: v1
  # source of the objects to be replicated
  source:
    type: minio # valid values are "s3" or "minio"
    bucket: mytest
    prefix: # 'PREFIX' is optional
      - object-prefix1 
      - object-prefix2
    # If your source is the 'local' alias specified to 'mc batch start', then the 'endpoint' and 'credentials' fields are optional and can be omitted
    # Either the 'source' or 'remote' *must* be the "local" deployment
#    endpoint: "http://127.0.0.1:9000"
#    # path: "on|off|auto" # "on" enables path-style bucket lookup. "off" enables virtual host (DNS)-style bucket lookup. Defaults to "auto"
#    credentials:
#      accessKey: minioadmin # Required
#      secretKey: minioadmin # Required
#    # sessionToken: SESSION-TOKEN # Optional only available when rotating credentials are used
    snowball: # automatically activated if the source is local
      disable: true # optionally turn-off snowball archive transfer
#      batch: 100 # upto this many objects per archive
#      inmemory: true # indicates if the archive must be staged locally or in-memory
#      compress: false # S2/Snappy compressed archive
#      smallerThan: 5MiB # create archive for all objects smaller than 5MiB
#      skipErrs: false # skips any source side read() errors

  # target where the objects must be replicated
  target:
    type: minio # valid values are "s3" or "minio"
    bucket: mytest
    prefix: stage # 'PREFIX' is optional
    # If your source is the 'local' alias specified to 'mc batch start', then the 'endpoint' and 'credentials' fields are optional and can be omitted

    # Either the 'source' or 'remote' *must* be the "local" deployment
    endpoint: "http://127.0.0.1:9001"
    # path: "on|off|auto" # "on" enables path-style bucket lookup. "off" enables virtual host (DNS)-style bucket lookup. Defaults to "auto"
    credentials:
      accessKey: minioadmin
      secretKey: minioadmin
    # sessionToken: SESSION-TOKEN # Optional only available when rotating credentials are used

  # NOTE: All flags are optional
  # - filtering criteria only applies for all source objects match the criteria
  # - configurable notification endpoints
  # - configurable retries for the job (each retry skips successfully previously replaced objects)
  flags:
    filter:
      newerThan: "7d10h31s" # match objects newer than this value (e.g. 7d10h31s)
      olderThan: "7d" # match objects older than this value (e.g. 7d10h31s)
#      createdAfter: "date" # match objects created after "date"
#      createdBefore: "date" # match objects created before "date"

      ## NOTE: tags are not supported when "source" is remote.
      tags:
         - key: "name"
           value: "pick*" # match objects with tag 'name', with all values starting with 'pick'

      metadata:
         - key: "content-type"
           value: "image/*" # match objects with 'content-type', with all values starting with 'image/'

#    notify:
#      endpoint: "https://notify.endpoint" # notification endpoint to receive job status events
#      token: "Bearer xxxxx" # optional authentication token for the notification endpoint
#
#    retry:
#      attempts: 10 # number of retries for the job before giving up
#      delay: "500ms" # least amount of delay between each retry

`
	var multiPrefixJob BatchJobRequest
	err = yaml.Unmarshal([]byte(multiPrefixReplicateYaml), &multiPrefixJob)
	if err != nil {
		t.Fatal("Failed to parse batch-job-replicate yaml", err)
	}
	if !slices.Equal(multiPrefixJob.Replicate.Source.Prefix.F(), []string{"object-prefix1", "object-prefix2"}) {
		t.Fatal("Failed to parse batch-job-replicate yaml")
	}
}
