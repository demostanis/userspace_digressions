MAKEFLAGS += -s

GO = go

all:
	$(GO) build init.go
	./vm.sh

clean:
	$(RM) alpine.iso boot/

re: clean all

.PHONY: all clean re
