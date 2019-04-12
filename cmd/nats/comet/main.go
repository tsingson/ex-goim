package main

import (
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/ex-goim/goim-nats/comet/conf"
	"github.com/tsingson/ex-goim/goim-nats/cometdaemon"
)

const (
	ver   = "2.0.0"
	appid = "goim.comet"
)

var (
	cfg *conf.CometConfig
)

func main() {

	// path, _ := utils.GetCurrentExecDir()
	// confPath := path + "/comet-config.toml"
	//
	// flag.Parse()
	// var err error
	// cfg, err = conf.Load(confPath)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// cfg.Env = &conf.Env{
	// 	Region:    "china",
	// 	Zone:      "gd",
	// 	DeployEnv: "dev",
	// 	Host:      "comet",
	// }

	cfg = conf.Default()
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	println(cfg.Debug)

	log.Infof("goim-comet [version: %s env: %+v] start", ver, cfg.Env)
	ss, err := cometdaemon.CometStart(cfg)
	if err != nil {
		panic(err)
	}
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("goim-comet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ss.Close()
			log.Infof("goim-comet [version: %s] exit", ver)
			// log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
