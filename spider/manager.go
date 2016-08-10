package spider

import (
	"fmt"
	"sync"
	"time"

	"github.com/btlike/crawl/utils"
	"github.com/btlike/spider"
)

//const
const (
	HashChanSize  = 512
	UniqHashSize  = 100 * 1000
	BlackListSize = 100 * 1000
	StoreChanSize = 1000
)

//Manager spider
var manage manager

func (p *manager) run() {
	p.initChan()
	p.initUniqHash()
	p.initBlackList()
	p.wire = NewWire()
	go p.monitor()
}

func (p *manager) initChan() {
	p.hashIDChan = make(chan spider.AnnounceData, HashChanSize)
	p.storeMap = make(map[string]chan string)
	for i := 0; i <= 15; i++ {
		p.storeMap[fmt.Sprintf("%X", i)] = make(chan string, StoreChanSize)
	}
}

func (p *manager) initBlackList() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for k := range p.blackList {
		delete(p.blackList, k)
	}
	p.blackList = nil
	p.blackList = make(map[string]bool, BlackListSize)
}

func (p *manager) initUniqHash() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for k := range p.uniqInfohash {
		delete(p.uniqInfohash, k)
	}
	p.uniqInfohash = nil
	p.uniqInfohash = make(map[string]bool, UniqHashSize)
}

func (p *manager) isHashinfoExist(hash string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.uniqInfohash[hash]; ok {
		return true
	}
	manage.uniqInfohash[hash] = true
	return false
}

type manager struct {
	wire         *Wire
	storeCount   int64
	storeMap     map[string]chan string
	uniqInfohash map[string]bool
	blackList    map[string]bool
	hashIDChan   chan spider.AnnounceData
	mutex        sync.Mutex
}

func (p *manager) monitor() {
	go spider.Monitor()

	go func() {
		for {
			utils.Log().Printf("从DHT网络获取资源数量(BEP9): %v\n", p.storeCount)
			if len(manage.uniqInfohash) >= UniqHashSize {
				p.initUniqHash()
			}
			if len(manage.blackList) >= BlackListSize {
				p.initBlackList()
			}
			time.Sleep(time.Minute)
		}
	}()
}
