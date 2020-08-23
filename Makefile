VER ?= 1.0
# version info
git_state  := $(shell (git status --porcelain | grep -q .) && echo \(dirty\))
commit_id  := $(shell git rev-parse HEAD)
build_time := $(shell date +%FT%T%z)

work_dir   := $(shell pwd)
ifneq "$(strip $(shell command -v go 2>/dev/null))" ""
	GOOS ?= $(shell go env GOOS)
	GOARCH ?= $(shell go env GOARCH)
else
	ifeq ($(GOOS),)
		# approximate GOOS for the platform if we don't have Go and GOOS isn't
		# set. We leave GOARCH unset, so that may need to be fixed.
		ifeq ($(OS),Windows_NT)
			GOOS = windows
		else
			UNAME_S := $(shell uname -s)
			ifeq ($(UNAME_S),Linux)
				GOOS = linux
			endif
			ifeq ($(UNAME_S),Darwin)
				GOOS = darwin
			endif
			ifeq ($(UNAME_S),FreeBSD)
				GOOS = freebsd
			endif
		endif
	else
		GOOS ?= $$GOOS
		GOARCH ?= $$GOARCH
	endif
endif

all: kubectl-inspection

kubectl-inspection:
	mkdir -p bin
	export GOOS=$(GOOS) GO111MODULE="on"; go build -o bin/kubectl-inspection -ldflags "-X main.commitId=${commit_id}${git_state} -X main.buildTime=${build_time}" cmd/kubectl-inspection.go

clean:
	rm -rf bin

