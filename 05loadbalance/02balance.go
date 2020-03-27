package main

import (
	"fmt"
	"math/rand"
	"time"
)

//负载均衡算法效果验证,设置随机种子
//map[0:224436 1:128780 5:129310 6:129194 2:129643 3:129384 4:129253]
//map[6:143275 5:143054 3:143584 2:143031 1:141898 0:142631 4:142527]
//分布结果和我们推导出的结论是一致的。

func init() {
	rand.Seed(time.Now().UnixNano())
}
//shuffle 算法
func shuffle1(slice []int) {
	for i := 0; i < len(slice); i++ {
		a := rand.Intn(len(slice))
		b := rand.Intn(len(slice))
		slice[a], slice[b] = slice[b], slice[a]
	}
}
//fisher yates 算法
func shuffle2(indexes []int) {
	for i := len(indexes); i > 0; i-- {
		lastIdx := i - 1
		idx := rand.Intn(i)
		indexes[lastIdx], indexes[idx] = indexes[idx], indexes[lastIdx]
	}
}
func main() {
	var cnt1 = map[int]int{}
	for i := 0; i < 1000000; i++ {
		var sl = []int{0, 1, 2, 3, 4, 5, 6}
		shuffle1(sl)
		cnt1[sl[0]]++
	}
	var cnt2 = map[int]int{}
	for i := 0; i < 1000000; i++ {
		var sl = []int{0, 1, 2, 3, 4, 5, 6}
		shuffle2(sl)
		cnt2[sl[0]]++
	}
	fmt.Println(cnt1, "\n", cnt2)
}
