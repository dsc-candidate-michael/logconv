package logconv

import (
	"fmt"
	"sync"
)

const (
	Key20xStatusCode = "20x"
	Key30xStatusCode = "30x"
	Key40xStatusCode = "40x"
	Key50xStatusCode = "50x"
)

type ReqDetail struct {
	remoteAddr string
	statusCode int
	route      string
}

func (rd *ReqDetail) RemoteAddr() string {
	return rd.remoteAddr
}

func (rd *ReqDetail) StatusCode() int {
	return rd.statusCode
}

func (rd *ReqDetail) Route() string {
	return rd.route
}

type ReqDetailConsumer struct {
	reqDetailChannel chan *ReqDetail
	quitChannel      chan bool
	// Todo (mk): Consider using sync map
	store     map[string]int
	storeLock sync.RWMutex
}

func (rdc *ReqDetailConsumer) Subscribe() {
	rdc.reset()
	for {
		select {
		case reqDetail := <-rdc.reqDetailChannel:
			fmt.Printf("req detail %v", reqDetail)
			rdc.consume(reqDetail)
		case <-rdc.quitChannel:
			// Producer closes the channel, we don't need to do anything here
			return
		}
	}
}

func (rdc *ReqDetailConsumer) Flush() map[string]int {
	rdc.storeLock.Lock()
	storeCopy := make(map[string]int)
	for key, value := range rdc.store {
		fmt.Printf("%s:%d|s\n", key, value)
		storeCopy[key] = value
	}
	rdc.storeLock.Unlock()
	rdc.reset()
	return storeCopy
}

func (rdc *ReqDetailConsumer) reset() {
	rdc.storeLock.Lock()
	defer rdc.storeLock.Unlock()
	// Create a new store and initialize the values
	rdc.store = make(map[string]int)
	rdc.store[Key20xStatusCode] = 0
	rdc.store[Key30xStatusCode] = 0
	rdc.store[Key40xStatusCode] = 0
	rdc.store[Key50xStatusCode] = 0
}

func (rdc *ReqDetailConsumer) consume(reqDetail *ReqDetail) {
	rdc.storeLock.Lock()
	defer rdc.storeLock.Unlock()
	fmt.Printf("consuming: %v", reqDetail)
	if reqDetail.statusCode < 600 && reqDetail.statusCode >= 500 {
		rdc.store[Key50xStatusCode]++
		rdc.store[reqDetail.route]++
	} else if reqDetail.statusCode < 500 && reqDetail.statusCode >= 400 {
		rdc.store[Key40xStatusCode]++
	} else if reqDetail.statusCode < 400 && reqDetail.statusCode >= 300 {
		rdc.store[Key30xStatusCode]++
	} else if reqDetail.statusCode < 300 && reqDetail.statusCode >= 200 {
		rdc.store[Key20xStatusCode]++
	}
}
