# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

all: test build
setup:
	rm -rdf ./dist/
	mkdir -p ./dist/linux
	mkdir -p ./dist/windows
	mkdir -p ./dist/mac
build:
	cp ./cmd/nats/comet/comet-config.toml ./dist/linux/comet-config.toml
	cp ./cmd/nats/logic/logic-config.toml ./dist/linux/logic-config.toml
	cp ./cmd/nats/job/job-config.toml ./dist/linux/job-config.toml
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -o ./dist/linux/comet ./cmd/nats/comet/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -o ./dist/linux/logic ./cmd/nats/logic/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -o ./dist/linux/job ./cmd/nats/job/main.go

build-win:
	cp ./cmd/nats/comet/comet-config.toml ./dist/windows/comet-config.toml
	cp ./cmd/nats/logic/logic-config.toml ./dist/windows/logic-config.toml
	cp ./cmd/nats/job/job-config.toml ./dist/windows/job-config.toml
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -o ./dist/windows/comet.ext ./cmd/nats/comet/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -o ./dist/windows/logic.ext ./cmd/nats/logic/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -o ./dist/windows/job.exe ./cmd/nats/job/main.go

build-mac:
	cp ./cmd/nats/comet/comet-config.toml ./dist/mac/comet-config.toml
	cp ./cmd/nats/logic/logic-config.toml ./dist/mac/logic-config.toml
	cp ./cmd/nats/job/job-config.toml ./dist/mac/job-config.toml
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin  go build -o ./dist/mac/comet ./cmd/nats/comet/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin  go build -o ./dist/mac/logic ./cmd/nats/logic/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin  go build -o ./dist/mac/job ./cmd/nats/job/main.go



test:
	$(GOTEST) -v ./...

clean:
	rm -rf dist/

run:
	nohup ./dist/logic   2>&1 > dist/logic.log &
	nohup ./dist/comet   2>&1 > dist/comet.log &
	nohup ./dist/job   2>&1 > dist/job.log &

stop:
	pkill -f ./dist/logic
	pkill -f ./dist/job
	pkill -f ./dist/comet
