server:
	go run cmd/api/main.go

test:
	go test -v -cover -short ./...

air:
	air -c .air.toml
	#air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

.PHONY: server test
