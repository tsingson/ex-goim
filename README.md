# (Experimental) goim via nats 
fork from [goim](https://github.com/Terry-Mao/goim)  and support nats  to replace kafka / zookeeper


## chinese note 中文说明
goim 是 非常成功的 IM ( 即时消息平台), 依赖项为 kafka ( 消息队列) + zookeeper ( 分布式扩展) + bilibili/discovery( 服务发现与均衡) , 由于 kafka / zk 在部署上与 golang 的单一可执行文件相比, 稍复杂, 这里 fork 了 goim 并修改为 nats ( 并抽象 dao 为接口, 提供其他队列支持的可能性) 

由于修改比较大, 暂时用新的 repo 来进行代码管理, 以后看情况是否能回归到 Terry-Mao 的主线版本上. 



同时, 与 [goim](https://github.com/Terry-Mao/goim) 有所差异的重要一点:

 **这个 fork 是实验性质练手项目, 请不要用于生产环境!!  this repo personal Experimental , DO NOT use in production!!**




### Movation 动机

作为一个曾经的架构师(2005~2014, Utstarcom IPTV/OTT 事业部) 与当前自由的技术类咨询与服务从业者, 有合作伙伴询问 IM 用在视频直播中的方案, 我作了一些学习与研究 

中国 B站( BiliBili ) 的技术领军 [毛剑](https://github.com/Terry-Mao/) 是我神交以久的技术专家,   [goim](https://github.com/Terry-Mao/goim)  是一个非常成功的架构示例, 其模块拆分, 接口设计, 技术造型 , 以及部署方式, 都是一个互联网商用项目典范.  

同时,  另一位技术专家 [Xin.zh](https://github.com/alexstocks) 的文章 [一套高可用实时消息系统实现](https://alexstocks.github.io/html/pubsub.html)  给我很大启发.  

在电信/广电的几年经历, 这一次,  闲来无事, 算是满怀着在巨人肩头的感谢与敬意, 尝试写一些代码来加深学习.

个人在 Utstarcom 以业务平台架构师/解决方案工程师/ IPTV播控产品线 release manager 角色有一些时间,  除了技术方案的原型代码撰写与演示以外, 甚少参与实际撰写代码的工作 ,  这次写写代码也是有趣的过程  :P

欢迎指点/交流....




### 主要变更

![arch](./docs/arch.png)


  - [x] 消息队列修改为 [nats](https://github.com/nats-io/gnatsd) + [liftbridge](https://github.com/liftbridge-io/liftbridge)  注:  [liftbridge](https://github.com/liftbridge-io/liftbridge) 替代了 [nats-streaming-server](https://github.com/nats-io/nats-streaming-server) , 相关信息参见[liftbridge介绍文章](https://bravenewgeek.com/introducing-liftbridge-lightweight-fault-tolerant-message-streams/)
  - [x] 日志替换为 [uber-go/zap](https://github.com/uber-go/zap), 替换原一是因为 zap 快一点, 二是个人更为熟悉这个日志库 
  - [x] 修改了三个应用程序的启动方式, 去除了所有启动参数, 改为读取指定的 toml 配置文件( 同时, 预留接口以久将来进行读取远程配置, 及配置参数动态加载) 
  - [ ] 深入 comet / logic 模块尽量抽象接口, 以及作一些外部对接, 以及 技术实现的替换 ,  例如, websocket 更换为 [ws](https://github.com/gobwas/ws))
  - [ ] comet 增加 gRPC 与 rpc 接口,  tcp /websocket 等扩展增加用户注册 /发送消息/ 变更聊天室 / 查看历史消息等
  - [ ] 增加 gRPC 拦截 ( 支持 chatbot 等), 增加支持 消息历史存储/ relay 等接口



### 文件结构
修改文件如下示意, 支持 nats 的应用程序在以下路径, 每个应用下的 toml 为对应的配置
/cmd/nats/discoveryd-config.toml 为  discovery 的配置

```
├── cmd
│   └── nats
│       ├── comet
│       │   ├── comet-config.toml
│       │   └── main.go
│       ├── discoveryd-config.toml
│       ├── job
│       │   ├── job-config.toml
│       │   └── main.go
│       └── logic
│           ├── logic-config.toml
│           └── main.go
```

支持 nats 的库文件在 
```
/internal/nats/
```
路径下, 除配置文件以外, 所有库的调用方式与原 goim 相同



###  goim guide 安装/编译/使用指南(WIP)
参见 [/goim-usage-cn.md](goim-usage-cn.md) ( chinese )



goim v2.0
==============

[![Build Status](https://travis-ci.org/Terry-Mao/goim.svg?branch=master)](https://travis-ci.org/Terry-Mao/goim) 
[![Go Report Card](https://goreportcard.com/badge/github.com/Terry-Mao/goim)](https://goreportcard.com/report/github.com/Terry-Mao/goim)
[![codecov](https://codecov.io/gh/Terry-Mao/goim/branch/master/graph/badge.svg)](https://codecov.io/gh/Terry-Mao/goim)

goim is a im server writen by golang.

## Features
 * Light weight
 * High performance
 * Pure Golang
 * Supports single push, multiple push and broadcasting
 * Supports one key to multiple subscribers (Configurable maximum subscribers count)
 * Supports heartbeats (Application heartbeats, TCP, KeepAlive, HTTP long pulling)
 * Supports authentication (Unauthenticated user can't subscribe)
 * Supports multiple protocols (WebSocket，TCP，HTTP）
 * Scalable architecture (Unlimited dynamic job and logic modules)
 * Asynchronous push notification based on Kafka

## Architecture
![arch](./docs/arch.png)

## Quick Start

### Build
```
    make build
```

### Run
```
    make run
    make stop

    // or
    nohup target/logic -conf=target/logic.toml -region=sh -zone=sh001 deploy.env=dev weight=10 2>&1 > target/logic.log &
    nohup target/comet -conf=target/comet.toml -region=sh -zone=sh001 deploy.env=dev weight=10 addrs=127.0.0.1 2>&1 > target/logic.log &
    nohup target/job -conf=target/job.toml -region=sh -zone=sh001 deploy.env=dev 2>&1 > target/logic.log &

```
### Environment
```
    env:
    export REGION=sh
    export ZONE=sh001
    export DEPLOY_ENV=dev

    supervisor:
    environment=REGION=sh,ZONE=sh001,DEPLOY_ENV=dev

    go flag:
    -region=sh -zone=sh001 deploy.env=dev
```
### Configuration
You can view the comments in target/comet.toml,logic.toml,job.toml to understand the meaning of the config.

### Dependencies
[Discovery](https://github.com/Bilibili/discovery)

[Kafka](https://kafka.apache.org/quickstart)

## Document
[Protocol](./docs/protocol.png)

[English](./README_en.md)

[中文](./README_cn.md)

## Examples
Websocket: [Websocket Client Demo](https://github.com/Terry-Mao/goim/tree/master/examples/javascript)

Android: [Android](https://github.com/roamdy/goim-sdk)

iOS: [iOS](https://github.com/roamdy/goim-oc-sdk)

## Benchmark
![benchmark](./docs/benchmark.jpg)

### Benchmark Server
| CPU | Memory | OS | Instance |
| :---- | :---- | :---- | :---- |
| Intel(R) Xeon(R) CPU E5-2630 v2 @ 2.60GHz  | DDR3 32GB | Debian GNU/Linux 8 | 1 |

### Benchmark Case
* Online: 1,000,000
* Duration: 15min
* Push Speed: 40/s (broadcast room)
* Push Message: {"test":1}
* Received calc mode: 1s per times, total 30 times

### Benchmark Resource
* CPU: 2000%~2300%
* Memory: 14GB
* GC Pause: 504ms
* Network: Incoming(450MBit/s), Outgoing(4.39GBit/s)

### Benchmark Result
* Received: 35,900,000/s

[中文](./docs/benchmark_cn.md)

[English](./docs/benchmark_en.md)

## LICENSE
goim is is distributed under the terms of the MIT License.
