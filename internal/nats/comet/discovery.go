package comet

import (
	"fmt"
	"time"

	"github.com/tsingson/discovery/naming"
)

func Register(cfg *naming.Config, appid string, addrs []string) {

	// dis := naming.New(cfg.Discovery)
	// resolver.Register(dis)

	// conf := &naming.Config{
	// 	Nodes: []string{"127.0.0.1:7171"}, // NOTE: 配置种子节点(1个或多个)，client内部可根据/discovery/nodes节点获取全部node(方便后面增减节点)
	// 	Zone:  "sh1",
	// 	Env:   "test",
	// }
	dis := naming.New(cfg)
	ins := &naming.Instance{
		Zone:  cfg.Zone,
		Env:   cfg.Env,
		AppID: "goim.comet",
		// Hostname:"", // NOTE: hostname 不需要，会优先使用discovery new时Config配置的值，如没有则从os.Hostname方法获取！！！
		Addrs:    []string{"http://172.0.0.1:8888"},
		LastTs:   time.Now().Unix(),
		Metadata: map[string]string{"weight": "10"},
	}
	cancel, _ := dis.Register(ins)
	defer cancel() // NOTE: 注意一般在进程退出的时候执行，会调用discovery的cancel接口，使实例从discovery移除
	fmt.Println("register")

}
