package main

import (
	"github.com/btlike/crawl/spider"
	"github.com/btlike/crawl/utils"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	utils.Init()
	spider.Run()
}
