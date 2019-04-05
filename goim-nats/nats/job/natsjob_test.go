package job

import (
	"os"
	"testing"

	"github.com/tsingson/ex-goim/goim-nats/nats/job/conf"
)

var (
	d *NatsJob
)

func TestMain(m *testing.M) {
	conf.Conf = conf.Default()
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestNatsJob_ConsumeCheck(t *testing.T) {
	d.ConsumeCheck()
}

// func TestNatsJob_Subscribe(t *testing.T) {
// 	d.Subscribe(d.c.Nats.Channel, d.c.Nats.ChannelID)
// }

// func TestNatsJob_WatchComet(t *testing.T) {
// 	d.WatchComet(d.c.Discovery)
// }
