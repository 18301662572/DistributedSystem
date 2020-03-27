package main

import "math/rand"

//基于洗牌算法的负载均衡

//考虑到我们需要随机选取每次发送请求的 endpoint，同时在遇到下游返回错误时换其它节点重试。
//所以我们设计一个大小和 endpoints 数组大小一致的索引数组，每次来新的请求，我们对索引数组做洗牌，
//然后取第一个元素作为选中的服务节点，如果请求失败，那么选择下一个节点重试，以此类推
//我们循环一遍 slice，两两交换，这个和我们平常打牌时常用的洗牌方法类似。看起来没有什么问题。

var endpoints = []string {
	"100.69.62.1:3232",
	"100.69.62.32:3232",
	"100.69.62.42:3232",
	"100.69.62.81:3232",
	"100.69.62.11:3232",
	"100.69.62.113:3232",
	"100.69.62.101:3232",
}

// 重点在这个 shuffle
//错误的洗牌导致的负载不均衡
func shuffle01(slice []int) {
	for i := 0; i < len(slice); i++ {
		a := rand.Intn(len(slice))
		b := rand.Intn(len(slice))
		slice[a], slice[b] = slice[b], slice[a]
	}
}

//修正洗牌算法
//fisher-yates 算法，主要思路为每次随机挑选一个值，放在数组末尾。然后在 n-1 个元素的数组中再随机挑选一个值，放在数组末尾，
func shuffle02(indexes []int) {
	for i:=len(indexes); i>0; i-- {
		lastIdx := i - 1
		idx := rand.Intn(i)
		indexes[lastIdx], indexes[idx] = indexes[idx], indexes[lastIdx]
	}
}

// Go实现fisher-yates 算法
func shuffle03(n int) []int {
	b := rand.Perm(n)
	return b
}

func request(params map[string]interface{}) error {
	var indexes = []int {0,1,2,3,4,5,6}
	var err error
	shuffle(indexes)
	maxRetryTimes := 3
	idx := 0
	for i := 0; i < maxRetryTimes; i++ {
		err = apiRequest(params, indexes[idx])
		if err == nil {
			break
		}
		idx++
	}
	if err != nil {
		// logging
		return err
	}
	return nil
}
