BUILD_DATE = $(shell date +%F-%Z/%T)
TAG=$(shell git describe --tags --abbrev=0)
GO_VERSION=$(shell go version)
GO_VERS := $(shell echo "${GO_VERSION}" | sed -e 's/ /_/g')
build:
	clear
	rm -rdf ./bin
	go fmt ./...
	# go build -ldflags "-w -s -X main.GoVersion=$(GO_VERS) -X main.Version=${TAG} -X main.Date=${BUILD_DATE}" -o ./bin/server ./cmd/server/main.go
	go build -o ./bin/wartank_prod ./cmd/server/main.go
	strip -s ./bin/wartank_prod
	upx -f ./bin/wartank_prod
prod2:
	clear
	rm -rdf ./bin_dev
	go fmt ./...
	go build -ldflags "-X main.GoVersion=${GO_VERS} -X main.Version=${TAG} -X main.Date=${BUILD_DATE}" -o ./bin2/server2 ./cmd/server/main.go
	strip -s ./bin2/server2
	upx -f ./bin2/server2
dev:
	clear
	rm -rdf ./bin_dev
	go fmt ./...
	# go build -ldflags "-w -s -X main.GoVersion=$(GO_VERS) -X main.Version=${TAG} -X main.Date=${BUILD_DATE}" -o ./bin_dev/wartank_dev ./cmd/server/main.go
	go build -race -o ./bin_dev/wartank_dev ./cmd/server/main.go
	./run_dev.sh
prod:
	clear
	go fmt ./...
	go build -o ./bin/wartank_prod ./cmd/server/main.go
	./run_prod.sh
.PHONY: test
test:
	clear
	go fmt ./...
	go test -shuffle=on -vet=all -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func=cover.out
mod:
	clear
	go get -u ./...
	go mod tidy -compat=1.22.0
	go mod vendor
	go fmt ./...
lint:
	clear
	go fmt ./...
	golangci-lint run ./cmd/server/...
	golangci-lint run ./app/...
	golangci-lint run ./kernel/...
.PHONY: graph
graph:
	# #go install github.com/kisielk/godepgraph@latest
	# #go mod graph | modgraphviz | dot -Tsvg -o graph.svg
	# godepgraph ./cmd/server | dot -Tsvg -o graph.svg
	clear
	go fmt ./...
	# http://localhost:7878/
	# go-callvis ./...
	#
	#
	godepgraph ./cmd/server | dot -Tsvg -o graph.svg
