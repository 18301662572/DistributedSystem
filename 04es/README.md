# 分布式搜索引擎 -es

**elasticsearch 是开源分布式搜索引擎的霸主，其依赖于 Lucene 实现**，在部署和运维方面做了很多优化。<br/>
简单配置客户端 ip 和端口就可以了<br/>

### 1.倒排列表
```text
虽然 es 是针对搜索场景来定制的，但实际应用中常常用 es 来作为 database 来使用，就是因为倒排列表的特性。
对 es 中的数据进行查询时，本质就是求多个排好序的序列求交集。
┌─────────────────┐       ┌─────────────┬─────────────┬─────────────┬─────────────┐
│  order_id: 103  │──────▶│ doc_id:4231 │ doc_id:4333 │ doc_id:5123 │ doc_id:9999 │
└─────────────────┘       └─────────────┴─────────────┴─────────────┴─────────────┘
┌─────────────────┐       ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐
│  sku_id: 30221  │──────▶│ doc_id:4231 │ doc_id:5123 │ doc_id:5644 │ doc_id:7801 │ doc_id:9999 │
└─────────────────┘       └─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘
┌─────────────────┐       ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐
│   city_id: 3    │──────▶│ doc_id:5123 │ doc_id:9999 │doc_id:10232 │doc_id:54321 │doc_id:63142 │doc_id:71230 │doc_id:90123 │
└─────────────────┘       └─────────────┴─────────────┴─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘

eg：当用户搜索 ‘天气很好’ 时，其实就是求：天气、气很、很好三组倒排列表的交集，但这里的相等判断逻辑有些特殊，用伪代码表示一下：
func equal() {
    if postEntry.docID of '天气' == postEntry.docID of '气很' && postEntry.offset + 1 of '天气' == postEntry.offset of '气很' {
        return true
    }
    if postEntry.docID of '气很' == postEntry.docID of '很好' && postEntry.offset + 1 of '气很' == postEntry.offset of '很好' {
        return true
    }
    if postEntry.docID of '天气' == postEntry.docID of '很好' && postEntry.offset + 2 of '天气' == postEntry.offset of '很好' {
        return true
    }
    return false
}
多个有序列表求交集的时间复杂度是：O(N * M)， N 为给定列表当中元素数最小的集合， M 为给定列表的个数。

在整个算法中起决定作用的一是最短的倒排列表的长度，其次是词数总和，一般词数不会很大
(想像一下，你会在搜索引擎里输入几百字来搜索么？)，所以起决定性作用的，一般是所有倒排列表中，最短的那一个的长度。

文档总数很多的情况下，搜索词的倒排列表最短的那一个不长时，搜索速度也是很快的。如果用关系型数据库，
那就需要按照索引(如果有的话)来慢慢扫描了。
```

### 2.查询 DSL
**es 定义了一套查询 DSL，当我们把 es 当数据库使用时，需要用到其 bool 查询。**<br/>

用 bool should must 可以表示 and 的逻辑：<br/>
```json
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "field_1": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_2": {
              "query": "2",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_3": {
              "query": "3",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_4": {
              "query": "4",
              "type": "phrase"
            }
          }
        }
      ]
    }
  },
  "from": 0,
  "size": 1
}
```
```text
看起来比较麻烦，但表达的意思很简单：
if field_1 == 1 && field_2 == 2 && field_3 == 3 && field_4 == 4 {
    return true
}
```
用 bool should query 可以表示 or 的逻辑：<br/>
```json
{
  "query": {
    "bool": {
      "should": [
        {
          "match": {
            "field_1": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_2": {
              "query": "3",
              "type": "phrase"
            }
          }
        }
      ]
    }
  },
  "from": 0,
  "size": 1
}
```
```text
这里表示的是类似：
if field_1 == 1 || field_2 == 2 {
    return true
}
```
```text
es 的 Bool Query 方案，实际上就是用 json 来表达了这种程序语言中的 Boolean Expression，为什么可以这么做呢？
因为 json 本身是可以表达树形结构的，我们的程序代码在被编译器 parse 之后，也会变成 AST，而 AST 抽象语法树，顾名思义，
就是树形结构。理论上 json 能够完备地表达一段程序代码被 parse 之后的结果。
这里的 Boolean Expression 被编译器 Parse 之后也会生成差不多的树形结构，而且只是整个编译器实现的一个很小的子集。
```

### 3.基于 client sdk 做开发 --elastic.go
```text
将 sql 转换为 DSL:
具体实现:github.com/cch123/elasticsql
eg:
有一段 bool 表达式，user_id = 1 and (product_id = 1 and (star_num = 4 or star_num = 5) and banned =1
sql:
select * from xxx where user_id = 1 and (product_id = 1 and (star_num = 4 or star_num = 5) and banned = 1)

es 的 DSL :
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "user_id": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "product_id": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "bool": {
            "should": [
              {
                "match": {
                  "star_num": {
                    "query": "4",
                    "type": "phrase"
                  }
                }
              },
              {
                "match": {
                  "star_num": {
                    "query": "5",
                    "type": "phrase"
                  }
                }
              }
            ]
          }
        },
        {
          "match": {
            "banned": {
              "query": "1",
              "type": "phrase"
            }
          }
        }
      ]
    }
  },
  "from": 0,
  "size": 1
}
```

### 4.异构数据同步
```text
我们很少直接向搜索引擎中写入数据。
更为常见的方式是，将 MySQL 或其它关系型数据中的数据同步到搜索引擎中。而搜索引擎的使用方只能对数据进行查询，无法进行修改和删除。

常见的同步方案有两种：
    1.通过时间戳进行增量数据同步
    逻辑实际上就是一条 SQL：
    select * from wms_orders where update_time >= date_sub(now(), interval 11 minute);
    取最近 11 分钟有变动的数据覆盖更新到 es 中。这种方案的缺点显而易见，我们必须要求业务数据严格遵守一定的规范。
    比如这里的，必须要有 update_time 字段，并且每次创建和更新都要保证该字段有正确的时间值。否则我们的同步逻辑就会丢失数据。
    2.通过 binlog 进行数据同步
                ┌────────────────────────┐
                │      MySQL master      │
                └────────────────────────┘
                             │
                             │
                             │
                             │
                             │
                             │
                             ▼
                   ┌───────────────────┐
                   │ row format binlog │
                   └───────────────────┘
                             │
                             │
                             │
             ┌───────────────┴──────────────┐
             │                              │
             │                              │
             ▼                              ▼
┌────────────────────────┐         ┌─────────────────┐
│      MySQL slave       │         │      canal      │
└────────────────────────┘         └─────────────────┘
                                            │
                                  ┌─────────┴──────────┐
                                  │   parsed binlog    │
                                  └─────────┬──────────┘
                                            │
                                            ▼
                                   ┌────────────────┐
                                   │     kafka      │─────┐
                                   └────────────────┘     │
                                                          │
                                                          │
                                                          │
                                                          │
                                              ┌───────────┴──────┐
                                              │  kafka consumer  │
                                              └───────────┬──────┘
                                                          │
                                                          │
                                                          │
                                                          │      ┌────────────────┐
                                                          └─────▶│ elasticsearch  │
                                                                 └────────────────┘
        业界使用较多的是阿里开源的 canal，来进行 binlog 解析与同步。canal 会伪装成 MySQL 的从库，
        然后解析好行格式的 binlog，再以更容易解析的格式(例如 json) 发送到消息队列。

        由下游的 kafka 消费者负责把上游数据表的自增主键作为 es 的 document 的 id 进行写入，这样可以保证每次接收到 binlog 时，
        对应 id 的数据都被覆盖更新为最新。MySQL 的 row 格式的 binlog 会将每条记录的所有字段都提供给下游，所以实际上
        在向异构数据目标同步数据时，不需要考虑数据是插入还是更新，只要一律按 id 进行覆盖即可。

        这种模式同样需要业务遵守一条数据表规范，即表中必须有唯一主键 id 来保证我们进入 es 的数据不会发生重复。
        一旦不遵守该规范，那么就会在同步时导致数据重复。当然，你也可以为每一张需要的表去定制消费者的逻辑，
        这就不是通用系统讨论的范畴了。
```