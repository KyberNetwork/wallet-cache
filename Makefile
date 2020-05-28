# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: cache android ios cache-cross swarm evm all test clean
.PHONY: cache-linux cache-linux-386 cache-linux-amd64 cache-linux-mips64 cache-linux-mips64le
.PHONY: cache-linux-arm cache-linux-arm-5 cache-linux-arm-6 cache-linux-arm-7 cache-linux-arm64
.PHONY: cache-darwin cache-darwin-386 cache-darwin-amd64
.PHONY: cache-windows cache-windows-386 cache-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

cache:
	build/env.sh go run build/ci.go install ./cmd/cache
	@echo "Done building."
	@echo "Run \"$(GOBIN)/cache\" to launch cache."

# swarm:
# 	build/env.sh go run build/ci.go install ./cmd/swarm
# 	@echo "Done building."
# 	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

# all:
# 	build/env.sh go run build/ci.go install

# android:
# 	build/env.sh go run build/ci.go aar --local
# 	@echo "Done building."
# 	@echo "Import \"$(GOBIN)/cache.aar\" to use the library."

# ios:
# 	build/env.sh go run build/ci.go xcode --local
# 	@echo "Done building."
# 	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

# test: all
# 	build/env.sh go run build/ci.go test

# lint: ## Run linters.
# 	build/env.sh go run build/ci.go lint

# clean:
# 	./build/clean_go_build_cache.sh
# 	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

# devtools:
# 	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
# 	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
# 	env GOBIN= go get -u github.com/fjl/gencodec
# 	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
# 	env GOBIN= go install ./cmd/abigen
# 	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
# 	@type "solc" 2> /dev/null || echo 'Please install solc'
# 	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# swarm-devtools:
# 	env GOBIN= go install ./cmd/swarm/mimegen

# Cross Compilation Targets (xgo)

cache-cross: cache-linux cache-darwin cache-windows cache-android cache-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/cache-*

cache-linux: cache-linux-386 cache-linux-amd64 cache-linux-arm cache-linux-mips64 cache-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-*

cache-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/cache
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep 386

cache-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/cache
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep amd64

cache-linux-arm: cache-linux-arm-5 cache-linux-arm-6 cache-linux-arm-7 cache-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep arm

cache-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/cache
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep arm-5

cache-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/cache
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep arm-6

cache-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/cache
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep arm-7

cache-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/cache
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep arm64

cache-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/cache
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep mips

cache-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/cache
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep mipsle

cache-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/cache
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep mips64

cache-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/cache
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/cache-linux-* | grep mips64le

cache-darwin: cache-darwin-386 cache-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/cache-darwin-*

cache-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/cache
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/cache-darwin-* | grep 386

cache-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/cache
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/cache-darwin-* | grep amd64

cache-windows: cache-windows-386 cache-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/cache-windows-*

cache-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/cache
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/cache-windows-* | grep 386

cache-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/cache
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/cache-windows-* | grep amd64
