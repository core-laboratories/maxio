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

func TestBatchJobPrefix_UnmarshalYAML(t *testing.T) {
	type args struct {
		yamlStr string
	}
	type PrefixTemp struct {
		Prefix BatchJobPrefix `yaml:"prefix"`
	}
	tests := []struct {
		name    string
		b       PrefixTemp
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test1",
			b:    PrefixTemp{},
			args: args{
				yamlStr: `
prefix: "foo"
`,
			},
			want:    []string{"foo"},
			wantErr: false,
		},
		{
			name: "test2",
			b:    PrefixTemp{},
			args: args{
				yamlStr: `
prefix:
  - "foo"
  - "bar"
`,
			},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := yaml.Unmarshal([]byte(tt.args.yamlStr), &tt.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !slices.Equal(tt.b.Prefix.F(), tt.want) {
				t.Errorf("UnmarshalYAML() = %v, want %v", tt.b.Prefix.F(), tt.want)
			}
		})
	}
}
