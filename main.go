package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/btlike/crawl/spider"
	"gitlab.com/btlike/crawl/utils"
)

func main() {
	utils.Init()
	spider.Run()
}
