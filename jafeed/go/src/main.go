package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	r "gopkg.in/dancannon/gorethink.v2"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"bufio"
	"os"
)

type FeedEntry struct {
	Id string `gorethink:"id,omitempty"`
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
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
		Database: "jafeed",
	})
    if err != nil {
        log.Fatalln(err.Error())
    }
    return session
}

func parseitems(feed *gofeed.Feed, items []*gofeed.Item, c chan string) string {
	conn := Connect()
	var allitems []*FeedEntry
	for _, item := range items {
		feeddata := new(FeedEntry)
		// hash an id and package data
		idstr := fmt.Sprintf("%s %s %s", item.Published, item.Title, item.Author.Name)
		hasher := md5.New()
		hasher.Write([]byte(idstr))
		feeddata.Id = hex.EncodeToString(hasher.Sum(nil))
		feeddata.Feed = feed.Link
		feeddata.Title = item.Title
		feeddata.Content = item.Content
		feeddata.Description = item.Description
		feeddata.Url = item.Link
		feeddata.Published = item.Published
		allitems = append(allitems, feeddata)
	}
	// fire the write query into Rethinkdb
	changes, err := r.Table("feeds").Insert(allitems, r.InsertOpts{Conflict: "update"}).RunWrite(conn)
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
	c := make(chan string)
	for _, url := range urls {
		txt := fmt.Sprintf("Parsing feed %s", url)
		fmt.Println(txt)
		feed := parsefeed(url)
		fmt.Println(feed.Title)
		items := feed.Items
		go parseitems(feed, items, c)
		output := <- c
		fmt.Println(output)
	}
	fmt.Println("Ok")
	return
}

