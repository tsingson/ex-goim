# goim-nats change log

## 说明
goim 是 非常成功的 IM ( 即时消息平台), 依赖项为 kafka ( 消息队列) + zookeeper ( 分布式扩展) + bilibili/discovery( 服务发现与均衡) , 由于 kafka / zk 在部署上与 golang 的单一可执行文件相比, 稍复杂, 加上为简化运维, 所以, 这里 fork 了 goim 并修改为 nats

1. 消息队列修改为 [nats](https://github.com/nats-io/gnatsd) + [liftbridge](https://github.com/liftbridge-io/liftbridge)  注:  [liftbridge](https://github.com/liftbridge-io/liftbridge) 替代了 [nats-streaming-server](https://github.com/nats-io/nats-streaming-server) , 相关信息参见[liftbridge介绍文章](https://bravenewgeek.com/introducing-liftbridge-lightweight-fault-tolerant-message-streams/)
2. 日志替换为 [uber-go/zap](https://github.com/uber-go/zap), 替换原一是因为 zap 快一点, 二是个人更为熟悉这个日志库


### Version 2.0.0
> 1.router has been changed to redis  
> 2.Support node with redis online heartbeat maintenance  
> 3.Support for gRPC and Discovery services  
> 4.Support node connection number and weight scheduling  
> 5.Support node scheduling by region  
> 6.Support instruction subscription  
> 7.Support the current connection room switch  
> 8.Support multiple room types ({type}://{room_id})  
> 9.Support sending messages by device_id  
> 10.Support for room message aggregation  
> 11.Supports IPv6  
