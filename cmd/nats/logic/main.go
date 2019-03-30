package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sanity-io/litter"
	"github.com/tsingson/discovery/naming"
	resolver "github.com/tsingson/discovery/naming/grpc"
	"github.com/tsingson/fastx/utils"
	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/goim/internal/nats/logic"
	"github.com/tsingson/goim/internal/nats/model"

	"github.com/tsingson/goim/internal/nats/logic/grpc"
	"github.com/tsingson/goim/internal/nats/logic/http"

	"github.com/tsingson/goim/internal/nats/logic/conf"
)

const (
	ver   = "2.0.0"
	appid = "goim.logic"
)

var cfg *conf.LogicConfig

func main() {
	path, _ := utils.GetCurrentExecDir()
	confPath := path + "/logic-config.toml"

	var err error
	cfg, err = conf.LoadToml(confPath)

	if err != nil {
		panic(err)
	}

	litter.Dump(cfg)

	{
		// env := &conf.Env{
		// 	Region:    "china",
		// 	Zone:      "gd",
		// 	DeployEnv: "sz",
		// 	Host:      "logic",
		// }
		// cfg.Env = env
	}

	var dis *naming.Discovery

	{
		log.Infof("goim-logic [version: %s env: %+v] start", ver, cfg.Env)
		// grpc register naming
		dis = naming.New(cfg.Discovery)
		resolver.Register(dis)
	}

	// logic
	srv := logic.New(cfg)
	httpSrv := http.New(cfg.HTTPServer, srv)
	rpcSrv := grpc.New(cfg.RPCServer, srv)
	cancel := register(dis, srv)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("goim-logic get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			srv.Close()
			httpSrv.Close()
			rpcSrv.GracefulStop()
			log.Infof("goim-logic [version: %s] exit", ver)
			// log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func register(dis *naming.Discovery, srv *logic.Logic) context.CancelFunc {
	env := cfg.Env
	addr := "127.0.0.1" //  ip.InternalIP()
	_, port, _ := net.SplitHostPort(cfg.RPCServer.Addr)
	// port := "3119"
	ins := &naming.Instance{
		Region:   env.Region,
		Zone:     env.Zone,
		Env:      env.DeployEnv,
		Hostname: env.Host,
		AppID:    appid,
		Addrs: []string{
			"grpc://" + addr + ":" + port,
		},
		Metadata: map[string]string{
			model.MetaWeight: strconv.FormatInt(env.Weight, 10),
		},
	}
	cancel, err := dis.Register(ins)
	if err != nil {
		panic(err)
	}
	return cancel
}
