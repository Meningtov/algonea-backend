
install:
	go mod download

test:
	go test -v handler/*_test.go

.PHONY: install test
