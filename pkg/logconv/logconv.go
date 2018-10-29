package logconv

import (
	"fmt"
	"sync"
	"time"
)

type LogConvType string

const (
	LogConvBatchType LogConvType = "Batch"
)

type LogConvConf struct {
	InputLogFilePath string
	BatchTime        int
	Type             LogConvType
	ServerType       string
}

func NewLogConv(config LogConvConf) *LogConv {
	return &LogConv{
		batchTime:        config.BatchTime,
		serverType:       config.ServerType,
		inputLogFilePath: config.InputLogFilePath,
		isRunning:        false,
	}
}

type LogConv struct {
	batchTime         int
	serverType        string
	inputLogFilePath  string
	quitChannel       chan bool
	reqDetailChannel  chan *ReqDetail
	logObserver       LogObserver
	reqDetailConsumer *ReqDetailConsumer
	batchTicker       *time.Ticker
	isRunning         bool
}

func (blc *LogConv) Start() error {
	blc.batchTicker = time.NewTicker(time.Duration(blc.batchTime) * time.Second)
	for {
		select {
		case <-blc.batchTicker.C:
			blc.reqDetailConsumer.Flush()
		case <-blc.quitChannel:
			blc.batchTicker.Stop()
			return nil
		}
	}
}

func (blc *LogConv) Subscribe() error {
	blc.quitChannel = make(chan bool)
	blc.reqDetailChannel = make(chan *ReqDetail)
	blc.reqDetailConsumer = &ReqDetailConsumer{
		reqDetailChannel: blc.reqDetailChannel,
		quitChannel:      blc.quitChannel,
		store:            make(map[string]int),
		storeLock:        sync.RWMutex{},
	}

	go blc.reqDetailConsumer.Subscribe()

	defer blc.Stop()

	parser, err := NewLogParser(LogParserConf{
		Type: blc.serverType,
	})

	if err != nil {
		return err
	}

	logObserverConfig := LogObserverConfig{
		Parser:           parser,
		Type:             LogObserverTypeFile,
		InputFile:        blc.inputLogFilePath,
		ReqDetailChannel: blc.reqDetailChannel,
		QuitChannel:      blc.quitChannel,
	}
	blc.logObserver, err = NewLogObserver(logObserverConfig)
	if err != nil {
		return err
	}

	err = blc.logObserver.Start()
	return err
}

func (blc *LogConv) Stop() error {
	if blc.isRunning {
		blc.isRunning = false
		close(blc.quitChannel)
		return nil
	}
	return fmt.Errorf("Can't stop what hasn't started...")
}
