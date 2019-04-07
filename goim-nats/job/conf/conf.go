package conf

import (
	"time"

	"github.com/tsingson/discovery/naming"

	xtime "github.com/tsingson/ex-goim/pkg/time"
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
