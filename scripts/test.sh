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
cp README.md sdk/nodejs/yarn.lock sdk/nodejs/package.json  sdk/nodejs/bin/

export PATH=$PATH:$PWD/cmd/pulumi-resource-command

pushd ./examples/project
yarn install
pulumi --non-interactive destroy && pulumi --non-interactive stack rm "command-test" -y --force
pulumi --non-interactive stack init "command-test"

echo 'Running pulumi up =================================='
pulumi --non-interactive -v ${PULUMI_LOGLEVEL:-1} --logflow --logtostderr up --skip-preview
echo 'Repeat pulumi up ================================'
pulumi --non-interactive -v ${PULUMI_LOGLEVEL:-1} --logflow --logtostderr up --skip-preview
echo 'Done with pulumi up ================================'

popd
