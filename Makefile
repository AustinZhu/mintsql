mod:
	go get -u ./...
	go mod tidy

server:
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o ./build/mintsql ./cmd/server/main.go

client:
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o ./build/mintcli ./cmd/client/main.go

test:
	go test -v ./... -cover

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client

image:
	docker build -t mintsql:latest .

docker-run:
	docker run -d --rm -p 7384:7384 mintsql:latest