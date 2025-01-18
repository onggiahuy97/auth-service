TIMEOUT := 300s


test-database:
	go test -v ./database/...

lint:
	staticcheck ./... ; go fmt ./... ; gosec ./... ;	golangci-lint run ./...


