package main

import (
	"fmt"
	"github.com/jeremymwells/torCrawler/database"
)

var(
	db = database.Get()
)

func main() {
	
	fmt.Print(database.Get());
}