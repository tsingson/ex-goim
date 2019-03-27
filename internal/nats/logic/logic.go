package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/tsingson/discovery/naming"
	log "github.com/tsingson/zaplogger"

	"github.com/tsingson/goim/internal/nats/logic/dao"

	"github.com/tsingson/goim/internal/nats/logic/conf"

	"github.com/tsingson/goim/internal/nats/logic/model"
)

const (
	_onlineTick     = time.Second * 10
	_onlineDeadline = time.Minute * 5
)

// NatsLogic struct
type NatsLogic struct {
	c   *conf.NatsConfig
	dis *naming.Discovery
	dao *dao.NatsDao
	// online
	totalIPs   int64
	totalConns int64
	roomCount  map[string]int32
	// load balancer
	nodes        []*naming.Instance
	loadBalancer *LoadBalancer
	regions      map[string]string // province -> region
}

// New init
func New(c *conf.NatsConfig) (l *NatsLogic) {
	l = &NatsLogic{
		c:            c,
		dao:          dao.New(c),
		dis:          naming.New(c.Discovery),
		loadBalancer: NewLoadBalancer(),
		regions:      make(map[string]string),
	}
	l.initRegions()
	l.initNodes()
	_ = l.loadOnline()
	go l.onlineproc()
	return l
}

// Ping ping resources is ok.
func (l *NatsLogic) Ping(c context.Context) (err error) {
	return l.dao.Ping(c)
}

// Close close resources.
func (l *NatsLogic) Close() {
	l.dao.Close()
}

func (l *NatsLogic) initRegions() {
	for region, ps := range l.c.Regions {
		for _, province := range ps {
			l.regions[province] = region
		}
	}
}

func (l *NatsLogic) initNodes() {
	res := l.dis.Build("goim.comet")
	event := res.Watch()
	select {
	case _, ok := <-event:
		if ok {
			l.newNodes(res)
		} else {
			panic("discovery watch failed")
		}
	case <-time.After(10 * time.Second):
		log.Error("discovery start timeout")
	}
	go func() {
		for {
			if _, ok := <-event; !ok {
				return
			}
			l.newNodes(res)
		}
	}()
}

func (l *NatsLogic) newNodes(res naming.Resolver) {
	if zoneIns, ok := res.Fetch(); ok {
		var (
			totalConns int64
			totalIPs   int64
			allIns     []*naming.Instance
		)
		for _, zins := range zoneIns.Instances {
			for _, ins := range zins {
				if ins.Metadata == nil {
					log.Errorf("node instance metadata is empty(%+v)", ins)
					continue
				}
				offline, err := strconv.ParseBool(ins.Metadata[model.MetaOffline])
				if err != nil || offline {
					log.Warningf("strconv.ParseBool(offline:%t) error(%v)", offline, err)
					continue
				}
				conns, err := strconv.ParseInt(ins.Metadata[model.MetaConnCount], 10, 32)
				if err != nil {
					log.Errorf("strconv.ParseInt(conns:%d) error(%v)", conns, err)
					continue
				}
				ips, err := strconv.ParseInt(ins.Metadata[model.MetaIPCount], 10, 32)
				if err != nil {
					log.Errorf("strconv.ParseInt(ips:%d) error(%v)", ips, err)
					continue
				}
				totalConns += conns
				totalIPs += ips
				allIns = append(allIns, ins)
			}
		}
		l.totalConns = totalConns
		l.totalIPs = totalIPs
		l.nodes = allIns
		l.loadBalancer.Update(allIns)
	}
}

func (l *NatsLogic) onlineproc() {
	for {
		time.Sleep(_onlineTick)
		if err := l.loadOnline(); err != nil {
			log.Errorf("onlineproc error(%v)", err)
		}
	}
}

func (l *NatsLogic) loadOnline() (err error) {
	var (
		roomCount = make(map[string]int32)
	)
	for _, server := range l.nodes {
		var online *model.Online
		online, err = l.dao.ServerOnline(context.Background(), server.Hostname)
		if err != nil {
			return
		}
		if time.Since(time.Unix(online.Updated, 0)) > _onlineDeadline {
			_ = l.dao.DelServerOnline(context.Background(), server.Hostname)
			continue
		}
		for roomID, count := range online.RoomCount {
			roomCount[roomID] += count
		}
	}
	l.roomCount = roomCount
	return
}
