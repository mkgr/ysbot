NAME:=$(shell go list .)
BINDIR:=bin
TARGET:=$(BINDIR)/$(NAME)
FILES:=$(shell find . -type f -name '*.go' -print)

_GOOS:=linux
_GOARCH:=amd64

.PHONY: build
.PHONY: clean
.PHONY: fmt

build: $(TARGET)
clean:
	rm -f $(TARGET)
fmt:
	go fmt *.go

$(TARGET): $(FILES)
	env GOOS=$(_GOOS) GOARCH=$(_GOARCH) go build -o $@

