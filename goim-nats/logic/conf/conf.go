package conf

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
	"github.com/tsingson/discovery/naming"
	"golang.org/x/xerrors"

	xtime "github.com/tsingson/ex-goim/pkg/time"
)

// Config config.
type Config struct {
	Env        *Env
	Discovery  *naming.Config
	RPCClient  *RPCClient
	RPCServer  *RPCServer
	HTTPServer *HTTPServer
	// Kafka      *Kafka
	Nats    *Nats
	Redis   *Redis
	Node    *Node
	Backoff *Backoff
	Regions map[string][]string
}

// LogicConfig as alias of Config
type LogicConfig = Config

// Env is env config.
type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
	Weight    int64
}

// Node node config.
type Node struct {
	DefaultDomain string
	HostDomain    string
	TCPPort       int
	WSPort        int
	WSSPort       int
	HeartbeatMax  int
	Heartbeat     xtime.Duration
	RegionWeight  float64
}

// Backoff backoff.
type Backoff struct {
	MaxDelay  int32
	BaseDelay int32
	Factor    float32
	Jitter    float32
}

// Redis .
type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
	Expire       xtime.Duration
}

// Kafka .
type Kafka struct {
	Topic   string
	Brokers []string
}

// Nats .
type Nats struct {
	NatsAddr  string // "nats://localhost:4222"
	LiftAddr  string // "localhost:9292" // address for lift-bridge
	Channel   string //  "channel"
	ChannelID string //  "channel-stream"
	AckInbox  string // "acks"
}

// NatsConfig as alias of Nats
type NatsConfig = Nats

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    xtime.Duration
	Timeout xtime.Duration
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string
	Addr              string
	Timeout           xtime.Duration
	IdleTimeout       xtime.Duration
	MaxLifeTime       xtime.Duration
	ForceCloseWait    xtime.Duration
	KeepAliveInterval xtime.Duration
	KeepAliveTimeout  xtime.Duration
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

var (
	confPath  string
	region    string
	zone      string
	deployEnv string
	host      string
	weight    int64

	// Conf config
	Conf *Config
)

func init() {
	Conf = Default()

}

// Load init config.
func Load(path string) (cfg *Config, err error) {

	if len(path) == 0 {
		return cfg, xerrors.New("config path is nil")
	}

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

// LoadToml init config.
func LoadToml(path string) (cfg *Config, err error) {
	Conf = Default()
	if len(path) == 0 {
		return Conf, xerrors.New("no configuration")
	}
	//
	_, err = toml.DecodeFile(path, &Conf)

	return Conf, nil
}

// Default new a config with specified defualt value.
func Default() *LogicConfig {
	cfg := &LogicConfig{
		Env: &Env{
			Region:    "china",
			Zone:      "gd",
			DeployEnv: "dev",
			Host:      "logic",
			Weight:    100,
		},
		Discovery: &naming.Config{
			Nodes:  []string{"127.0.0.1:7171"},
			Region: "china",
			Zone:   "gd",
			Env:    "dev",
			Host:   "discovery",
		},
		Nats: &Nats{
			NatsAddr:  "nats://localhost:4222",
			LiftAddr:  "localhost:9292", // address for lift-bridge
			Channel:   "goim-push-topic",
			ChannelID: "goim-push-topic-stream",
			AckInbox:  "acks",
		},
		HTTPServer: &HTTPServer{
			Network:      "tcp",
			Addr:         "127.0.0.1:3111",
			ReadTimeout:  xtime.Duration(time.Second),
			WriteTimeout: xtime.Duration(time.Second),
		},
		RPCClient: &RPCClient{
			Dial:    xtime.Duration(time.Second),
			Timeout: xtime.Duration(time.Second),
		},
		RPCServer: &RPCServer{
			Network:           "tcp",
			Addr:              "127.0.0.1:3119",
			Timeout:           xtime.Duration(time.Second),
			IdleTimeout:       xtime.Duration(time.Second * 60),
			MaxLifeTime:       xtime.Duration(time.Hour * 2),
			ForceCloseWait:    xtime.Duration(time.Second * 20),
			KeepAliveInterval: xtime.Duration(time.Second * 60),
			KeepAliveTimeout:  xtime.Duration(time.Second * 20),
		},
		Backoff: &Backoff{MaxDelay: 300,
			BaseDelay: 3,
			Factor:    1.8,
			Jitter:    1.3,
		},
		Redis: &Redis{
			Network:      "tcp",
			Addr:         "127.0.0.1:6379",
			Active:       60000,
			Idle:         1024,
			DialTimeout:  xtime.Duration(200 * time.Second),
			ReadTimeout:  xtime.Duration(500 * time.Microsecond),
			WriteTimeout: xtime.Duration(500 * time.Microsecond),
			IdleTimeout:  xtime.Duration(120 * time.Second),
			Expire:       xtime.Duration(30 * time.Minute),
		},
	}
	cfg.Regions = make(map[string][]string, 0)
	cfg.Regions["gz"] = []string{"广东","福建","广西","海南","湖南","四川","贵州","云南","西藏","香港","澳门"}




	return cfg
}
