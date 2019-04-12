package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/ex-goim/goim-nats/logic/conf"
	"github.com/tsingson/ex-goim/goim-nats/logicdaemon"
)

const (
	ver   = "2.0.0"
	appID = "goim.logic"
)

var cfg *conf.LogicConfig

func main() {
	cfg = conf.Default()

	// _ = file.SaveToml(cfg, "/Users/qinshen/go/bin/logic-config.toml")
	daemon := logicdaemon.LogicStart(cfg)
	// signal
	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			s := <-c
			log.Infof("goim-logic get a signal %s", s.String())
			//
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				daemon.Close()
				log.Infof("goim-logic [version: %s] exit", ver)
				// log.Flush()
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	}
}
