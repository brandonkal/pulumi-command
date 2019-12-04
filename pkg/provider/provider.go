package provider

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/golang/glog"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	pkgerrors "github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
)

const commandType = "command:exec:command"

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

// CheckConfig validates the configuration for this resource provider.
func (p *commandProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig checks the impact a hypothetical change to this provider's configuration will have on the provider.
func (p *commandProvider) DiffConfig(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.DiffConfig(%s)", p.label(), urn)
	glog.V(9).Infof("%s executing", label)

	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.olds", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, err
	}
	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		RejectAssets: true,
	})
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "diffconfig failed because malformed resource inputs")
	}
	glog.V(9).Infof("OLDS:  %s", olds)
	glog.V(9).Infof("NEWS:  %s", news)

	return &pulumirpc.DiffResponse{
		Changes: pulumirpc.DiffResponse_DIFF_NONE,
	}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (p *commandProvider) Configure(ctx context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	return &pulumirpc.ConfigureResponse{
		AcceptSecrets: true,
	}, nil
}

func (p *commandProvider) label() string {
	return fmt.Sprintf("Provider[%s]", p.name)
}

// Invoke dynamically executes a built-in command in the provider.
func (p *commandProvider) Invoke(ctx context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	panic("Invoke not implemented")
}

// StreamInvoke dynamically executes a built-in function in the provider, which returns a stream
// of responses.
func (p *commandProvider) StreamInvoke(*pulumirpc.InvokeRequest, pulumirpc.ResourceProvider_StreamInvokeServer) error {
	panic("StreamInvoke not implemented")
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (p *commandProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Check(%s)", p.label(), urn)
	glog.V(9).Infof("%s executing", label)

	if urn.Type() != commandType {
		return nil, errors.Errorf("unknown resource type %v", urn.Type())
	}

	// map[create:{map[command:{[{echo} {hello}]}]} diff:{map[command:{[{echo} {hello}]}]} update:{map[command:{[{echo} {hello}]}]}]

	// news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
	// 	Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// Check the schema.
	// failures, err := checkProperties(news, command{})
	// failures := []*pulumirpc.CheckFailure
	// if err != nil {
	// 	return nil, err
	// }

	// We currently don't change the inputs during check.
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *commandProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	var diff []string
	return &pulumirpc.DiffResponse{
		Replaces:            diff,
		Changes:             pulumirpc.DiffResponse_DIFF_NONE,
		Stables:             diff,
		DeleteBeforeReplace: true,
	}, nil
}

func (p *commandProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Create(%s)", p.label(), urn)
	glog.V(9).Infof("%s executing", label)

	if urn.Type() != commandType {
		return nil, errors.Errorf("unknown resource type %v", urn.Type())
	}

	newResInputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.properties", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}

	glog.V(1).Infoln("newResInputs are:")
	glog.V(1).Info(newResInputs)

	// map[create:{map[command:{[{echo} {hello}]}]} diff:{map[command:{[{echo} {hello}]}]} update:{map[command:{[{echo} {hello}]}]}]
	glog.V(1).Info(newResInputs["create"])
	cr := newResInputs["create"]
	// createArgs := cr["command"]
	theMap := newResInputs.Copy()

	// Execute the Command
	cmd := exec.CommandContext(ctx, "echo", "hello world")
	cmd.Run()
	cmd.Env = []string{"SOMETHING=true"}

	// var f function
	// if err := decodeProperties(newResInputs, &f); err != nil {
	// 	return nil, err
	// }

	// clientFunc := &client.Function{
	// 	Service:      f.Service,
	// 	Network:      f.Network,
	// 	Image:        f.Image,
	// 	EnvProcess:   f.EnvProcess,
	// 	EnvVars:      f.EnvVars,
	// 	Labels:       f.Labels,
	// 	Annotations:  f.Annotations,
	// 	Secrets:      f.Secrets,
	// 	RegistryAuth: f.RegistryAuth,
	// }

	// if err := p.client.CreateFunction(p.canceler.context, clientFunc); err != nil {
	// 	return nil, err
	// }

	return &pulumirpc.CreateResponse{
		Id: "id", Properties: req.GetProperties(),
	}, nil
}

// Read the current live state associated with a resource.  Enough state must be include in the inputs to uniquely
// identify the resource; this is typically just the resource ID, but may also include some properties.
func (p *commandProvider) Read(context.Context, *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	return &pulumirpc.ReadResponse{}, errors.New("Read not implemented")
}

func (p *commandProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	return &pulumirpc.UpdateResponse{Properties: req.GetNews()}, nil
}

// Delete tears down an existing resource with the given ID.
// If it fails, the resource is assumed to still exist.
func (p *commandProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
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
