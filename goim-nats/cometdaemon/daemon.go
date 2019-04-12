package cometdaemon

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/tsingson/discovery/naming"
	discovery "github.com/tsingson/discovery/naming/grpc"
	log "github.com/tsingson/zaplogger"
	"google.golang.org/grpc"

	"github.com/tsingson/ex-goim/goim-nats/comet"
	"github.com/tsingson/ex-goim/goim-nats/comet/conf"
	"github.com/tsingson/ex-goim/goim-nats/cometgrpc"
	"github.com/tsingson/ex-goim/goim-nats/model"
	"github.com/tsingson/ex-goim/pkg/ip"
)

const (
	ver   = "2.0.0"
	appid = "goim.comet"
)

type CometDaemon struct {
	cometServer *comet.CometServer
	grpcServer  *grpc.Server
	CancelFunc  context.CancelFunc
}

func CometStart(cfg *conf.CometConfig) (ss *CometDaemon, err error) {

	// register discovery
	dis := naming.New(cfg.Discovery)
	discovery.Register(dis)
	// new comet server
	srv := comet.New(cfg)
	ss.cometServer = srv
	if err = comet.InitWhitelist(cfg.Whitelist); err != nil {
		return
	}

	// tcp
	if err = comet.InitTCP(srv, cfg.TCP.Bind, runtime.NumCPU()); err != nil {
		return
	}
	//
	if err = comet.InitWebsocket(srv, cfg.Websocket.Bind, runtime.NumCPU()); err != nil {
		return
	}
	//
	if cfg.Websocket.TLSOpen {
		if err = comet.InitWebsocketWithTLS(srv, cfg.Websocket.TLSBind, cfg.Websocket.CertFile, cfg.Websocket.PrivateFile, runtime.NumCPU()); err != nil {
			return
		}
	}
	// new grpc server
	rpcSrv := cometgrpc.New(cfg.RPCServer, srv)
	ss.grpcServer = rpcSrv
	var cancelFunc context.CancelFunc
	cancelFunc, err = register(dis, srv, cfg)
	if err != nil {
		return
	}
	ss.CancelFunc = cancelFunc
	return
}

func (s *CometDaemon) Close() {
	if s.CancelFunc != nil {
		s.CancelFunc()
	}
	s.grpcServer.GracefulStop()
	s.cometServer.Close()
}

func register(dis *naming.Discovery, srv *comet.Server, cfg *conf.Config) (cancelFunc context.CancelFunc, err error) {
	env := cfg.Env
	addr := ip.InternalIP()
	_, port, _ := net.SplitHostPort(cfg.RPCServer.Addr)
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
			model.MetaWeight:  strconv.FormatInt(env.Weight, 10),
			model.MetaOffline: strconv.FormatBool(env.Offline),
			model.MetaAddrs:   strings.Join(env.Addrs, ","),
		},
	}
	cancelFunc, err = dis.Register(ins)
	if err != nil {
		return
	}
	// renew discovery metadata
	go func() {
		for {
			var (
				err   error
				conns int
				ips   = make(map[string]struct{})
			)
			for _, bucket := range srv.Buckets() {
				for ip := range bucket.IPCount() {
					ips[ip] = struct{}{}
				}
				conns += bucket.ChannelCount()
			}
			ins.Metadata[model.MetaConnCount] = fmt.Sprint(conns)
			ins.Metadata[model.MetaIPCount] = fmt.Sprint(len(ips))
			if err = dis.Set(ins); err != nil {
				log.Errorf("dis.Set(%+v) error(%v)", ins, err)
				time.Sleep(time.Second)
				continue
			}
			time.Sleep(time.Second * 10)
		}
	}()
	return
}
