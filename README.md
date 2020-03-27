# DistributedSystem 分布式系统

架构
```text
分布式系统
01snowflake                     分布式 id 生成器   --snowflake/sonyflake
02lock                          分布式 锁          --redis etcd zk redlock  sync.Mutex(单机锁)
03timer                         延时任务系统       ----四叉树，消息队列MQ设置集群
04es                            分布式搜索引擎     --es
05loadbalance                   负载均衡          --shuffle算法  fisher yates算法 、 rand.Seed(time.Now().UnixNano())
06configurationmanagement       分布式配置管理     --etcd  
07distributedcrawler            分布式爬虫        --nats(消息队列),colly（爬虫）
08pubsub                        发布订阅模型       --包pubsub实现了一个简单的多主题pub-sub库。
09gatefs                        控制并发数         --vfs/gatefs包，为了控制访问该虚拟文件系统的最大并发数。
```