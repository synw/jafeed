package main

import (
	"fmt"
	"log"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"bufio"
	"os"
	"time"
	r "gopkg.in/dancannon/gorethink.v2"
	"github.com/spf13/viper"
	"github.com/mmcdole/gofeed"
)

type FeedEntry struct {
	Id string `gorethink:"id,omitempty"`
	Timestamp int64
	Feed string
	Url string
	Title string
	Published string
	Content string
	Description string
}

func getdomain() string {
	f, _ := os.Open("jafeed.cfg")
    scanner := bufio.NewScanner(f)
    var domain string
    for scanner.Scan() {
		domain = scanner.Text()	
	}
	return domain
}

func getConf() map[string]interface{} {
	viper.SetConfigName("jafeed_config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./etc/jafeed")
	viper.AddConfigPath("$HOME/.jafeed")
	viper.SetDefault("centrifugo_host", "localhost")
	viper.SetDefault("centrifugo_port", 8001)
	viper.SetDefault("rethinkdb_host", "localhost")
	viper.SetDefault("rethinkdb_port", 28015)
	viper.SetDefault("rethinkdb_user", "admin")
	viper.SetDefault("rethinkdb_password", "")
	viper.SetDefault("database", "jafeed")
	viper.SetDefault("table", "feeds")
	viper.SetDefault("frequency", 15)
	err := viper.ReadInConfig()
	if err != nil {
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf := make(map[string]interface{})
	conf["centrifugo_host"] = viper.Get("centrifugo_host")
	conf["centrifugo_port"] = viper.Get("centrifugo_port")
	conf["centrifugo_secret_key"] = viper.Get("centrifugo_secret_key")
	conf["rethinkdb_host"] = viper.Get("rethinkdb_host")
	conf["rethinkdb_port"] = viper.Get("rethinkdb_port")
	conf["rethinkdb_user"] = viper.Get("rethinkdb_user")
	conf["rethinkdb_password"] = viper.Get("rethinkdb_password")
	conf["database"] = viper.Get("database")
	conf["table"] = viper.Get("table")
	conf["frequency"] = viper.GetInt("frequency")
	return conf
}

func geturls() []string {
	// get the list of urls from http
	confurl := fmt.Sprintf("%s%s", getdomain(), "/jafeed/config/")
	response, err := http.Get(confurl)
	var urls []string
    if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			line := scanner.Text()
			urls = strings.Split(line, "#!#")
		}
	}
	return urls
}

func parsefeed(feedurl string) *gofeed.Feed {
	// grab the feed
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedurl)
	return feed
}

func Connect() (*r.Session) {
	conf := getConf()
	host := conf["rethinkdb_host"].(string)
	port := conf["rethinkdb_port"].(string)
	db := conf["database"].(string)
	user := conf["rethinkdb_user"].(string)
	pwd := conf["rethinkdb_password"].(string)
	addr := host+":"+port
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: addr,
		Database: db,
		Username: user,
		Password: pwd,
	})
    if err != nil {
        log.Fatalln(err.Error())
    }
    return session
}
/*
TODO: send info over websocket
func broadcast(conf map[string]interface{}, channel string, message string, c chan string) {
	secret := conf["centrifugo_secret_key"].(string)
	host := conf["centrifugo_host"].(string)
	port := conf["centrifugo_port"].(string)
}
*/
func parseitems(feed *gofeed.Feed, items []*gofeed.Item, c chan string, conf map[string]interface{}) string {
	conn := Connect()
	table := conf["table"].(string)
	var allitems []*FeedEntry
	for _, item := range items {
		feeddata := new(FeedEntry)
		// hash an id and package data
		idstr := fmt.Sprintf("%s %s %s", item.Published, item.Title, item.Author.Name)
		hasher := md5.New()
		hasher.Write([]byte(idstr))
		feeddata.Id = hex.EncodeToString(hasher.Sum(nil))
		feeddata.Timestamp = time.Now().UnixNano()
		feeddata.Feed = feed.Link
		feeddata.Title = item.Title
		feeddata.Content = item.Content
		feeddata.Description = item.Description
		feeddata.Url = item.Link
		feeddata.Published = item.Published
		allitems = append(allitems, feeddata)
	}
	// fire the write query into Rethinkdb
	changes, err := r.Table(table).Insert(allitems, r.InsertOpts{Conflict: "update"}).RunWrite(conn)
	var output string
	if err != nil {
		log.Fatal(err)
	} else {
		output = output+fmt.Sprintf("%d items inserted\n", changes.Inserted)
		output = output+fmt.Sprintf("%d items updated\n", changes.Replaced)
		output = output+fmt.Sprintf("%d items unchanged", changes.Unchanged)
	}
	c <- output
	return ""
}
	
func main() {
	urls := geturls()
	conf := getConf()
	c := make(chan string)
	for _, url := range urls {
		txt := fmt.Sprintf("Parsing feed %s", url)
		fmt.Println(txt)
		feed := parsefeed(url)
		fmt.Println(feed.Title)
		items := feed.Items
		go parseitems(feed, items, c, conf)
		output := <- c
		fmt.Println(output)
	}
	fmt.Println("Ok")
	return
}

