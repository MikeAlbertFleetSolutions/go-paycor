.PHONY: default get codetest test fmt lint vet cyclo secure

default: fmt codetest

get:
	go get -t -v ./...
	go get -u github.com/tcnksm/ghr
	go get -u golang.org/x/lint/golint
	go get -u github.com/fzipp/gocyclo
	go get -u github.com/securego/gosec/cmd/gosec/...

codetest: lint vet cyclo secure

test:
	go test -v -cover

fmt:
	go fmt ./...

lint:
	@echo golint ./...
	@OUTPUT=`golint ./... 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "golint errors:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	fi

vet:
	go vet -all .

cyclo:
	gocyclo -over 20 .

secure:
	gosec -quiet ./...
