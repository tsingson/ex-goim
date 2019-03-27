package comet

import (
	"testing"

	"github.com/tsingson/discovery/naming"
)

func TestRegister(t *testing.T) {

	cfg := &naming.Config{
		Nodes: []string{"127.0.0.1:7171"}, // NOTE: 配置种子节点(1个或多个)，client内部可根据/discovery/nodes节点获取全部node(方便后面增减节点)
		Zone:  "sh1",
		Env:   "test",
	}

	AppID := "goim.comet"
	// Hostname:"", // NOTE: hostname 不需要，会优先使用discovery new时Config配置的值，如没有则从os.Hostname方法获取！！！
	Addrs := []string{"http://172.0.0.1:8888"}

	Register(cfg, AppID, Addrs)
}
