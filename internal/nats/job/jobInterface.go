package job

import (
	"github.com/tsingson/discovery/naming"
)

type Job interface {
	WatchComet(c *naming.Config)
	Subscribe(channel, channelID string) error
	Consume()
	Close() error
}
