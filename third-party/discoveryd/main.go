package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/tsingson/go-daemon"
	log "github.com/tsingson/zaplogger"
	"go.uber.org/zap"

	"github.com/tsingson/discovery/discovery"
	"github.com/tsingson/discovery/http"
)

func main() {

	runtime.MemProfileRate = 0
	runtime.GOMAXPROCS(128)


	/**
	tw = timingwheel.NewTimingWheel(time.Minute, 60)
	tw.StartCron()
	defer tw.StopCron()
	*/

	var cntxt = &daemon.Context{
		PidFileName: "pid-discoveryd",
		PidFilePerm: 0644,
		LogFileName: logPath + "/discoveryd-daemon.log",
		LogFilePerm: 0640,
		WorkDir:     path,
		Umask:       027,
		// 	Args:        []string{"aaa-demo"},
	}

	var d, err = cntxt.Reborn()
	if err != nil {
		log.Fatal("cat's reborn ", zap.Error(err))
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Info("trying to start daemon")

	svr, cancel := discovery.New(cfg)

	 http.Init(cfg, svr)

	// if err != nil {
	// 	cancel()
	// 	os.Exit(-1)
	// }
	runtime.Goexit()
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("discovery get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancel()
			time.Sleep(time.Second)
			log.Info("discovery quit !!!")
			// log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
