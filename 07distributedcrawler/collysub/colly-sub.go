package collysub
import (
	"net/url"
	"github.com/gocolly/colly"
	"github.com/nats-io/nats.go"
	"time"
	"os"
)

//nats结合 colly 的消息消费

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
	return c
}

func initV2fxCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("www.v2fx.com"),
		colly.MaxDepth(maxDepth),
	)
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

func startConsumer() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		// log error
		return
	}
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
		urlStr := string(msg.Data)
		ins := factory(urlStr)
		// 因为最下游拿到的一定是对应网站的落地页
		// 所以不用进行多余的判断了，直接爬内容即可
		ins.Visit(urlStr)
	}
}
func main() {
	startConsumer()
}