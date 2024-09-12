server:
	go run cmd/api/main.go

test:
	go test -v -cover -short ./...

.PHONY: server test
