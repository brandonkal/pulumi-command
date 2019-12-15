//
// conv.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package structpbconv

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
)

var json string = "{\"id\":\"id\",\"urn\":\"urn:pulumi:command-test::command-test::command:exec:command::demo\",\"olds\":{\"inputs\":{\"compare\":\"eb045d78d273107348b0300c01d29b7552d622abbc6faf81b3ec55359aa9950c\",\"create\":{\"command\":[\"ls\"]},\"diff\":{\"command\":[\"true\"]},\"update\":{\"command\":[\"bash\",\"-c\",\"echo $VAR\"],\"environment\":{\"VAR\":\"Hello Pulumi!\"}}},\"stderr\":\"\",\"stdout\":\"Pulumi.command-test.yaml\\nPulumi.yaml\\nindex.ts\\nnode_modules\\npackage.json\\ntsconfig.json\\nyarn.lock\\n\"},\"news\":{\"compare\":\"eb045d78d273107348b0300c01d29b7552d622abbc6faf81b3ec55359aa9950c\",\"create\":{\"command\":[\"ls\"]},\"diff\":{\"command\":[\"true\"]},\"update\":{\"command\":[\"bash\",\"-c\",\"echo $VAR\"],\"environment\":{\"VAR\":\"Hello Pulumi!\"}}}}"

type Input struct {
	Create cmd
	Read   cmd
	Update cmd
	Delete cmd
}

type cmd struct {
	Command     []string
	Stdin       string
	Environment map[string]string
}

func TestConvert(t *testing.T) {
	var req *pulumirpc.DiffRequest = &pulumirpc.DiffRequest{}
	err := jsonpb.UnmarshalString(json, req)
	if err != nil {
		panic("Could not unmarshal json string")
	}
	type args struct {
		// src *structpb.Struct
		dst interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Maps",
			args: args{
				// src: *req,
				dst: Input{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			news := req.GetNews()
			var dst = Input{}
			if err := Convert(news, &dst); (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Print(dst)
		})
	}
}
