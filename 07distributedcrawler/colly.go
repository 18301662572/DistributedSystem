package main
import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"time"
)

//基于 colly 的单机爬虫
//用 colly 来实现爬取 v2ex 前十页内容：
//直接命令行爬到 v2ex 在 Go tag 下的新贴，

var visited = map[string]bool{}
func main() {
	// Instantiate default collector 实例化默认收集器
	c := colly.NewCollector(
		colly.AllowedDomains("www.v2ex.com"),
		colly.MaxDepth(1),
	)
	detailRegex, _ := regexp.Compile(`/go/go\?p=\d+$`)
	listRegex, _ := regexp.Compile(`/t/\d+#\w+`)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// 已访问过的详情页或列表页，跳过
		if visited[link] && (detailRegex.Match([]byte(link)) || listRegex.Match([]byte(link))) {
			return
		}
		// 匹配下列两种 url 模式的，才去 visit
		// https://www.v2ex.com/go/go?p=2
		// https://www.v2ex.com/t/472945#reply3
		if !detailRegex.Match([]byte(link)) && !listRegex.Match([]byte(link)) {
			println("not match", link)
			return
		}
		time.Sleep(time.Second)
		println("match", link)
		println(e.Request.AbsoluteURL(link))
		visited[link] = true
		time.Sleep(time.Millisecond * 2)
		c.Visit(e.Request.AbsoluteURL(link))
	})
	err := c.Visit("https://www.v2ex.com/go/go")
	if err != nil {
		fmt.Println(err)
	}
}
