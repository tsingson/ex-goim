#  (Experimental) goim via nats  guide 安装/编译/使用指南(WIP)

## 1. dependency and run prepare 依赖与测试(运行)环境准备

1.1.  原 goim 依赖的消息队列为 kafka + zookeeper ,  现修改为 [nats](https://github.com/nats-io/gnatsd) + [liftbridge](https://github.com/liftbridge-io/liftbridge)  

	注1:  [liftbridge](https://github.com/liftbridge-io/liftbridge) 替代了 [nats-streaming-server](https://github.com/nats-io/nats-streaming-server) , 相关信息参见[liftbridge介绍文章](https://bravenewgeek.com/introducing-liftbridge-lightweight-fault-tolerant-message-streams/)
	注2:  TODO:  本 repo 代码并没有覆盖 goim 支持 kafka + zookeeper 代码, 同时支持两种 消息队列的方式稍后再写

1. 依赖 [bilibili/discovery](https://github.com/bilibili/discovery) , 一个在 [Netflix Eureka](https://github.com/Netflix/eureka)基础上增强实现的 AP 类型服务注册/发现单元.  相关的安装配置, 请自行参考官方文档




### 1. nats 安装与运行
源码编译安装
注: 以下 mkdir 创建路径指令可以不需要
```
mkdir ~/go/src/goim
cd ~/go/src/goim

go get github.com/nats-io/gnatsd
```
运行 ( 默认配置,  端口 :4222)
```
gnatsd
```

### 2. liftbridge 安装与运行
源码编译安装
注: 以下 mkdir 创建路径指令可以不需要
```
cd ~/go/src/goim

go get github.com/liftbridge-io/liftbridge
```
运行 ( 默认单例配置,  端口 :9292)
```
liftbridge --raft-bootstrap-seed
```

注: 关于 liftbridge 集群多例运行, 请参考原文档

### 3. discovery 运行

请下载 [编译好的discoveryd用于测试](https://github.com/tsingson/discovery/releases/download/v0.1.1/discoveryd.tar.bz2)  并用以下命令解码到 linux ( 验证环境为 cent os 7 ) 到某一路径下 

	注: 压缩包的打包方式 tar -cvjpf discoveryd.tar.bz2 ./discoveryd


```
tar -xvjf    discoveryd.tar.bz2
cd discoveryd
chmod +x ./discoveryd
ps -ef | grep discoveryd
```

 运行, 注: 该程序为后台 daemon 运行, 用 ps -ef | grep discoveryd 查看是否成功运行

```
discoveryd
```



## 2. testing and running 测试与运行 

分别在不同的 terminal 窗口中, 依次运行以下命令

```
gnatsd


liftbridge --raft-bootstrap-seed

redis-server 

discoveryd


```

## 3. running GOIM  运行GOIM

### 编译版本运行

请下载 [编译好的 goim 三个应用](https://github.com/tsingson/ex-goim/releases/download/v2.0.1/goim.tar.bz2)  并用以下命令解码到 linux ( 验证环境为 cent os 7 ) 到某一路径下 

分别在三个不同的 terminal 运行程序 


### 编译与配置
克隆本项目 repo 
```
$ git clone https://github.com/tsingson/ex-goim goim
cd goim
```


复制配置示例文件

在 /cmd/nats 路径下, 将三个 toml 配置文件 copy 到 $GOPATH/bin



分别在三个不同的 terminal 运行程序 

```
$ cd $GOPATH/bin
./logic
./comet
./job 
```


### Build
```
    make build
```

### Run
```
    make run
    make stop

    // or
	nohup ./target/logic   2>&1 > target/logic.log &
	nohup ./target/comet   2>&1 > target/comet.log &
	nohup ./target/job   2>&1 > target/job.log &

```




( to be continue...   待续...) 





## 3. Change log  修订日志
* 2019/03/28  初稿
