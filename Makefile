all: install

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
