GOFILES = $(shell find ./*.go | grep -v _test)
TEST = go test -v -failfast -cover ./...

all: install

docs:
	@go doc

gosec:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b ~/bin v2.3.0

install:
	go get -u golang.org/x/lint/golint
	go get github.com/fzipp/gocyclo
	go get github.com/securego/gosec/cmd/gosec/...

lint:
	go fmt -x
	gocyclo -over 10 .
	golint
	go vet
	gosec ./...

run:
	go run $(GOFILES)

test:
	$(TEST)

test-run:
ifdef test
	$(TEST) -run $(test)
else
	@echo Syntax is 'make $@ test=<test name>'
endif
