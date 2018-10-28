package logconv

import (
	"fmt"
)

type LogConvType string

const (
	LogConvBatchType LogConvType = "Batch"
)

type LogConvConf struct {
	InputLogFilePath string
	BatchTime        int64
	Type             LogConvType
	ServerType       ServerType
}

type LogConv interface {
	Run()
	Stop()
}

func NewLogConv(config LogConvConf) (LogConv, error) {
	switch config.Type {
	case LogConvBatchType:
		return &BatchLogConv{
			batchTime:        config.BatchTime,
			serverType:       config.ServerType,
			inputLogFilePath: config.InputLogFilePath,
		}, nil
	default:
		return nil, fmt.Errorf("%s LogConv not supported", config.ServerType)
	}
}

type BatchLogConv struct {
	batchTime        int64
	serverType       ServerType
	inputLogFilePath string
}

func (blc *BatchLogConv) Run() {
	fmt.Println("running")
}

func (blc *BatchLogConv) Stop() {
	fmt.Println("stopping")
}
