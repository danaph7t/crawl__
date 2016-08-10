package spider

import (
	"fmt"
	"sync"
	"time"

	"github.com/btlike/spider"
)

//const
const (
	HashChanLength  = 512
	UniqHashSize    = 100 * 1000
	StoreChanLength = 1000
)

//Manager spider
var manage manager

func (p *manager) run() {
	p.initChan()
	p.initUniqHash()
	p.wire = NewWire()
	go p.monitor()
}

func (p *manager) initChan() {
	p.hashIDChan = make(chan spider.AnnounceData, HashChanLength)
	p.storeMap = make(map[string]chan string)
	for i := 0; i <= 15; i++ {
		p.storeMap[fmt.Sprintf("%X", i)] = make(chan string, StoreChanLength)
	}
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
	storeMap     map[string]chan string
	uniqInfohash map[string]bool
	hashIDChan   chan spider.AnnounceData
	mutex        sync.Mutex
}

func (p *manager) monitor() {
	go spider.Monitor()

	go func() {
		for {
			if len(manage.uniqInfohash) >= UniqHashSize {
				p.initUniqHash()
			}
			time.Sleep(time.Minute)
		}
	}()
}
