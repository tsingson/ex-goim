package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/tsingson/discovery/naming"
	discovery "github.com/tsingson/discovery/naming/grpc"
	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/ex-goim/goim-nats/comet"
	"github.com/tsingson/ex-goim/goim-nats/comet/conf"
	"github.com/tsingson/ex-goim/goim-nats/comet/grpc"
	"github.com/tsingson/ex-goim/goim-nats/model"
	"github.com/tsingson/ex-goim/pkg/ip"
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
	// register discovery
	dis := naming.New(cfg.Discovery)
	discovery.Register(dis)
	// new comet server
	srv := comet.NewServer(cfg)
	if err := comet.InitWhitelist(cfg.Whitelist); err != nil {
		panic(err)
	}
	// tcp
	if err := comet.InitTCP(srv, cfg.TCP.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}
	//
	if err := comet.InitWebsocket(srv, cfg.Websocket.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}
	//
	if cfg.Websocket.TLSOpen {
		if err := comet.InitWebsocketWithTLS(srv, cfg.Websocket.TLSBind, cfg.Websocket.CertFile, cfg.Websocket.PrivateFile, runtime.NumCPU()); err != nil {
			panic(err)
		}
	}
	// new grpc server
	rpcSrv := grpc.New(cfg.RPCServer, srv)
	cancel := register(dis, srv)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("goim-comet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			rpcSrv.GracefulStop()
			srv.Close()
			log.Infof("goim-comet [version: %s] exit", ver)
			// log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func register(dis *naming.Discovery, srv *comet.Server) context.CancelFunc {
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
	cancel, err := dis.Register(ins)
	if err != nil {
		panic(err)
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
	return cancel
}
