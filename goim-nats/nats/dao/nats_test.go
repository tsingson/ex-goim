package dao

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/tsingson/ex-goim/goim-nats/nats/logic/conf"
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
	d.PushMsg(context.TODO(), 122, "room111", []string{"test", "tttt"}, []byte("test"))
}

func BenchmarkNatsDao_PushMsg(b *testing.B) {
	// b.StopTimer()
	//
	// b.StartTimer()
	for n := 0; n < b.N; n++ {
		d.PushMsg(context.TODO(), 122, "room111", []string{strconv.Itoa(n), "tttt"}, []byte("test"))
	}
}
