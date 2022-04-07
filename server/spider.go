package server

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"net/http"
	"os"
	"strings"
)

type Session struct {
	session *colly.Collector
	file    *os.File
	cookie  string
}

func (c *Session) GetSession() *colly.Collector {
	return c.session
}

func (c *Session) Init() *colly.Collector {
	c.session = colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"),
	)

	c.session.AllowURLRevisit = true

	extensions.RandomUserAgent(c.session)

	return c.session
}

func (c *Session) CookieLogin(cookies string) {
	url := "https://zhibo.homeway.com.cn/11605/vips.html"

	c.session.OnRequest(func(request *colly.Request) {
		err := c.session.SetCookies(url, setCookieRaw(cookies))
		if err != nil {
			fmt.Println(err)
		}
	})
}

func setCookieRaw(cookieRaw string) []*http.Cookie {
	var cookies []*http.Cookie
	cookieList := strings.Split(cookieRaw, "; ")
	for _, item := range cookieList {
		key := strings.Split(item, "=")
		name := key[0]
		valueList := key[1:]
		cookieItem := http.Cookie{
			Name:  name,
			Value: valueList[0],
		}
		cookies = append(cookies, &cookieItem)
	}
	return cookies
}

func (c *Session) GetText() {
	url := "https://zhibo.homeway.com.cn/11605/vips.html"

	c.session.OnHTML("#l702945692116623360 > dl > dd > div.ts.ft18", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.session.Visit(url)
}
