package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/jeremymwells/easyConfig"
	"github.com/jeremymwells/crawlers/pastebin"
)

type Configuration struct{
	StartAddress string
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
	fetchQueue = Queue{}
	processQueue = Queue{}
	currentAddress = &config.StartAddress //todo: coalesce last crawled address and startingAddress
)







func main(){
	pastebin.Crawl();
	//TODO initialize any/all web crawlers from here
}