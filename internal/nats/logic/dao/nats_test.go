package dao

import (
	"context"
	"os"
	"testing"

	"github.com/tsingson/goim/internal/nats/logic/conf"
)

var (
	d *NatsDao
)

func TestMain(m *testing.M) {
	conf.Conf = conf.Default()
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestPushMsg(t *testing.T) {
      d.PushMsg(context.TODO(), 122, "room111", []string {"test", "tttt"}, []byte("test"))
}
