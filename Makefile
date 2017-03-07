# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

CURDIR := $(shell pwd)
GO        := go
GOBUILD   := $(GO) build
GOTEST    := $(GO) test

OS        := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
PACKAGES  := $$(go list ./...| grep -vE 'vendor')
FILES     := $$(find . -name '*.go' | grep -vE 'vendor')
TARGET	  := "server"

test:
	$(GOTEST) $(PACKAGES) -cover

build:
	$(GOBUILD) -o $(TARGET)

dev: test build

clean:
	rm $(TARGET)

update:
	which glide >/dev/null || curl https://glide.sh/get | sh
	which glide-vc || go get -v -u github.com/sgotti/glide-vc
	glide update
	@echo "removing test files"
	glide-vc --only-code --no-tests