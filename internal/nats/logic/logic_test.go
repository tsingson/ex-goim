package logic

import (
	"context"
	"os"
	"testing"

	"github.com/tsingson/goim/internal/nats/logic/conf"
)

var (
	lg *NatsLogic
)

func TestMain(m *testing.M) {

	lg = New(conf.Conf)
	if err := lg.Ping(context.TODO()); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
