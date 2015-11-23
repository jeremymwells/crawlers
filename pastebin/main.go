package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/jeremymwells/easyConfig"
	"github.com/jeremymwells/crawlers/database"
	"github.com/jeremymwells/crawlers/hasher"
)

type Configuration struct{
	StartAddress string
}

type Link struct {
	Value string
	Next *Link
}

type Queue struct {
	First *Link
	Last *Link
}

func (this *Queue) Enqueue(link Link) {
	this.Last.Next = &link
	this.Last = this.Last.Next
}

func (this *Queue) Dequeue() *Link {
	first := this.First
	this.First = this.First.Next
	return first
}

var (
	config = easyConfig.New(&Configuration{}, "./config.json").(*Configuration)
	db = database.Get()
	fetchQueue = Queue{}
	processQueue = Queue{}
	currentAddress = &config.StartAddress //todo: coalesce last crawled address and startingAddress
)




func Crawl(){
	
	resp, err := http.Get(*currentAddress)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	
	sha1, md5 := hasher.Hash(body)
	
	file := database.DbFile{0, sha1, md5, int64(len(body)), "text/html"}
	
	pastebinFile := &database.DbPastebinFile{0,file,*currentAddress}
	
	pbf := db.WritePastebinFile(pastebinFile, false)
	
	fmt.Println(pbf)
}