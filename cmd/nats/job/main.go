package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/ex-goim/goim-nats/job"
	"github.com/tsingson/ex-goim/goim-nats/job/conf"
)

var (
	ver = "2.0.0"
	cfg *conf.JobConfig
)

func main() {

	// path, _ := utils.GetCurrentExecDir()
	// confPath := path + "/job-config.toml"
	// flag.Parse()
	// var err error
	// cfg, err = conf.Load(confPath)
	// if err != nil {
	// 	panic(err)
	// }
	cfg = conf.Default()

	log.Infof("goim-job [version: %s env: %+v] start", ver, cfg.Env)
	jobServer := job.JobStart(cfg)
	// signal
	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			s := <-c
			log.Infof("goim-job get a signal %s", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				_ = jobServer.Close()
				log.Infof("goim-job [version: %s] exit", ver)
				// log.Flush()
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	}
}
