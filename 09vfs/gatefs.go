package main

import (
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/gatefs"
)

//其中vfs.OS("/path")基于本地文件系统构造一个虚拟的文件系统，然后gatefs.New基于现有的虚拟文件系统构造一个并发受控的虚拟文件系统。
// 并发数控制的原理在前面一节已经讲过，就是通过带缓存管道的发送和接收规则来实现最大并发阻塞：

func main() {
	fs := gatefs.New(vfs.OS("/path"), make(chan bool, 8))
	// ...
}
