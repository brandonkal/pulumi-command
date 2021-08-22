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
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strings"

	"github.com/brandonkal/pulumi-command/provider/pkg/structpbconv"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type cmd struct {
	Command     []string          `pulumi:"command"`
	Stdin       string            `pulumi:"stdin,optional"`
	Environment map[string]string `pulumi:"environment,optional"`
}

const (
	commandType               = "command:v1:Command"
	backwardCompatCommandType = "command:v1:exec"
)

type cancellationContext struct {
	context context.Context
	cancel  context.CancelFunc
}

func makeCancellationContext() *cancellationContext {
	var ctx, cancel = context.WithCancel(context.Background())
	return &cancellationContext{
		context: ctx,
		cancel:  cancel,
	}
}

type commandProvider struct {
	canceler *cancellationContext
	name     string
	version  string
}

func makeCommandProvider(name, version string) (pulumirpc.ResourceProviderServer, error) {
	return &commandProvider{
		canceler: makeCancellationContext(),
		name:     name,
		version:  version,
	}, nil
}

func (p *commandProvider) prepare(req hasUrn, op string, props *structpb.Struct, path string) (resource.PropertyMap, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.%s(%s)", p.label(), op, urn)
	logging.V(9).Infof("%s executing", label)
	if urn.Type() != commandType && urn.Type() != backwardCompatCommandType {
		return nil, errors.Errorf("unknown resource type %v", urn.Type())
	}
	input, err := plugin.UnmarshalProperties(props, plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.%s", label, path), KeepUnknowns: true, SkipNulls: true,
	})
	return input, err
}

// execCommand runs the specified command and returns a proto structure containing stderr and stdout
// if exitCode is zero and an error is returned, it is an internal error
func (p *commandProvider) execCommand(ctx context.Context, req hasUrn, op string, props *structpb.Struct, path string) (out *structpb.Struct, err error, code int) {
	code = 0
	input, err := p.prepare(req, op, props, path)
	if err != nil {
		return nil, err, code
	}

	if path == "olds" {
		if inputs, ok := input["inputs"]; ok {
			input = inputs.ObjectValue()
		}
	}

	what := input[(resource.PropertyKey)(op)]
	var this cmd
	if what.V == nil {
		return nil, errors.Errorf("%s command unspecified", op), code
	}
	err = decodeProperty("", what, reflect.ValueOf(&this))
	if err != nil {
		return nil, err, code
	}

	envs := this.Environment
	var environment = []string{}
	for k, v := range envs {
		environment = append(environment, fmt.Sprintf("%s=%s", k, v))
	}

	args := this.Command

	// Prepare the Command
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	if len(environment) > 0 {
		cmd.Env = environment
	}
	if len(this.Stdin) > 0 {
		r := strings.NewReader(this.Stdin)
		cmd.Stdin = r
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			code = exitError.ExitCode()
			err = errors.Wrap(err, stderr.String())
			logging.V(1).Infof("Command exit with code: %v", code)
		} else {
			return nil, err, code
		}
	}

	m := make(map[string]*structpb.Value)
	m["stdout"] = &structpb.Value{
		Kind: &structpb.Value_StringValue{StringValue: stdout.String()},
	}
	m["stderr"] = &structpb.Value{
		Kind: &structpb.Value_StringValue{StringValue: stderr.String()},
	}
	out = &structpb.Struct{
		Fields: m,
	}

	return out, err, code
}

// Call dynamically executes a method in the provider associated with a component resource.
func (k *commandProvider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Call is not yet implemented")
}

// Construct creates a new component resource.
func (k *commandProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

// CheckConfig validates the configuration for this resource provider.
func (p *commandProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig checks the impact a hypothetical change to this provider's configuration will have on the provider.
func (p *commandProvider) DiffConfig(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return &pulumirpc.DiffResponse{
		Changes: pulumirpc.DiffResponse_DIFF_NONE,
	}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (p *commandProvider) Configure(ctx context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	return &pulumirpc.ConfigureResponse{
		AcceptSecrets: false,
	}, nil
}

func (p *commandProvider) label() string {
	return fmt.Sprintf("Provider[%s]", p.name)
}

// Invoke dynamically executes a built-in command in the provider.
func (p *commandProvider) Invoke(ctx context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Invoke is not yet implemented")
}

// StreamInvoke dynamically executes a built-in function in the provider, which returns a stream
// of responses.
func (p *commandProvider) StreamInvoke(*pulumirpc.InvokeRequest, pulumirpc.ResourceProvider_StreamInvokeServer) error {
	return status.Error(codes.Unimplemented, "StreamInvoke is not yet implemented")
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (p *commandProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	// We currently don't change the inputs during check.
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

type Input struct {
	Compare string
	Create  cmd
	Read    cmd
	Update  cmd
	Delete  cmd
}

func isEmpty(item Input) bool {
	if item.Compare == "" && len(item.Create.Command) == 0 && item.Read.Command == nil && item.Update.Command == nil && item.Delete.Command == nil {
		return true
	}
	return false
}

type OldDiff struct {
	Inputs Input
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
// It first checks to see if inputs have changed. If they have not, it executes the diff command.
func (p *commandProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Diff(%s)", p.label(), urn)
	logging.V(9).Infof("%s executing", label)

	var oldDiff = OldDiff{}
	olds := req.GetOlds()
	news := req.GetNews()
	err := structpbconv.Convert(olds, &oldDiff)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert input")
	}
	var newInput = Input{}
	err = structpbconv.Convert(news, &newInput)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert input")
	}
	logging.V(9).Info("===OLD-DIFF===")
	logging.V(9).Info(oldDiff)
	logging.V(9).Info("===newInput===")
	logging.V(9).Info(newInput)

	var needsUpdate = false
	wasEmpty := isEmpty(oldDiff.Inputs)
	if !wasEmpty {
		depChanged := oldDiff.Inputs.Compare != newInput.Compare
		updateCmdChanged := !reflect.DeepEqual(oldDiff.Inputs.Update, newInput.Update)
		logging.V(1).Infof("Diff check: depChanged: %v. updateCmdChanged: %v", depChanged, updateCmdChanged)
		needsUpdate = depChanged || updateCmdChanged
	} else {
		logging.V(1).Info("oldDiff empty")
	}
	if !needsUpdate {
		_, err, code := p.execCommand(ctx, req, "diff", news, "news")
		// If the user doesn't provide a diff command, we never run update
		if err != nil && err.Error() != "diff command unspecified" && code == 0 {
			return nil, err
		}
		unspecified := (err != nil && err.Error() == "diff command unspecified")
		if code == 0 && !unspecified {
			// code is zero if no command specified or process returned success (update)
			needsUpdate = true
			logging.V(1).Infof("Diff check update required: return code: %v. unspecified? %v", code, unspecified)
		}
	}
	diff := pulumirpc.DiffResponse_DIFF_SOME
	if !needsUpdate {
		diff = pulumirpc.DiffResponse_DIFF_NONE
	}
	logging.V(1).Infof("Diff check needs update: %v", needsUpdate)

	return &pulumirpc.DiffResponse{
		Replaces:            []string{},
		Changes:             diff,
		Stables:             []string{},
		DeleteBeforeReplace: true,
	}, nil
}

func (p *commandProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	out, err, _ := p.execCommand(ctx, req, "create", req.GetProperties(), "properties")
	if err != nil {
		return nil, err
	}
	// Save inputs to state
	out.Fields["inputs"] = &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: req.Properties}}

	return &pulumirpc.CreateResponse{
		Id: "id", Properties: out,
	}, nil
}

type hasUrn interface {
	GetUrn() string
}

// Read the current live state associated with a resource.  Enough state must be include in the inputs to uniquely
// identify the resource; this is typically just the resource ID, but may also include some properties.
func (p *commandProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	out, err, _ := p.execCommand(ctx, req, "read", req.GetInputs(), "olds")
	if err != nil && err.Error() != "read command unspecified" {
		return nil, err
	}
	properties := req.GetProperties()
	if err.Error() != "read command unspecified" {
		properties = out
	}
	return &pulumirpc.ReadResponse{Id: req.GetId(), Properties: properties}, nil
}

func (p *commandProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	news := req.GetNews()
	out, err, _ := p.execCommand(ctx, req, "update", news, "properties")
	if err != nil {
		return nil, err
	}
	// Save inputs to state
	out.Fields["inputs"] = &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: news}}
	return &pulumirpc.UpdateResponse{Properties: out}, nil
}

// Delete tears down an existing resource with the given ID.
// If it fails, the resource is assumed to still exist.
func (p *commandProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	_, err, _ := p.execCommand(ctx, req, "delete", req.GetProperties(), "olds")
	if err != nil && err.Error() != "delete command unspecified" {
		return nil, err
	}

	if err != nil && err.Error() == "delete command unspecified" {
		logging.V(9).Infof("Skipping deleting resource: %v", err)
	}

	return &pbempty.Empty{}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
func (p *commandProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	p.canceler.cancel()
	return &pbempty.Empty{}, nil
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (p *commandProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: p.version,
	}, nil
}

// GetSchema returns the JSON-serialized schema for the provider.
func (k *commandProvider) GetSchema(ctx context.Context, req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	return &pulumirpc.GetSchemaResponse{}, nil
}
