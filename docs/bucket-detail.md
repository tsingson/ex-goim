#  bucket detail 
--------------------- 
作者：nini_boom 
来源：CSDN 
原文：https://blog.csdn.net/nini_boom/article/details/82940237 

```
package main
import(
	"goim/libs/define"
	"goim/libs/proto" //libs包含了对bufio,bytes,net等原生包的改写，
		 	  		  //个人觉得还是很考验阅读能力的
	"sync"
	"sync/atomic"
)

type BucketOptions struct { //Bucket操作
	ChannelSize 		int		
	RoomSize 			int		
	RoutineAmount 		uint64	//Size和Amount这些貌似是固定大小的
}

type Bucket struct { 
	cLock 		sync.RWMutex		//确保chs的协程安全
	chs 		map[string]*Channel	//订阅字符为key，channel实体为value
	boptions 	BucketOptions		
	rooms 		map[int32]*Room
	routines 	[]chan *proto.BoardcaseRoomArg //用于广播给bucket下所有的room
	routinesNum uint64
}

//创建Bucket实体
func NewBucket(boptions BucketOptions) (b *Bucket){
  b = new(Bucket)
  b.chs = make(map[string]*Channel, boptions.ChannelSize)
  b.boptions = boptions
  
  //room
  b.rooms = make(map[int32]*Room, boptions.RoomSize)
  //为广播房间数创建对应数量的gorounite?
  b.routines = make([]chan *proto.BoardcastRoomArg, boptions.RoutineAmount)
  for i:= uint64(0); i < boptions.RoutineAmount; i++ {
    c := make(chan *proto.BoardcastRoomArg, boptions.RoutineSize)
    b.routines[i] = c
    go b.roomproc(c)
  }
  return 
}

//bucket中Channel个数
func (b *Bucket) ChannelCount() int {
  return len(b.chs)
}

//bucket中Room个数
func (b *Bucket) RoomCount() int {
  return len(b.rooms)
}

//sub key为key，ch为channel，rid是Room号，一并初始化并put到Bucket中
func (b *Bucket) Put(key string, rid int32, ch *Channel) (err error){
  var (
  	room *Room
    ok bool
  )
  //map非线程安全，故加锁保护
  b.cLock.Lock()
  b.chs[key] = ch
  if rid != define.NoRoom { //define.NoRoom = -1，意思是没有房间信息
    if room, ok = b.rooms[rid]; !ok { 
      room = NewRoom(rid)
      b.rooms[rid] = room
    }
    ch.Room = room
  }
  b.cLock.Unlock()
  if room != nil {
    err = room.Put(ch)
  }
  return
}

//删除channel和room
//思考，channel和room有什么样的关系
//猜想，一个channel只对应唯一room,一个room有多个channel
func (b *Bucket) Del(key string) {
  var (
  	ok 		bool
    ch 		*Channel
    room 	*Room
  )
  //
  b.cLock.Lock()
  if ch, ok = b.chs[key]; ok {
    room = ch.Room
    delete(b.chs, key)
  }
  b.cLock.Unlock()
  if room != nil && room.Del(ch){
    //空room必须从bucket中删除
    b.DelRoom(room)
  }
}

//Channel get a channel by sub key
func (b *Bucket) Channel(key string) (ch *Channel) {
  b.cLock.RLock()//Lock用于读写不确定的情况下，有强制性
  				//RLock用于读多写少的情况，这就是使用RLock的原因
  				//也可理解互斥锁和读写锁的区别
  				//课后作业，到底底层区别在哪
  ch = b.chs[key]
  b.cLock.RUnlock()
  return
}

//广播消息给Bucket下所有的channels
//由此可见，Proto是消息单位体
func (b *Bucket) Broadcast(p *proto.Proto) {
  var ch *Channel
  b.cLock.RLock()
  for _, ch = range b.chs {
    ch.Puch(p)
  }
  b.cLock.RUnlock()
}

//get a room by rid
func (b *Bucket) Room(rid int32) (room *Room) {
  b.cLock.RLock()
  room, _ = b.rooms[rid]
  b.cLock.RUnlock()
  return
}

//delete a room of bucket by room pointer
func (b *Bucket) DelRoom(room *Room) {
  b.cLock.Lock()
  delete(b.rooms, room.Id)
  b.cLock.Unlock()
  room.Close()
  return
}

//向Bucket下所有Room广播信息
func (b *Bucket) BroadcastRoom(arg *proto.BoardcastRoomArg) {
  //原子增加，最终的数量不超过RountineAmount
  //课后作业，原子操作和加锁保护的区别在哪？
  num := atomic.AddUint64(&b.routinesNum, 1) % b.boptions.RoutineAmount
  b.routines[num] <- arg
}

//获取在线数据大于1(Online > 1)的所有room
func (b *Bucket) Rooms() (res map[int32]struct{}) {
  var (
  	roomId 	int32
    room 	*Room
  )
  res = make(map[int32]struct{})
  b.cLock.RLock()
  for roomId, room = range b.rooms {
    if room.Online > 0 {
      res[roomId] = struct{}{}
    }
  }
  b.cLock.RUnlock()
  return
}

func (b *Bucket) roomproc(c chan *proto.BoardcastRoomArg) {
  for {
    var (
    	arg 	*proto.BoardcastRoomArg
      	room 	*Room
    )
    arg = <- c
    if room = b.Room(arg.RoomId); room != nil {
      room.Push(&arg.P)
    }
  }
}

 


```
