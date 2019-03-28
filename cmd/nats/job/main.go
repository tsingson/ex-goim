package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/tsingson/discovery/naming"
	"github.com/tsingson/fastx/utils"

	"github.com/tsingson/goim/internal/nats/job/conf"

	resolver "github.com/tsingson/discovery/naming/grpc"
	log "github.com/tsingson/zaplogger"
)

var (
	ver = "2.0.0"
	cfg *conf.JobConfig
)

func main() {

	path, _ := utils.GetCurrentExecDir()
	confPath := path + "/job-config.toml"
	flag.Parse()

	cfg, err := conf.Init(confPath)
	if err != nil {
		panic(err)
	}

	env := &conf.Env{
		Region:    "test",
		Zone:      "test",
		DeployEnv: "test",
		Host:      "localhost",
	}
	cfg.Env = env

	log.Infof("goim-job [version: %s env: %+v] start", ver, cfg.Env)
	// grpc register naming
	dis := naming.New(cfg.Discovery)
	resolver.Register(dis)
	// job
	j := natsjob.New(cfg)
	go j.Consume()
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("goim-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			j.Close()
			log.Infof("goim-job [version: %s] exit", ver)
			// log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
