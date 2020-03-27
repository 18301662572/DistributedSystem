package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

//nats 是 Go 实现的一个高性能分布式消息队列，适用于高并发高吞吐量的消息分发场景。

//基本消息生产
func publish(){
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		// log error
		return
	}
	// 指定 subject 为 tasks，消息内容随意
	err = nc.Publish("tasks", []byte("your task content"))
	nc.Flush()
}

//基本消息消费 queue subscribe
//直接使用 nats 的 subscribe api 并不能达到任务分发的目的，因为 pub sub 本身是广播性质的。所有消费者都会收到完全一样的所有消息。
// 除了普通的 subscribe 之外，nats 还提供了 queue subscribe 的功能。
// 只要提供一个 queue group 名字(类似 kafka 中的 consumer group)，即可均衡地将任务分发给消费者。
func queueSubscribe(){
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		// log error
		return
	}
	// queue subscribe 相当于在消费者之间进行任务分发的分支均衡
	// 前提是所有消费者都使用 workers 这个 queue
	// nats 中的 queue 概念上类似于 kafka 中的 consumer group
	sub, err := nc.QueueSubscribeSync("tasks", "workers")
	if err != nil {
		// log error
		return
	}
	var msg *nats.Msg
	for {
		msg, err = sub.NextMsg(time.Hour * 10000)
		if err != nil {
			// log error
			break
		}
		// 正确地消费到了消息
		// 可用 nats.Msg 对象处理任务
	}
	fmt.Println(msg)
}

