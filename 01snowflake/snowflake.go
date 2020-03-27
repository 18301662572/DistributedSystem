package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"os"
)

//snowflake 雪花算法
//github.com/bwmarrin/snowflake 是一个相当轻量化的 snowflake 的 Go 实现。
//Epoch 就是本节开头讲的起始时间，NodeBits 指的是机器编号的位长，StepBits 指的是自增序列的位长。

func main() {
	n, err := snowflake.NewNode(1)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	for i := 0; i < 3; i++ {
		id := n.Generate()
		fmt.Println("id", id)
		fmt.Println("node: ", id.Node(), "step: ", id.Step(), "time: ", id.Time(), "\n")
	}
}
