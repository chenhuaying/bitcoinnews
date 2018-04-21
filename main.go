package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	//"strconv"
	"time"

	//"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type BitcoinData struct {
	Name              string
	Symbol            string
	MarketCap         int64
	Price             float64
	Volume_24h        int64
	CirculatingSupply int64
	Change_24h        float64
}

const (
	DEFAULT_USER            = "chy"
	DEFAULT_PASSWD          = "123456"
	DEFAULT_IP              = "192.168.56.102"
	DEFAULT_PORT            = "3306"
	DEFAULT_DATABASE        = "test"
	DEFAULT_TABLE           = "news"
	DEFAULT_TIMEOUT         = 30
	DEFAULT_MARKETCAP_FLUSH = 1200

	CONF_KEY_DBUSER          = "dbUser"
	CONF_KEY_DBPASSWD        = "dbPasswd"
	CONF_KEY_DBIP            = "dbIP"
	CONF_KEY_DBPORT          = "dbPort"
	CONF_KEY_DATABASE        = "database"
	CONF_KEY_TIMEOUT         = "timeout"
	CONF_KEY_MARKETCAP_FLUSH = "marketcap_flush"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.seebitcoin")
	viper.AddConfigPath("/etc/seebitcoin")

	viper.SetDefault(CONF_KEY_DBUSER, DEFAULT_USER)
	viper.SetDefault(CONF_KEY_DBPASSWD, DEFAULT_PASSWD)
	viper.SetDefault(CONF_KEY_DBIP, DEFAULT_IP)
	viper.SetDefault(CONF_KEY_DBPORT, DEFAULT_PORT)
	viper.SetDefault(CONF_KEY_DATABASE, DEFAULT_DATABASE)
	viper.SetDefault(CONF_KEY_TIMEOUT, DEFAULT_TIMEOUT)
	viper.SetDefault(CONF_KEY_MARKETCAP_FLUSH, DEFAULT_MARKETCAP_FLUSH)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config error: ", err)
	}

	dbUser := viper.GetString(CONF_KEY_DBUSER)
	dbPasswd := viper.GetString(CONF_KEY_DBPASSWD)
	dbAddr := viper.GetString(CONF_KEY_DBIP)
	dbPort := viper.GetString(CONF_KEY_DBPORT)
	database := viper.GetString(CONF_KEY_DATABASE)
	timeout := viper.GetInt(CONF_KEY_TIMEOUT)
	//mcapFlush := viper.GetInt64(CONF_KEY_MARKETCAP_FLUSH)

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, database)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	f, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open file error:", err)
		os.Exit(2)
	}
	defer f.Close()

	collector := NewCollector(dataSource)

	c := colly.NewCollector()
	c.SetRequestTimeout(time.Duration(timeout) * time.Second)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 30 * time.Second,
	})

	//bitcoinlist := make([]*BitcoinData, 0)
	//timestamp := time.Now().Unix()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
		fmt.Println(string(r.Body))
		any := jsoniter.Get(r.Body, "data")
		for i := 0; i < any.Size(); i++ {
			seo_url := any.Get(i, "seo_url")
			if seo_url.ValueType() != 0 {
				url := "https://toutiao.com" + seo_url.ToString()
				fmt.Println(url)
				data := any.Get(i)
				fmt.Println(data.Get("title").ToString(), data.Get("create_time").ToInt(), data.Get("article_url").ToString(), data.Get("media_name").ToString(), data.Get("media_url").ToString())
				title, ctime, media := data.Get("title").ToString(), data.Get("create_time").ToInt(), data.Get("media_name").ToString()
				collector.Collect(title, media, url, ctime)
			}
		}
	})

	c.OnScraped(func(_ *colly.Response) {
		//bData, _ := json.MarshalIndent(bitcoinlist, "", "  ")
		//bData, _ := json.Marshal(bitcoinlist)
		//f.Write(bData)
	})

	urlFmt := "https://www.toutiao.com/search_content/?offset=%d&format=json&keyword=%E5%8C%BA%E5%9D%97%E9%93%BE&autoload=true&count=20&cur_tab=1&from=search_tab"
	for i := 0; i <= 10; i += 20 {
		fmt.Println(fmt.Sprintf(urlFmt, i))
		err = c.Visit(urlFmt)
		if err != nil {
			fmt.Println("vist coins all error:", err)
		}
	}
}
