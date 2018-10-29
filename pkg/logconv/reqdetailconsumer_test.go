package logconv

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestReqDetailConsumer(t *testing.T) {
	reqDetailChannel := make(chan *ReqDetail)
	quitChannel := make(chan bool)

	reqDetails := []*ReqDetail{
		{
			remoteAddr: "72.34.110.66",
			statusCode: 200,
			route:      "/eye",
		},
		{
			remoteAddr: "72.34.43.56",
			statusCode: 404,
			route:      "/cream",
		},
		{
			remoteAddr: "72.34.33.12",
			statusCode: 200,
			route:      "/onion",
		},
		{
			remoteAddr: "72.34.110.13",
			statusCode: 403,
			route:      "/cheese",
		},
		{
			remoteAddr: "72.34.130.14",
			statusCode: 500,
			route:      "/potato",
		},
		{
			remoteAddr: "72.34.130.14",
			statusCode: 500,
			route:      "/potato",
		},
		{
			remoteAddr: "72.34.130.14",
			statusCode: 500,
			route:      "/cheese",
		},
	}

	expectedResult := map[string]int{
		Key20xStatusCode: 2,
		Key30xStatusCode: 0,
		Key40xStatusCode: 2,
		Key50xStatusCode: 3,
		"/potato":        2,
		"/cheese":        1,
	}

	reqDetailConsumer := &ReqDetailConsumer{
		reqDetailChannel: reqDetailChannel,
		quitChannel:      quitChannel,
		store:            make(map[string]int),
		storeLock:        sync.RWMutex{},
	}

	go reqDetailConsumer.Subscribe()

	for _, reqDetail := range reqDetails {
		reqDetailChannel <- reqDetail
	}
	time.Sleep(250 * time.Millisecond)
	results := reqDetailConsumer.Flush()
	equal := reflect.DeepEqual(results, expectedResult)
	if !equal {
		t.Errorf("Flush did not return the expected result (%v)", results)
	}
}
