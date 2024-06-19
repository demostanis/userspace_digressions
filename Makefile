MAKEFLAGS += -s

GO = go
GOFLAGS = -tags netgo

all: build
	./vm.sh

build:
	$(GO) build $(GOFLAGS) cmd/init/init.go
	$(GO) build $(GOFLAGS) cmd/initctl/initctl.go

clean:
	$(RM) alpine.iso boot/

re: clean all

.PHONY: all clean re
