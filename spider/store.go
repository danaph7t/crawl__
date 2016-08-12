package spider

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/btlike/crawl/utils"
	"github.com/btlike/database/torrent"
)

type torrentContent struct {
	Infohash   string
	Name       string
	CreateTime time.Time
	Length     int64
	FileCount  int64

	Files []file
}

type files []file

type file struct {
	Name   string
	Length int64
}

func (a files) Len() int           { return len(a) }
func (a files) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a files) Less(i, j int) bool { return a[i].Length > a[j].Length }

type torrentSearch struct {
	Name       string
	Length     int64
	CreateTime time.Time
}

func storeTorrent(data interface{}, infohash []byte) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	var t torrentContent
	t.CreateTime = time.Now()

	if info, ok := data.(map[string]interface{}); ok {
		//get name
		if name, ok := info["name"].(string); ok {
			t.Name = name
			if t.Name == "" {
				return fmt.Errorf("store name len is 0")
			}
		}
		//get infohash
		t.Infohash = hex.EncodeToString(infohash)
		if len(t.Infohash) != 40 {
			return fmt.Errorf("store infohash len is not 40")
		}
		//get files
		if v, ok := info["files"]; !ok {
			t.Length = int64(info["length"].(int))
			t.FileCount = 1
			t.Files = append(t.Files, file{Name: t.Name, Length: t.Length})
		} else {
			var tmpFiles files
			fls := v.([]interface{})
			tmpFiles = make([]file, len(fls))
			for i, item := range fls {
				fl := item.(map[string]interface{})
				flName := fl["path"].([]interface{})
				tmpFiles[i] = file{
					Name:   flName[0].(string),
					Length: int64(fl["length"].(int)),
				}
			}
			sort.Sort(tmpFiles)

			for k, v := range tmpFiles {
				if len(v.Name) > 0 {
					t.Length += v.Length
					t.FileCount++
					if k < 5 {
						t.Files = append(t.Files, file{
							Name:   v.Name,
							Length: v.Length,
						})
					}
				}
			}
		}
		b, _ := json.Marshal(t)
		if len(string(b)) > 1024 {
			return
		}
		err = insertData(t.Infohash, string(b))
		if err == nil {
			manage.storeCount++
			data := torrentSearch{
				Name:       t.Name,
				Length:     t.Length,
				CreateTime: time.Now(),
			}
			utils.Config.ElasticClient.Index().
				Index("torrent").
				Type("infohash").
				Id(t.Infohash).
				BodyJson(data).
				Refresh(false).
				Do()
		}
	}
	return
}

func insertData(hash string, content string) (err error) {
	switch hash[0] {
	case '0':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash0{Infohash: hash, Data: content})
	case '1':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash1{Infohash: hash, Data: content})
	case '2':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash2{Infohash: hash, Data: content})
	case '3':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash3{Infohash: hash, Data: content})
	case '4':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash4{Infohash: hash, Data: content})
	case '5':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash5{Infohash: hash, Data: content})
	case '6':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash6{Infohash: hash, Data: content})
	case '7':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash7{Infohash: hash, Data: content})
	case '8':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash8{Infohash: hash, Data: content})
	case '9':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash9{Infohash: hash, Data: content})
	case 'A':
		_, err = utils.Config.Engine.Insert(&torrent.Infohasha{Infohash: hash, Data: content})
	case 'B':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashb{Infohash: hash, Data: content})
	case 'C':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashc{Infohash: hash, Data: content})
	case 'D':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashd{Infohash: hash, Data: content})
	case 'E':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashe{Infohash: hash, Data: content})
	case 'F':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashf{Infohash: hash, Data: content})
	}
	return
}

func isTorrentExist(hash string) (has bool) {
	switch hash[0] {
	case '0':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash0{Infohash: hash})
	case '1':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash1{Infohash: hash})
	case '2':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash2{Infohash: hash})
	case '3':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash3{Infohash: hash})
	case '4':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash4{Infohash: hash})
	case '5':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash5{Infohash: hash})
	case '6':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash6{Infohash: hash})
	case '7':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash7{Infohash: hash})
	case '8':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash8{Infohash: hash})
	case '9':
		has, _ = utils.Config.Engine.Get(&torrent.Infohash9{Infohash: hash})
	case 'A':
		has, _ = utils.Config.Engine.Get(&torrent.Infohasha{Infohash: hash})
	case 'B':
		has, _ = utils.Config.Engine.Get(&torrent.Infohashb{Infohash: hash})
	case 'C':
		has, _ = utils.Config.Engine.Get(&torrent.Infohashc{Infohash: hash})
	case 'D':
		has, _ = utils.Config.Engine.Get(&torrent.Infohashd{Infohash: hash})
	case 'E':
		has, _ = utils.Config.Engine.Get(&torrent.Infohashe{Infohash: hash})
	case 'F':
		has, _ = utils.Config.Engine.Get(&torrent.Infohashf{Infohash: hash})
	}
	return
}
