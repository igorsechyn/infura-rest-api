.PHONY: build install fmt lint test test-unit clean watch test-all test-unit stop start
MAIN_PKG := ethereum-api/cmd/ethereum-api
BINARY_NAME := ethereum-api
GIT_HASH := $$(git rev-parse --short HEAD)
IMAGE_NAME := igorsechyn/ethereum-api
BIN=$(shell pwd)/bin
PERF_FOLDER=./test/performance/${shell /bin/date "+%Y-%m-%d---%H-%M-%S"}

clean:
	rm -rf build/bin/*

install:
	mkdir -p $(BIN)
	GOBIN=$(BIN) go get github.com/githubnemo/CompileDaemon@v.1.1.0
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0
	go mod download
	go mod tidy

build: clean fmt
	env GOOS=darwin GOARCH=amd64 go build -o build/bin/darwin.amd64/$(BINARY_NAME) $(GOBUILD_VERSION_ARGS) $(MAIN_PKG)
	chmod +x build/bin/darwin.amd64/$(BINARY_NAME)

fmt:
	gofmt -w=true -s $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -w=true -d $$(find . -type f -name '*.go' -not -path "./vendor/*")

lint-code:
	./bin/golangci-lint run ./... --skip-dirs vendor

test-unit:
	go test ./... -timeout 120s -count 1

test-all:
	./scripts/wait-for.sh localhost:8080 -t 5 && go test -tags integration ./... -timeout 120s -count 1

test-performance:
	mkdir -p ${PERF_FOLDER}
	ETHEREUM_API_BASE_URL=${ETHEREUM_API_BASE_URL} artillery run -o ${PERF_FOLDER}/report.json  ./test/performance/config.yaml

test: start test-all
	docker-compose logs
	docker-compose down

start: stop
	BINARY_NAME=${BINARY_NAME} MAIN_PKG=${MAIN_PKG} docker-compose up --build -d

stop: 
	docker-compose down

docker:
	docker build --build-arg binary_name=$(BINARY_NAME) --build-arg main_pkg=$(MAIN_PKG) -t $(IMAGE_NAME):$(GIT_HASH) .

docker-publish: docker
	docker push $(IMAGE_NAME):$(GIT_HASH)

watch:
	CompileDaemon -color=true -exclude-dir=.git -build="make test-unit"

verify: fmt lint-code test