# 分布式爬虫 --一套任务分发和执行系统。(一个基于消息队列的爬虫)
# nats(消息生产+消费)，colly(网站爬虫)
作为收集数据的前置工作，有能力去写一个简单的或者复杂的爬虫，对于我们来说依然非常重要。

nats：https://github.com/nats-io/nats.go

架构
```text
colly.go         单机colly爬虫 
nats.go          分布式 nats基本消息生产消费
collypub         分布式 nats结合colly爬虫的消息生产
collysub         分布式 nats结合colly爬虫的消息消费
```

### 1.基于 colly 的单机爬虫 --colly.go
```text
直接命令行爬到 v2ex 在 Go tag 下的新贴，只要简单写一个爬虫即可。
images/ch6-dist-crawler.png
```

### 2.分布式爬虫 --一套任务分发和执行系统。
```text
想像一下，你们的信息分析系统运行非常之快。
获取信息的速度成为了瓶颈，虽然可以用上 Go 语言所有优秀的并发特性，将单机的 CPU 和网络带宽都用满，
但还是希望能够加快爬虫的爬取速度。在很多场景下，速度是有意义的：
    1.对于价格战期间的电商们来说，希望能够在对手价格变动后第一时间获取到其最新价格，
    再靠机器自动调整本家的商品价格。
    2.对于类似头条之类的 feed 流业务，信息的时效性也非常重要。如果我们慢吞吞地爬到的新闻是昨天的新闻，
    那对于用户来说就没有任何意义。
所以我们需要分布式爬虫。
从本质上来讲，分布式爬虫是一套任务分发和执行系统。
而常见的任务分发，因为上下游存在速度不匹配问题，必然要借助消息队列。
```

**本节我们来简单实现一个基于消息队列的爬虫，本节我们使用 nats 来做任务分发。**
### 3.nats 是 Go 实现的一个高性能分布式消息队列，适用于高并发高吞吐量的消息分发场景
https://github.com/nats-io/nats.go
```text
nats 的服务端项目是 gnatsd，客户端与 gnatsd 的通信方式为基于 tcp 的文本协议，非常简单：
    向 subject 为 task 发消息：
        images/ch6-09-nats-protocol-pub.png
    以 workers 的 queue 从 tasks subject 订阅消息：
        images/ch6-09-nats-protocol-sub.png 
    其中的 queue 参数是可选的，如果希望在分布式的消费端进行任务的负载均衡，而不是所有人都收到同样的消息，
    那么就要给消费端指定相同的 queue 名字。

1.nats基本消息生产           --nats-demo.go
    生产消息只要指定 subject 即可
2.nats基本消息消费           --nats-demo.go
    直接使用 nats 的 subscribe api 并不能达到任务分发的目的，因为 pub sub 本身是广播性质的。所有消费者都会收到完全一样的所有消息。    
    除了普通的 subscribe 之外，nats 还提供了 queue subscribe 的功能。
    只要提供一个 queue group 名字(类似 kafka 中的 consumer group)，即可均衡地将任务分发给消费者。
3.nats结合 colly 的消息生产  --colly-pub.go
    我们为每一个网站定制一个对应的 collector，并设置相应的规则，比如 v2ex，v2fx(虚构的)，
    再用简单的工厂方法来将该 collector 和其 host 对应起来 
4.nats结合 colly 的消息消费  --colly-sub.go
从代码层面上来讲，这里的生产者和消费者其实本质上差不多。
如果日后我们要灵活地支持增加、减少各种网站的爬取的话，应该思考如何将这些爬虫的策略、参数尽量地配置化。

```
