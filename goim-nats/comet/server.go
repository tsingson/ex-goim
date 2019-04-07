package comet

import (
	"context"
	"math/rand"
	"time"

	log "github.com/tsingson/zaplogger"
	"github.com/zhenjl/cityhash"

	logic "github.com/tsingson/ex-goim/api/logic/grpc"

	"github.com/tsingson/ex-goim/goim-nats/comet/conf"
)

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
)

// Server is comet server.
type Server struct {
	c         *conf.CometConfig
	round     *Round    // accept round store
	buckets   []*Bucket // subkey bucket
	bucketIdx uint32

	serverID  string
	rpcClient logic.LogicClient
}

// NewServer returns a new Server.
func NewServer(cfg *conf.CometConfig) *Server {
	s := &Server{
		c:         cfg,
		round:     NewRound(cfg),
		rpcClient: NewLogicClient(cfg.RPCClient),
	}
	// init bucket
	s.buckets = make([]*Bucket, cfg.Bucket.Size)
	s.bucketIdx = uint32(cfg.Bucket.Size)
	for i := 0; i < cfg.Bucket.Size; i++ {
		s.buckets[i] = NewBucket(cfg.Bucket)
	}
	s.serverID = cfg.Env.Host
	go s.onlineproc()
	return s
}

// Buckets return all buckets.
func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

// Bucket get the bucket by subkey.
func (s *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	if conf.Conf.Debug {
		log.Infof("%s hit channel bucket index: %d use cityhash", subKey, idx)
	}
	return s.buckets[idx]
}

// RandServerHearbeat rand server heartbeat.
func (s *Server) RandServerHearbeat() time.Duration {
	return (minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat))))
}

// Close close the server.
func (s *Server) Close() (err error) {
	return
}

func (s *Server) onlineproc() {
	for {
		var (
			allRoomsCount map[string]int32
			err           error
		)
		roomCount := make(map[string]int32)
		for _, bucket := range s.buckets {
			for roomID, count := range bucket.RoomsCount() {
				roomCount[roomID] += count
			}
		}
		if allRoomsCount, err = s.RenewOnline(context.Background(), s.serverID, roomCount); err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, bucket := range s.buckets {
			bucket.UpRoomsCount(allRoomsCount)
		}
		time.Sleep(time.Second * 10)
	}
}
