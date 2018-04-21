package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

const (
	TITLE_PREFIX   = "title:"
	CONTENT_PREFIX = "content:"
)

type Task struct {
	Title string
	Ctime int
	Media string
	Url   string
}

type Collector struct {
	queue   chan *Task
	db      *sql.DB
	handler *colly.Collector
	data    string
}

func NewCollector(dbSource string) *Collector {
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err.Error())
	}

	collector := &Collector{
		queue: make(chan *Task),
		db:    db,
	}

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*toutiao.*",
		RandomDelay: 2 * time.Minute,
	})

	// On every a element which has href attribute call callback
	c.OnHTML("body", func(e *colly.HTMLElement) {
		//fmt.Println(string(e.Response.Body))
		//fmt.Println("-----------------")
		//fmt.Println(e.Text)
		bodyBytes := []byte(e.Text)[:]
		ai := bytes.Index(bodyBytes, []byte("articleInfo: {"))
		if ai != -1 {
			aiBytes := bodyBytes[ai:]

			ti := bytes.Index(aiBytes, []byte(TITLE_PREFIX))
			tb := aiBytes[ti+len(TITLE_PREFIX):]
			comma := bytes.Index(tb, []byte("',"))
			tb = bytes.TrimSpace(tb[:comma])[1:]

			tmp := aiBytes[ti:]
			ci := bytes.Index(tmp, []byte(CONTENT_PREFIX))
			cb := tmp[ci+len(CONTENT_PREFIX):]
			comma = bytes.Index(cb, []byte("',"))
			cb = bytes.TrimSpace(cb[:comma])[1:]
			fmt.Println(string(tb), string(cb))

			unesSrc := html.UnescapeString(string(cb))
			buf := bytes.NewBufferString(unesSrc)
			doc, err := goquery.NewDocumentFromReader(buf)
			if err != nil {
				fmt.Println(err)
			} else {
				c := doc.Contents().Text()
				fmt.Println(">>>>>>>", c)
				collector.data = c
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Headers.Set("accept-language", "zh-CN,zh;q=0.9")
		r.Headers.Set("accept-encoding", "deflate, br")
		r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("upgrade-insecure-requests", "1")
		r.Headers.Set("cache-control", "max-age=0")
		r.Headers.Set("cookie", "tt_webid=6540885012588267022; UM_distinctid=15ffc32d2b5a0c-0a3d9ec3eb9e0d-173f6d55-fa000-15ffc32d2b64a6; tt_webid=6540885012588267022; WEATHER_CITY=%E5%8C%97%E4%BA%AC; tt_webid=6540885012588267022; sso_login_status=0; __tasessionId=n0v527z0j1523804124171; CNZZDATA1259612802=2047701018-1511763536-%7C1523803829")
		r.Headers.Set("referer", "https://www.toutiao.com/search/?keyword=%E5%8C%BA%E5%9D%97%E9%93%BE")
		r.Headers.Set("connection", "keep-alive")
		//fmt.Println(r.Headers)
	})

	collector.handler = c

	return collector
}

func (c *Collector) Collect(title, media, url string, ctime int) {

	c.handler.Visit(url)

	insert, err := c.db.Exec("insert into cc_news (title, media, url, ctime, time, content) values (?, ?, ?, FROM_UNIXTIME(?), NOW(), ?)",
		title, media, url, ctime, c.data)
	if err == nil {
		_, err = insert.LastInsertId()
		fmt.Println("insert cc_news:", title, media, url, ctime)
	}
	c.data = ""
	fmt.Println("insert cc_new error:", err)
}
