MAKEFLAGS += --jobs all
GO := GOEXPERIMENT=jsonv2,greenteagc CGO_ENABLED=0 go

VERSION := $(shell git describe --tags --always || echo 'dev')
DATE := $(shell date +%Y-%m-%d)
BINARY_NAME := apihub-mcp



test:
	GOEXPERIMENT=jsonv2 CGO_ENABLED=1 go vet ./...
	GOEXPERIMENT=jsonv2 CGO_ENABLED=1 go test ./...

# 构建当前平台的二进制文件
build:
	$(GO) build -trimpath -ldflags "-X main.Version=$(VERSION) -X main.Date=$(DATE)" -o bin/$(BINARY_NAME) ./cmd/apihub-mcp


# 交叉编译所有平台
build-all:
	mkdir -p bin
	@$(foreach platform,$(PLATFORMS), \
		$(call build_platform,$(platform)))

# 目标平台配置
PLATFORMS ?= \
	darwin/amd64 \
	darwin/arm64 \
	linux/386 \
	linux/amd64 \
	linux/arm \
	linux/arm64 \
	windows/386 \
	windows/amd64 \
	windows/arm64

# 构建单个平台的内部函数
define build_platform
	$(eval OS := $(firstword $(subst /, ,$(1))))
	$(eval ARCH := $(lastword $(subst /, ,$(1))))
	$(eval OUTPUT := bin/$(BINARY_NAME)_$(OS)_$(ARCH))
	$(eval EXT := $(if $(filter windows,$(OS)),.exe,))
	@echo "Building for $(OS)/$(ARCH)..."
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -trimpath -ldflags "-X main.Version=$(VERSION) -X main.Date=$(DATE)" -o $(OUTPUT)$(EXT) ./cmd/apihub-mcp
endef

# 清理构建产物
clean:
	rm -rf bin/
