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

type LogConv interface {
	Start() error
	Stop() error
}

func NewLogConv(config LogConvConf) (LogConv, error) {
	switch config.Type {
	case LogConvBatchType:
		return &BatchLogConv{
			batchTime:        config.BatchTime,
			serverType:       config.ServerType,
			inputLogFilePath: config.InputLogFilePath,
			isRunning:        false,
		}, nil
	default:
		return nil, fmt.Errorf("%s LogConv not supported", config.ServerType)
	}
}

type BatchLogConv struct {
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

func (blc *BatchLogConv) Start() error {
	blc.quitChannel = make(chan bool)
	blc.reqDetailChannel = make(chan *ReqDetail)
	blc.reqDetailConsumer = &ReqDetailConsumer{
		reqDetailChannel: blc.reqDetailChannel,
		quitChannel:      blc.quitChannel,
		store:            make(map[string]int),
		storeLock:        sync.RWMutex{},
	}

	go blc.reqDetailConsumer.Subscribe()

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
	if err != nil {
		return err
	}

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

func (blc *BatchLogConv) Stop() error {
	if blc.isRunning {
		close(blc.quitChannel)
		return nil
	}
	return fmt.Errorf("Can't stop what hasn't started...")
}
