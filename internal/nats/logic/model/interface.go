package model

import (
	"context"
)

type Dao interface {
	PushMsg(c context.Context, op int32, server string, keys []string, msg []byte) (err error)
	BroadcastRoomMsg(c context.Context, op int32, room string, msg []byte) (err error)
	AddMapping(c context.Context, mid int64, key, server string) (err error)
	ExpireMapping(c context.Context, mid int64, key string) (has bool, err error)
	DelMapping(c context.Context, mid int64, key, server string) (has bool, err error)
	ServersByKeys(c context.Context, keys []string) (res []string, err error)
	KeysByMids(c context.Context, mids []int64) (ress map[string]string, olMids []int64, err error)
	AddServerOnline(c context.Context, server string, online *Online) (err error)
	ServerOnline(c context.Context, server string) (online *Online, err error)
	DelServerOnline(c context.Context, server string) (err error)
	Close() error
}
