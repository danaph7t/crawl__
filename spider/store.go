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
		}
	}
	return
}

func insertData(hash string, content string) (err error) {
	switch hash[0] {
	case '0':
		var data torrent.Infohash0
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '1':
		var data torrent.Infohash1
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '2':
		var data torrent.Infohash2
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '3':
		var data torrent.Infohash3
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '4':
		var data torrent.Infohash4
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '5':
		var data torrent.Infohash5
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '6':
		var data torrent.Infohash6
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '7':
		var data torrent.Infohash7
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '8':
		var data torrent.Infohash8
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '9':
		var data torrent.Infohash9
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'A':
		var data torrent.Infohasha
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'B':
		var data torrent.Infohashb
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'C':
		var data torrent.Infohashc
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'D':
		var data torrent.Infohashd
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'E':
		var data torrent.Infohashe
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'F':
		var data torrent.Infohashf
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	}
	return
}

func getTorrent(hash string) (exist bool, content string) {
	var has bool
	var err error
	switch hash[0] {
	case '0':
		var data torrent.Infohash0
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '1':
		var data torrent.Infohash1
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '2':
		var data torrent.Infohash2
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '3':
		var data torrent.Infohash3
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '4':
		var data torrent.Infohash4
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '5':
		var data torrent.Infohash5
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '6':
		var data torrent.Infohash6
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '7':
		var data torrent.Infohash7
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '8':
		var data torrent.Infohash8
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case '9':
		var data torrent.Infohash9
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case 'A':
		var data torrent.Infohasha
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case 'B':
		var data torrent.Infohashb
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case 'C':
		var data torrent.Infohashc
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case 'D':
		var data torrent.Infohashd
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case 'E':
		var data torrent.Infohashe
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	case 'F':
		var data torrent.Infohashf
		data.Infohash = hash
		has, err = utils.Config.Engine.Get(&data)
		if err != nil {
			return false, ""
		}
		if !has {
			return false, ""
		}
		content = data.Data
	}
	return true, content
}
