package model

import (
	"context"

	"github.com/tsingson/discovery/naming"

	"github.com/tsingson/goim/api/comet/grpc"
	pb "github.com/tsingson/goim/api/logic/grpc"
)

// Action  interface for logic
type LogicProcess interface {
	Connect(c context.Context, server, cookie string, token []byte) (mid int64, key, roomID string, accepts []int32, hb int64, err error)
	Disconnect(c context.Context, mid int64, key, server string) (has bool, err error)
	Heartbeat(c context.Context, mid int64, key, server string) (err error)
	RenewOnline(c context.Context, server string, roomCount map[string]int32) (map[string]int32, error)
	Receive(c context.Context, mid int64, proto *grpc.Proto) (err error)
	PushKeys(c context.Context, op int32, keys []string, msg []byte) (err error)
	PushMids(c context.Context, op int32, mids []int64, msg []byte) (err error)
	PushRoom(c context.Context, op int32, typ, room string, msg []byte) (err error)
	PushAll(c context.Context, op, speed int32, msg []byte) (err error)
	NodesInstances(c context.Context) (res []*naming.Instance)
	NodesWeighted(c context.Context, platform, clientIP string) *pb.NodesReply
	Ping(c context.Context) (err error)
	Close()
	OnlineTop(c context.Context, typ string, n int) (tops []*Top, err error)
	OnlineRoom(c context.Context, typ string, rooms []string) (res map[string]int32, err error)
	OnlineTotal(c context.Context) (int64, int64)
}
