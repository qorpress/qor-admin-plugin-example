.SILENT :

export GO111MODULE=on

# Base package
BASE_PACKAGE=github.com/qorpress

# App name
APPNAME=qor-admin-plugin-example

# Go configuration
GOOS?=$(shell go env GOHOSTOS)
GOARCH?=$(shell go env GOHOSTARCH)

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")
LDLAGS_LAUNCHER:=$(if $(is_windows),-ldflags "-H=windowsgui",)

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Plugin name
PLUGIN?=oniontree

# Plugin filename
PLUGIN_SO=$(APPNAME)-$(PLUGIN).so

# Extract version infos
#PKG_VERSION:=github.com/qorpress/$(APPNAME)/v1/pkg/version
VERSION:=`git describe --tags --always`
GIT_COMMIT:=`git rev-list -1 HEAD --abbrev-commit`
#BUILT:=`date`
#define LDFLAGS
#-X '$(PKG_VERSION).Version=$(VERSION)' \
#-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
#-X '$(PKG_VERSION).Built=$(BUILT)'
#endef

dep:
	go mod vendor
.PHONY: dep

run:
	go run main.go
.PHONY: run

run_bindatafs:
	go run config/compile/compile.go
	go run -tags bindatafs main.go
.PHONY: run_bindatafs

build: dep
	go build main.go
.PHONY: build

build_bindatafs: dep
	go run config/compile/compile.go
	go build -tags bindatafs main.go
.PHONY: build_bindatafs

## Bulid plugin (defined by PLUGIN variable)
plugin:
	-mkdir -p release
	echo ">>> Building: $(PLUGIN_SO) $(VERSION) for $(GOOS)-$(GOARCH) ..."
	cd plugins/$(PLUGIN) && GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildmode=plugin -o ../../plugins/$(PLUGIN_SO)
.PHONY: plugin

## Build all plugins
plugins:
	GOARCH=amd64 PLUGIN=bubble_sort make plugin
	GOARCH=amd64 PLUGIN=quick_sort make plugin
	GOARCH=amd64 PLUGIN=oniontree make plugin
.PHONY: plugins   
