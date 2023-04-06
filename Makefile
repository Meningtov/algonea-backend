
install:
	go mod download

test:
	bash -c 'diff -u <(echo -n) <(gofmt -s -d .)'
	go vet ./...
	go test -race ./...

.PHONY: install test
