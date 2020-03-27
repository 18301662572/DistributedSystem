package collypub
import (
	"net/url"
	"github.com/gocolly/colly"
	"github.com/nats-io/nats.go"
	"time"
	"os"
)

//nats 结合 colly 的消息生产
//我们为每一个网站定制一个对应的 collector，并设置相应的规则，
//比如 www.v2ex.com，www.v2fx.com(虚构的网站)，再用简单的工厂方法来将该 collector 和其 host 对应起来：

var domain2Collector = map[string]*colly.Collector{}

var nc *nats.Conn

var maxDepth = 10

var natsURL = "nats://localhost:4222"

func factory(urlStr string) *colly.Collector {
	u, _ := url.Parse(urlStr)
	return domain2Collector[u.Host]
}

func initV2exCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("www.v2ex.com"),
		colly.MaxDepth(maxDepth),
	)
	c.OnResponse(func(resp *colly.Response) {
		// 做一些爬完之后的善后工作
		// 比如页面已爬完的确认存进 MySQL
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// 基本的反爬虫策略
		time.Sleep(time.Second * 2)
		// TODO, 正则 match 列表页的话，就 visit
		// TODO, 正则 match 落地页的话，就发消息队列
		//err = nc.Publish("tasks", []byte("your task content"))
		//nc.Flush()
		var link string
		c.Visit(e.Request.AbsoluteURL(link))
	})
	return c
}

func initV2fxCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("www.v2fx.com"),
		colly.MaxDepth(maxDepth),
	)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	})
	return c
}

func init() {
	domain2Collector["www.v2ex.com"] = initV2exCollector()
	domain2Collector["www.v2fx.com"] = initV2fxCollector()
	var err error
	nc, err = nats.Connect(natsURL)
	if err != nil {
		// log fatal
		os.Exit(1)
	}

}

func main() {
	urls := []string{"https://www.v2ex.com", "https://www.v2fx.com"}
	for _, url := range urls {
		instance := factory(url)
		instance.Visit(url)
	}
}
