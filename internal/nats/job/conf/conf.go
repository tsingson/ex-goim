package conf

import (
	"time"

	"github.com/tsingson/discovery/naming"

	xtime "github.com/tsingson/goim/pkg/time"
)

// Config is job config.
type Config struct {
	Env       *Env
	Nats      *Nats
	Discovery *naming.Config
	Comet     *Comet
	Room      *Room
}

// JobConfig is alias of Config
type JobConfig = Config

// Room is room config.
type Room struct {
	Batch  int
	Signal xtime.Duration
	Idle   xtime.Duration
}

// Comet is comet config.
type Comet struct {
	RoutineChan int
	RoutineSize int
}

// Kafka is kafka config.
type Kafka struct {
	Topic   string
	Group   string
	Brokers []string
}

// Env is env config.
type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
}

type Nats struct {
	Channel   string
	ChannelID string
	Group     string
	NatsAddr  string
	LiftAddr  string
}

type NatsConfig = Nats

var (
	confPath  string
	region    string
	zone      string
	deployEnv string
	host      string

	Conf *JobConfig
)

// func init() {
// 	var (
// 		defHost, _ = os.Hostname()
// 	)
// 	flag.StringVar(&confPath, "conf", "job-example.toml", "default config path")
// 	flag.StringVar(&region, "region", os.Getenv("REGION"), "avaliable region. or use REGION env variable, value: sh etc.")
// 	flag.StringVar(&zone, "zone", os.Getenv("ZONE"), "avaliable zone. or use ZONE env variable, value: sh001/sh002 etc.")
// 	flag.StringVar(&deployEnv, "deploy.env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
// 	flag.StringVar(&host, "host", defHost, "machine hostname. or use default machine hostname.")
// }

// // Init init config.
// func Init(path string) (cfg *JobConfig, err error) {
// 	cfg = Default()
// 	if len(path) > 0 {
// 		_, err = toml.DecodeFile(path, &cfg)
// 	} else {
// 		_, err = toml.DecodeFile(confPath, &cfg)
// 	}
//
// 	return
// }

func init() {
	Conf = Default()
}

// Default new a config with specified defualt value.
func Default() *Config {

	return &Config{
		Nats: &Nats{
			Channel:   "channel",
			ChannelID: "channel-stream",
			Group:     "group",
			LiftAddr:  "localhost:9292", // address for lift-bridge
			NatsAddr:  "localhost:4222",
		},
		Env: &Env{
			Region:    "china",
			Zone:      "gd",
			DeployEnv: "dev",
			Host:      "job",
		},
		Discovery: &naming.Config{
			Nodes:  []string{"127.0.0.1:7171"},
			Region: "china",
			Zone:   "gd",
			Env:    "dev",
			Host:   "discovery",
		},
		Comet: &Comet{
			RoutineChan: 1024,
			RoutineSize: 32,
		},
		Room: &Room{
			Batch:  20,
			Signal: xtime.Duration(time.Second),
			Idle:   xtime.Duration(time.Minute * 15),
		},
	}
}
