# Makefile
.PHONY: build build-mac build-linux build-windows build-all clean test version

DIST := .dist
VERSION := 0.0.1
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X main.cliVersion=$(VERSION)
PROFILE ?= Release

ifeq ($(PROFILE),Release)
GO_BUILD_FLAGS := -trimpath -ldflags "-s -w $(LDFLAGS)"
else ifeq ($(PROFILE),Debug)
GO_BUILD_FLAGS := -gcflags "all=-N -l" -ldflags "$(LDFLAGS)"
else
$(error PROFILE must be Debug or Release (got $(PROFILE)))
endif

# Platform shortcuts always publish a Release binary, regardless of PROFILE.
RELEASE_BUILD_FLAGS := -trimpath -ldflags "-s -w $(LDFLAGS)"

# 打印版本信息
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"

# 构建当前平台版本
build:
	mkdir -p $(DIST)
	go build $(GO_BUILD_FLAGS) -o $(DIST)/backstage .

# 构建 macOS 版本
build-mac:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 go build $(RELEASE_BUILD_FLAGS) -o $(DIST)/backstage .

# 构建 Linux 版本
build-linux:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(DIST)/backstage-linux-amd64 .

# 构建 Windows 版本
build-windows:
	mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 go build $(RELEASE_BUILD_FLAGS) -o $(DIST)/backstage.exe .

# 构建所有平台
build-all: build-mac build-linux build-windows

# 清理构建产物
clean:
	rm -rf $(DIST)

# 运行测试
test:
	go test ./...

# 代码检查
lint:
	golangci-lint run
