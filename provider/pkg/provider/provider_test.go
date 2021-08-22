// Copyright 2019, Brandon Kalinowski.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var json string = "{\"id\":\"id\",\"urn\":\"urn:pulumi:command-test::command-test::command:v1:exec::demo\",\"olds\":{\"inputs\":{\"compare\":\"eb045d78d273107348b0300c01d29b7552d622abbc6faf81b3ec55359aa9950c\",\"create\":{\"command\":[\"ls\"]},\"diff\":{\"command\":[\"bash\",\"-c\",\"exit 1\"]},\"update\":{\"command\":[\"bash\",\"-c\",\"echo $VAR\"],\"environment\":{\"VAR\":\"Hello Pulumi!\"}}},\"stderr\":\"\",\"stdout\":\"Pulumi.command-test.yaml\\nPulumi.yaml\\nindex.ts\\nnode_modules\\npackage.json\\ntsconfig.json\\nyarn.lock\\n\"},\"news\":{\"compare\":\"eb045d78d273107348b0300c01d29b7552d622abbc6faf81b3ec55359aa9950c\",\"create\":{\"command\":[\"ls\"]},\"diff\":{\"command\":[\"bash\",\"-c\",\"exit 1\"]},\"update\":{\"command\":[\"bash\",\"-c\",\"echo $VAR\"],\"environment\":{\"VAR\":\"Hello Pulumi!\"}}}}"

func Test_commandProvider_Diff(t *testing.T) {
	var req *pulumirpc.DiffRequest = &pulumirpc.DiffRequest{}
	err := jsonpb.UnmarshalString(json, req)
	if err != nil {
		panic("Could not unmarshal json string")
	}
	type fields struct {
		canceler *cancellationContext
		name     string
		version  string
	}
	type args struct {
		ctx context.Context
		req *pulumirpc.DiffRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pulumirpc.DiffResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test One",
			fields: fields{
				name:     "command",
				version:  "dev",
				canceler: makeCancellationContext(),
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want: &pulumirpc.DiffResponse{
				Changes:             pulumirpc.DiffResponse_DIFF_NONE,
				DeleteBeforeReplace: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &commandProvider{
				canceler: tt.fields.canceler,
				name:     tt.fields.name,
				version:  tt.fields.version,
			}
			got, err := p.Diff(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandProvider.Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Changes, tt.want.Changes) {
				t.Errorf("commandProvider.Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
