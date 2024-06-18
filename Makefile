MAKEFLAGS += -s

GO = go

all:
	$(GO) build init.go
	./vm.sh

clean:
	$(RM) -r alpine.iso boot/

re: clean all

.PHONY: all clean re
