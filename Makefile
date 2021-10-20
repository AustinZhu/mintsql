.PHONY: mod
mod:
	go get -u ./...
	go mod tidy
	go mod vendor

.PHONY: build-server
build-server: mod
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o ./build/mintsql ./cmd/server

.PHONY: build-client
build-client: mod
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o ./build/mintcli ./cmd/client

.PHONY: install-server
install-server: mod
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o $(GOPATH)/bin/mintsql ./cmd/server

.PHONY: install-client
install-client: mod
	CGO_ENABLED=0 go build -mod=vendor -ldflags '-w -s' -a -o $(GOPATH)/bin/mintcli ./cmd/client

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: run-server
run-server:
	go run ./cmd/server $(PORT)

.PHONY: run-client
run-client:
	go run ./cmd/client $(HOST) $(PORT)

.PHONY: image
image:
	docker build -t mintsql:latest .

.PHONY: docker-run
docker-run:
	docker run -d --rm -p 7384:7384 austinzhu666/mintsql:latest