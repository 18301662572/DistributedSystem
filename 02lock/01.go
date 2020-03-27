package main

import (
	"sync"
)

//进程内加锁（单机锁）

//❯❯❯ go run local_lock.go
//1000

func main1(){
	var counter int
	// ... 省略之前部分
	var wg sync.WaitGroup
	var l sync.Mutex
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Lock()
			counter++
			l.Unlock()
		}()
	}
	wg.Wait()
	println(counter)
	// ... 省略之后部分
}

