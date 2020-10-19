.PHONY: coverage.txt

GOFILES = $(shell find ./*.go | grep -v _test)
TESTFILES = $(shell find ./*.go | grep -v cmd/)
TEST = go test -v -failfast -cover $(TESTFILES)

all: install

bench:
	cd cmd/teleporter/ && go test -bench .

cover: coverage.txt
	go tool cover -html=coverage.txt

coverage.txt:
	$(TEST) -coverprofile=$@ -covermode=atomic

dasm:
	go run cmd/dasm/main.go

docs:
	@go doc

gosec:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b ~/bin v2.3.0

install: gosec
	go get -u golang.org/x/lint/golint
	go get github.com/fzipp/gocyclo/cmd/gocyclo

lint:
	go fmt -x
	gocyclo -over 10 .
	golint
	go vet
	gosec ./...

teleporter:
	go run cmd/teleporter/main.go

test:
	$(TEST)

test-run:
ifdef test
	$(TEST) -run $(test)
else
	@echo Syntax is 'make $@ test=<test name>'
endif

vault:
	go run cmd/vault/main.go

vm:
	go run cmd/vm/main.go 2> vm.log
