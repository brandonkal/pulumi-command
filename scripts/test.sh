#!/bin/bash
set -o -v -e nounset -o errexit -o pipefail

export CI=1
export PULUMI_CONFIG_PASSPHRASE=password

echo Running Integration Tests...

pushd cmd/pulumi-resource-command
go build
popd

pushd ./sdk/nodejs
yarn install
yarn build
yarn link
popd

export PATH=$PATH:$PWD/cmd/pulumi-resource-command
echo $PATH

pushd ./examples/project
yarn install
pulumi --non-interactive destroy && pulumi --non-interactive stack rm "command-test" -y --force
pulumi --non-interactive stack init "command-test"

echo $PATH
echo 'Running pulumi up =================================='
pulumi --non-interactive up --skip-preview -v 9 --logflow
echo 'Done with pulumi up ================================'

pulumi --non-interactive destroy
pulumi --non-interactive stack rm "command-test" -y
popd
