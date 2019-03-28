package job

import (
	comet "github.com/tsingson/goim/api/comet/grpc"
)

type CometProcess interface {
	Push(arg *comet.PushMsgReq) (err error)
	BroadcastRoom(arg *comet.BroadcastRoomReq) (err error)
	Broadcast(arg *comet.BroadcastReq) (err error)
	Close() (err error)
}
