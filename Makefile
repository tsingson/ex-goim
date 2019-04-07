GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

all: setup  build
setup:
	rm -rdf ./dist/
	mkdir -p ./dist/linux
	mkdir -p ./dist/windows
	mkdir -p ./dist/mac
	mkdir -p ./dist/log
build-linux:
	cp ./cmd/nats/comet/comet-config.toml ./dist/linux/comet-config.toml
	cp ./cmd/nats/logic/logic-config.toml ./dist/linux/logic-config.toml
	cp ./cmd/nats/job/job-config.toml ./dist/linux/job-config.toml
	cp ./third-party/discoveryd/discoveryd-config.toml ./dist/linux/

	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/linux/gnatsd ./third-party/gnatsd/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/linux/discoveryd ./third-party/discoveryd/
    ${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/linux/liftbridge ./third-party/liftbridge/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/linux/comet ./cmd/nats/comet/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/linux/logic ./cmd/nats/logic/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/linux/job ./cmd/nats/job/main.go

build-win:
	cp ./cmd/nats/comet/comet-config.toml ./dist/windows/comet-config.toml
	cp ./cmd/nats/logic/logic-config.toml ./dist/windows/logic-config.toml
	cp ./cmd/nats/job/job-config.toml ./dist/windows/job-config.toml
	cp ./third-party/discoveryd/discoveryd-config.toml ./dist/windows/

	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/windows/gnatsd ./third-party/gnatsd/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/windows/discoveryd ./third-party/discoveryd/
    ${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/windows/liftbridge ./third-party/liftbridge/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/windows/comet.exe ./cmd/nats/comet/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/windows/logic.exe ./cmd/nats/logic/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/windows/job.exe ./cmd/nats/job/main.go

build-mac:
	cp ./cmd/nats/comet/comet-config.toml ./dist/mac/comet-config.toml
	cp ./cmd/nats/logic/logic-config.toml ./dist/mac/logic-config.toml
	cp ./cmd/nats/job/job-config.toml ./dist/mac/job-config.toml
	cp ./third-party/discoveryd/discoveryd-config.toml ./dist/mac/

	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/mac/gnatsd ./third-party/gnatsd/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/mac/discoveryd ./third-party/discoveryd/
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/mac/liftbridge ./third-party/liftbridge/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin  go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/mac/comet ./cmd/nats/comet/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin  go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/mac/logic ./cmd/nats/logic/main.go
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin  go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./dist/mac/job ./cmd/nats/job/main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -rf dist/

run-linux:
	nohup ./dist/linux/gnatsd   2>&1 > dist/log/gnatsd.log &
	nohup ./dist/linux/liftbridge --raft-bootstrap-seed   2>&1 > dist/log/liftbridge.log &
	nohup ./dist/linux/discoveryd   2>&1 > dist/log/discoveryd.log &
	nohup ./dist/linux/logic   2>&1 > dist/log/logic.log &
	nohup ./dist/linux/comet   2>&1 > dist/log/comet.log &
	nohup ./dist/linux/job   2>&1 > dist/log/job.log &

run-mac:
	nohup ./dist/mac/gnatsd   2>&1 > dist/log/gnatsd.log &
	nohup ./dist/mac/liftbridge --raft-bootstrap-seed   2>&1 > dist/log/liftbridge.log &
	nohup ./dist/mac/discoveryd   2>&1 > dist/log/discoveryd.log &
	nohup ./dist/mac/logic   2>&1 > dist/log/logic.log &
	nohup ./dist/mac/comet   2>&1 > dist/log/comet.log &
	nohup ./dist/mac/job   2>&1 > dist/log/job.log &

stop:
	pkill -f ./dist/linux/gnatsd
	pkill -f ./dist/linux/liftbridge
	pkill -f ./dist/linux/discoveryd
	pkill -f ./dist/linux/logic
	pkill -f ./dist/linux/job
	pkill -f ./dist/linux/comet

