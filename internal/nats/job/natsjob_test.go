package job

import (
	"os"
	"testing"

	"github.com/tsingson/goim/internal/nats/job/conf"
)

var (
	d *NatsJob
)

func TestMain(m *testing.M) {
	conf.Conf = conf.Default()
	d = New(conf.Conf)
	os.Exit(m.Run())
}

// func TestNatsJob_Consume(t *testing.T) {
//  d.Consume()
// }

// func TestNatsJob_Subscribe(t *testing.T) {
// 	d.Subscribe(d.c.Nats.Channel, d.c.Nats.ChannelID)
// }

// func TestNatsJob_WatchComet(t *testing.T) {
// 	d.WatchComet(d.c.Discovery)
// }
