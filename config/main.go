package main

import (
	"github.com/imdario/mergo"
	"github.com/sanity-io/litter"
	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/ex-goim/goim-nats/nats/comet/conf"
)

type networkConfig struct {
	Protocol   string
	Address    string
	ServerType string `json: "server_type"`
	Port       uint16
}

type FssnConfig struct {
	Network networkConfig
}

type merge struct {
	FssnConfig FssnConfig
	conf.CometConfig
}

func main() {

	var fssnDefault = FssnConfig{
		networkConfig{
			"tcp",
			"127.0.0.1",
			"http",
			31560,
		},
	}
	confPath := "/Users/qinshen/go/bin/comet-config.toml"
	con, _ := conf.Init(confPath)

	config := merge{
		fssnDefault, *con,
	}
	cc := merge{
		fssnDefault, *con,
	}

	if err := mergo.Merge(&config, cc, mergo.WithOverride); err != nil {
		log.Fatal(err)
	}

	litter.Dump(config.RPCClient)

}
