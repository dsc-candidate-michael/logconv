package logconv

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLogFileObserverStart(t *testing.T) {
	reqDetailChannel := make(chan *ReqDetail)
	quitChannel := make(chan bool)
	fakeLogFilePath := "../../test-artifacts/fake.log"
	fakeLogFileHandle, _ := os.OpenFile(fakeLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	parser, _ := NewLogParser(LogParserConf{
		Type: ServerTypeFake,
	})
	logFileObserverConfig := LogObserverConfig{
		Parser:           parser,
		Type:             LogObserverTypeFile,
		InputFile:        fakeLogFilePath,
		ReqDetailChannel: reqDetailChannel,
		QuitChannel:      quitChannel,
	}
	logFileObserver, _ := NewLogObserver(logFileObserverConfig)

	testLogs := []string{
		"72.34.110.66,/eye,200\n",
		"72.34.43.56,/cream,404\n",
		"72.34.33.12,/onion,200\n",
		"72.34.110.13,/cheese,403\n",
		"72.34.130.14,/potato,500\n",
	}

	expectedResults := []*ReqDetail{
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
	}

	var results []*ReqDetail

	logFileObserver.Start()

	go func() {
		for _, log := range testLogs {
			fakeLogFileHandle.WriteString(log)
		}
	}()
	go func() {
		for i := 0; i < len(expectedResults); i++ {
			select {
			case reqDetail := <-reqDetailChannel:
				results = append(results, reqDetail)
			}
		}
		quitChannel <- true
	}()

	select {
	case <-quitChannel:
		equal := reflect.DeepEqual(results, expectedResults)
		fmt.Printf("%v", results)
		if !equal {
			t.Errorf("Did not receive the expected results")
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout occured waiting for the channel to close")
	}
}
