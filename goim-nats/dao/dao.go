package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/liftbridge-io/go-liftbridge"
	"github.com/nats-io/go-nats"

	"github.com/tsingson/ex-goim/goim-nats/logic/conf"
)

// Dao dao for nats
type Dao struct {
	c           *conf.LogicConfig
	natsClient  *nats.Conn
	liftClient  liftbridge.Client
	redis       *redis.Pool
	redisExpire int32
}

// NatsDao alias name of dao
type NatsDao = Dao

// LogicConfig configuration for nats / liftbridge queue
type Config struct {
	Channel   string
	ChannelID string
	Group     string
	NatsAddr  string
	LiftAddr  string
}

// NatsCOnfig alias name of Config
type NatsConfig = Config
