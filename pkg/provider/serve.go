package provider

import (
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	lumirpc "github.com/pulumi/pulumi/sdk/proto/go"
)

// Serve launches the gRPC server for the Pulumi Kubernetes resource provider.
func Serve(providerName, version string) {
	// Start gRPC service.
	err := provider.Main(
		providerName, func(host *provider.HostClient) (lumirpc.ResourceProviderServer, error) {
			return makeCommandProvider(providerName, version)
		})

	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}
