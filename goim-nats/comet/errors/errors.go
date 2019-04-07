package errors

import (
	"golang.org/x/xerrors"
)

// .
var (
	// server
	ErrHandshake = xerrors.New("handshake failed")
	ErrOperation = xerrors.New("request operation not valid")
	// ring
	ErrRingEmpty = xerrors.New("ring buffer empty")
	ErrRingFull  = xerrors.New("ring buffer full")
	// timer
	ErrTimerFull   = xerrors.New("timer full")
	ErrTimerEmpty  = xerrors.New("timer empty")
	ErrTimerNoItem = xerrors.New("timer item not exist")
	// channel
	ErrPushMsgArg   = xerrors.New("rpc pushmsg arg error")
	ErrPushMsgsArg  = xerrors.New("rpc pushmsgs arg error")
	ErrMPushMsgArg  = xerrors.New("rpc mpushmsg arg error")
	ErrMPushMsgsArg = xerrors.New("rpc mpushmsgs arg error")
	// bucket
	ErrBroadCastArg     = xerrors.New("rpc broadcast arg error")
	ErrBroadCastRoomArg = xerrors.New("rpc broadcast  room arg error")

	// room
	ErrRoomDroped = xerrors.New("room droped")
	// rpc
	ErrLogic = xerrors.New("logic rpc is not available")
)
