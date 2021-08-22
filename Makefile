.PHONY: test

test:
	./scripts/test.sh

PROJECT_NAME := Pulumi Command Resource Provider

PACK             := command
PACKDIR          := sdk
PROJECT          := github.com/brandonkal/pulumi-command
NODE_MODULE_NAME := @brandonkal/pulumi-command
NUGET_PKG_NAME   := Pulumi.Command

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         ?= $(shell pulumictl get version)
PROVIDER_PATH   := provider
VERSION_PATH    := ${PROVIDER_PATH}/pkg/version.Version

SCHEMA_FILE     := provider/cmd/pulumi-resource-command/schema.json
GOPATH			:= $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

PROVIDER_BIN   := $(WORKING_DIR)/bin/
PLUGIN_ROOT     := ~/.pulumi/plugins
PLUGIN_LINK     := $(PLUGIN_ROOT)/resource-$(PACK)-v$(VERSION)

ensure::
	cd provider && go mod tidy
	cd ${PACKDIR} && go mod tidy
	cd tests && go mod tidy

gen::
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${CODEGEN} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))

provider::
	(cd provider && VERSION=${VERSION} go generate cmd/${PROVIDER}/main.go)
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

provider_debug::
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider::
	cd provider/pkg && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} ./...

dotnet_sdk:: DOTNET_VERSION := $(shell pulumictl get version --language dotnet)
dotnet_sdk:: export SDK_DIR = ${PACKDIR}
dotnet_sdk::
	rm -rf ${PACKDIR}/dotnet && mkdir -p ${PACKDIR}/dotnet && touch ${PACKDIR}/dotnet/go.mod
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${DOTNET_VERSION} dotnet $(SCHEMA_FILE) $(CURDIR)
	cd ${PACKDIR}/dotnet/&& \
		echo "${DOTNET_VERSION}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

go_sdk:: export SDK_DIR = ${PACKDIR}
go_sdk::
	rm -rf ${PACKDIR}/go/command && mkdir -p ${PACKDIR}/go/command
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} go $(SCHEMA_FILE) $(CURDIR)

nodejs_sdk:: NODE_VERSION := $(shell pulumictl get version --language javascript)
nodejs_sdk:: export SDK_DIR = ${PACKDIR}
nodejs_sdk::
	rm -rf ${PACKDIR}/nodejs && mkdir -p ${PACKDIR}/nodejs && touch ${PACKDIR}/nodejs/go.mod
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} nodejs $(SCHEMA_FILE) $(CURDIR)
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

python_sdk:: PYPI_VERSION := $(shell pulumictl get version --language python)
python_sdk:: export SDK_DIR = ${PACKDIR}
python_sdk::
	# Delete only files and folders that are generated.
	mkdir -p ${PACKDIR}/python && touch ${PACKDIR}/python/go.mod && cp README.md ${PACKDIR}/python
	rm -rf ${PACKDIR}/python/pulumi_comand/*/ ${PACKDIR}/python/pulumi_command/__init__.py
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} python $(SCHEMA_FILE) $(CURDIR)
	cp README.md ${PACKDIR}/python/
	cd ${PACKDIR}/python/ && \
		python3 setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(PYPI_VERSION)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist

.PHONY: build
build:: gen provider dotnet_sdk go_sdk nodejs_sdk python_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done


install::
	mkdir -p ${PLUGIN_LINK} && rm -rf ${PLUGIN_LINK}
	cp -r ${PROVIDER_BIN} ${PLUGIN_LINK}

GO_TEST_FAST := go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}
GO_TEST 	 := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

test_fast::
# TODO: re-enable this test once https://github.com/pulumi/pulumi/issues/4954 is fixed.
#	./${PACKDIR}/nodejs/node_modules/mocha/bin/mocha ./${PACKDIR}/nodejs/bin/tests
	cd provider/pkg && $(GO_TEST_FAST) ./...
	# cd tests/sdk/nodejs && $(GO_TEST_FAST) ./...
	cd tests/sdk/python && $(GO_TEST_FAST) ./...
	# cd tests/sdk/dotnet && $(GO_TEST_FAST) ./...
	cd tests/sdk/go && $(GO_TEST_FAST) ./...

test_all::
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...
	cd tests/sdk/python && $(GO_TEST) ./...
	cd tests/sdk/dotnet && $(GO_TEST) ./...
	cd tests/sdk/go && $(GO_TEST) ./...

install_dotnet_sdk::
	# rm -rf $(WORKING_DIR)/nuget/$(NUGET_PKG_NAME).*.nupkg
	# mkdir -p $(WORKING_DIR)/nuget
	# find . -name '*.nupkg' -print -exec cp -p {} ${WORKING_DIR}/nuget \;

install_python_sdk::
	#target intentionally blank

install_go_sdk::
	#target intentionally blank

install_nodejs_sdk::
	# -yarn unlink --cwd $(WORKING_DIR)/${PACKDIR}/nodejs/bin
	# yarn link --cwd $(WORKING_DIR)/${PACKDIR}/nodejs/bin