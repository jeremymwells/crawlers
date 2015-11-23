package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/jeremymwells/easyConfig"
)

type Configuration struct{
	CurrentAddress string
}

var (
	configInstance = easyConfig.New(&Configuration{}, "./config.json").(*Configuration)
	startingAddress = configInstance.CurrentAddress //todo: coalesce last crawled address and startingAddress
)

func main(){
	
	resp, err := http.Get(startingAddress)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}