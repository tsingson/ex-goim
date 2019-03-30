# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

all: test build
build:
	rm -rdf ./target/
	mkdir ./target/
	cp ./cmd/nats/comet/comet-config.toml ./target/comet-config.toml
	cp ./cmd/nats/logic/logic-config.toml ./target/logic-config.toml
	cp ./cmd/nats/job/job-config.toml ./target/job-config.toml
	$(GOBUILD) -o ./target/comet ./cmd/nats/comet/main.go
	$(GOBUILD) -o ./target/logic ./cmd/nats/logic/main.go
	$(GOBUILD) -o ./target/job ./cmd/nats/job/main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -rf target/

run:
	nohup ./target/logic   2>&1 > target/logic.log &
	nohup ./target/comet   2>&1 > target/comet.log &
	nohup ./target/job   2>&1 > target/job.log &

stop:
	pkill -f ./target/logic
	pkill -f ./target/job
	pkill -f ./target/comet
