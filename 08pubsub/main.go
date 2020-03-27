package main

import (
	"code.oldbody.com/studygolang/distributedsystem/08pubsub/pubsub"
	"fmt"
	"strings"
	"time"
)

//有两个订阅者分别订阅了全部主题和含有”golang”的主题：

func main() {
	p := pubsub.NewPublisher(100*time.Millisecond, 10)
	defer p.Close()
	//订阅所有主题
	all := p.Subscribe()
	//订阅含有“golang”的主题
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})
	p.Publish("hello,  world!")
	p.Publish("hello, golang!")
	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()
	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()
	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}
