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
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// Serve launches the gRPC server for the Pulumi Kubernetes resource provider.
func Serve(providerName, version string) {
	// Start gRPC service.
	err := provider.Main(
		providerName, func(host *provider.HostClient) (pulumirpc.ResourceProviderServer, error) {
			return makeCommandProvider(providerName, version)
		})

	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}
