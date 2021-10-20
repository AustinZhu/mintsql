.PHONY: mod
mod:
	go get -u ./...
	go mod tidy

.PHONY: server
server:
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o ./build/mintsql ./cmd/server/main.go

.PHONY: client
client:
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o ./build/mintcli ./cmd/client/main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: run-server
run-server:
	go run ./cmd/server

.PHONY: run-client
run-client:
	go run ./cmd/client

.PHONY: image
image:
	docker build -t mintsql:latest .

.PHONY: docker-run
docker-run:
	docker run -d --rm -p 7384:7384 mintsql:latest