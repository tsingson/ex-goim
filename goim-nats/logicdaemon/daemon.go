package logicdaemon

import (
	"context"
	"net"
	"strconv"

	"github.com/tsingson/discovery/naming"
	resolver "github.com/tsingson/discovery/naming/grpc"
	log "github.com/tsingson/zaplogger"

	"google.golang.org/grpc"

	"github.com/tsingson/ex-goim/goim-nats/logic"
	"github.com/tsingson/ex-goim/goim-nats/logic/conf"
	"github.com/tsingson/ex-goim/goim-nats/logicgrpc"
	"github.com/tsingson/ex-goim/goim-nats/logichttp"
	"github.com/tsingson/ex-goim/goim-nats/model"
	"github.com/tsingson/ex-goim/pkg/ip"
)

const (
	ver   = "2.0.0"
	appID = "goim.logic"
)

type LogicDaemon struct {
	logicServer *logic.LogicServer
	httpServer  *logichttp.Server
	grpcServer  *grpc.Server
	CancelFunc  context.CancelFunc
}

func LogicStart(cfg *conf.LogicConfig) *LogicDaemon {
	serv := &LogicDaemon{}
	ver := "2.0"
	var discovery *naming.Discovery
	{
		log.Infof("goim-logic [version: %s env: %+v] start", ver, cfg.Env)
		// grpc register naming
		discovery = naming.New(cfg.Discovery)
		resolver.Register(discovery)
	}

	// logic
	serv.logicServer = logic.New(cfg)
	// http server for push message
	serv.httpServer = logichttp.New(cfg.HTTPServer, serv.logicServer)
	// grpc server for comet client
	serv.grpcServer = logicgrpc.New(cfg.RPCServer, serv.logicServer)
	// register grpc server to discovery
	serv.CancelFunc = register(discovery, cfg)
	return serv
}

func (s *LogicDaemon) Close() {
	if s.CancelFunc != nil {
		s.CancelFunc()
	}
	s.logicServer.Close()
	s.httpServer.Close()
	// grpc
	s.grpcServer.GracefulStop()
}

func register(discovery *naming.Discovery, cfg *conf.LogicConfig) context.CancelFunc {
	env := cfg.Env
	// 	addr := "127.0.0.1" //  ip.InternalIP()
	addr := ip.InternalIP()
	_, port, _ := net.SplitHostPort(cfg.RPCServer.Addr)
	// port := "3119"
	ins := &naming.Instance{
		Region:   env.Region,
		Zone:     env.Zone,
		Env:      env.DeployEnv,
		Hostname: env.Host,
		AppID:    appID,
		Addrs: []string{
			"grpc://" + addr + ":" + port,
		},
		Metadata: map[string]string{
			model.MetaWeight: strconv.FormatInt(env.Weight, 10),
		},
	}
	cancel, err := discovery.Register(ins)
	if err != nil {
		panic(err)
	}
	return cancel
}
