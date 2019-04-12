package job

import (
	"github.com/tsingson/discovery/naming"
	resolver "github.com/tsingson/discovery/naming/grpc"

	"github.com/tsingson/ex-goim/goim-nats/job/conf"
)

func JobStart(cfg *conf.JobConfig) *NatsJob {
	// gRPC register naming
	discovery := naming.New(cfg.Discovery)
	resolver.Register(discovery)
	// job
	jobServer := New(cfg)
	// run jobServer
	go jobServer.Consume()
	return jobServer
}
