package comet

import (
	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/goim/internal/nats/comet/conf"
)

var whitelist *Whitelist

// Whitelist .
type Whitelist struct {
	log  *log.Logger
	list map[int64]struct{} // whitelist for debug
}

// InitWhitelist a whitelist struct.
func InitWhitelist(c *conf.Whitelist) (err error) {
	var mid int64

	whitelist = new(Whitelist)
	whitelist.log = log.New("", true)
	whitelist.list = make(map[int64]struct{})
	for _, mid = range c.Whitelist {
		whitelist.list[mid] = struct{}{}
	}

	return
}

// Contains whitelist contains a mid or not.
func (w *Whitelist) Contains(mid int64) (ok bool) {
	if mid > 0 {
		_, ok = w.list[mid]
	}
	return
}

// Printf calls l.Output to print to the logger.
func (w *Whitelist) Printf(format string, v ...interface{}) {
	w.log.Printf(format, v...)
}
