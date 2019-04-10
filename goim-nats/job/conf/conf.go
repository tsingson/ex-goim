package conf

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
	"github.com/tsingson/discovery/naming"
	"golang.org/x/xerrors"

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
			Channel:   "goim-push-topic",
			ChannelID: "goim-push-topic-stream",
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

// Load init config.
func Load(path string) (cfg *Config, err error) {

	if len(path) == 0 {
		return cfg, xerrors.New("config path is nil")
	}

	Conf = Default()
	cfg = Default()

	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		return
	}
	err = mergo.Merge(&Conf, cfg, mergo.WithOverride)
	if err != nil {
		return Conf, err
	}

	return Conf, nil
}
