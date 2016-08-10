package spider

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/btlike/crawl/utils"
	"github.com/btlike/database/torrent"
	"github.com/btlike/spider"
)

//Run the spider
func Run() {
	manage.run()

	idList := spider.GenerateIDList(int64(utils.Config.SpiderNumber))

	//开启的dht节点
	for k, id := range idList {
		go func(port int, id spider.ID) {
			dhtNode := spider.NewDhtNode(&id, manage.hashIDChan, fmt.Sprintf(":%v", utils.Config.SpiderListenPort+port))
			dhtNode.Run()
		}(k, id)
	}

	store()

	for result := range manage.hashIDChan {
		if len(result.Infohash) == 40 {
			hash := strings.ToUpper(result.Infohash)
			if manage.isHashinfoExist(hash) {
				continue
			}
			if result.IsAnnouncePeer {
				var r Request
				r.Port = result.Port
				r.IP = result.IP.String()
				r.InfoHash = []byte(result.Infohash)
				manage.wire.fetchMetadata(r)
			}
			receive(hash)
		}
	}
}

//根据infohash的首字符(0~F)，将infohash写入到对应chan中
func receive(hash string) {
	if c, ok := manage.storeMap[string(hash[0])]; ok {
		c <- hash
	}
}

func store() {
	for k, v := range manage.storeMap {
		go storeSingle(k, v)
	}
}

//批量处理爬取到的infohash，如果此infohash已经抓取过了，则资源热度+1，否则存入预处理表
func storeSingle(k string, v chan string) {
	var hashs []string
	for hash := range v {
		hashs = append(hashs, hash)
		if len(hashs) >= 100 {
			table := getTable(k)
			if table != nil {
				err := utils.Config.Engine.In("infohash", hashs).Cols("infohash").Find(table)
				if err != nil {
					utils.Log().Println(err)
					hashs = make([]string, 0)
					continue
				}
				nt, err := json.Marshal(table)
				if err != nil {
					utils.Log().Println(err)
					hashs = make([]string, 0)
					continue
				}
				var data []torrent.Infohash0
				err = json.Unmarshal(nt, &data)
				if err != nil {
					utils.Log().Println(err)
					hashs = make([]string, 0)
					continue
				}

				var (
					hashMap = make(map[string]int)
				)

				for _, item := range hashs {
					hashMap[item] = 0
				}
				for _, item := range data {
					hashMap[item.Infohash]++
				}

				for key, value := range hashMap {
					if value > 0 && len(key) == 40 {
						searchResult, err := utils.Config.ElasticClient.Get().Index("torrent").Type("infohash").Id(key).Do()
						if err == nil && searchResult != nil && searchResult.Source != nil {
							var tdata torrentData
							err = json.Unmarshal(*searchResult.Source, &tdata)
							if err == nil {
								tdata.Heat++
								_, err = utils.Config.ElasticClient.Index().
									Index("torrent").
									Type("infohash").
									Id(key).
									BodyJson(tdata).
									Refresh(false).
									Do()
								if err != nil {
									utils.Log().Println(err)
								}
							}
						}
					} else {
						go StorePreInfohash(key)
					}
				}
				hashs = make([]string, 0)
			}
		}
	}
}

func getTable(k string) interface{} {
	switch k {
	case "0":
		return &[]torrent.Infohash0{}
	case "1":
		return &[]torrent.Infohash1{}
	case "2":
		return &[]torrent.Infohash2{}
	case "3":
		return &[]torrent.Infohash3{}
	case "4":
		return &[]torrent.Infohash4{}
	case "5":
		return &[]torrent.Infohash5{}
	case "6":
		return &[]torrent.Infohash6{}
	case "7":
		return &[]torrent.Infohash7{}
	case "8":
		return &[]torrent.Infohash8{}
	case "9":
		return &[]torrent.Infohash9{}
	case "A":
		return &[]torrent.Infohasha{}
	case "B":
		return &[]torrent.Infohashb{}
	case "C":
		return &[]torrent.Infohashc{}
	case "D":
		return &[]torrent.Infohashd{}
	case "E":
		return &[]torrent.Infohashe{}
	case "F":
		return &[]torrent.Infohashf{}
	}
	return nil
}

type torrentData struct {
	Name       string
	Length     int64
	Heat       int64
	CreateTime time.Time
}

//StorePreInfohash into temp table
func StorePreInfohash(infohash string) (err error) {
	if len(infohash) == 40 {
		var pre torrent.PreInfohash
		pre.Infohash = infohash
		pre.CreateTime = time.Now()
		_, err = utils.Config.Engine.Insert(&pre)
	}
	return
}
