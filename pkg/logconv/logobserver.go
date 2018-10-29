package logconv

import (
	"fmt"

	"github.com/hpcloud/tail"
)

type LogObserverType string

const (
	LogObserverTypeFile = "LogObserverTypeFile"
)

type LogObserverConfig struct {
	Parser           LogParser
	Type             LogObserverType
	InputFile        string
	ReqDetailChannel chan *ReqDetail
	QuitChannel      chan bool
}

func NewLogObserver(config LogObserverConfig) (LogObserver, error) {
	switch config.Type {
	case LogObserverTypeFile:
		return &LogFileObserver{
			parser:           config.Parser,
			inputFile:        config.InputFile,
			reqDetailChannel: config.ReqDetailChannel,
			quitChannel:      config.QuitChannel,
		}, nil
	default:
		return nil, fmt.Errorf("%s LogObserver not supported", config.Type)
	}
}

type LogObserver interface {
	Start() error
	Stop() error
}

type LogFileObserver struct {
	observer         *tail.Tail
	parser           LogParser
	inputFile        string
	reqDetailChannel chan *ReqDetail
	quitChannel      chan bool
}

func (logFileObserver *LogFileObserver) Start() error {
	observerConfig := tail.Config{
		MustExist: true,
		Follow:    true,
		Logger:    tail.DiscardingLogger,
		ReOpen:    true,
	}
	observer, err := tail.TailFile(logFileObserver.inputFile, observerConfig)
	if err != nil {
		return fmt.Errorf("Could not observe file %s (%v)", logFileObserver.inputFile, err)
	}
	logFileObserver.observer = observer
	go logFileObserver.produceReqDetails()
	return nil
}

// produceReqDetails is blocking function which will watch for new logs and
// will server as a producer of request details.
func (logFileObserver *LogFileObserver) produceReqDetails() error {
	for {
		select {
		case line := <-logFileObserver.observer.Lines:
			// Todo (mk): Consider performing better error handling here
			if line.Err != nil {
				fmt.Printf("there was an error with this line...")
				continue
			}
			reqDetail, err := logFileObserver.parser.Parse(line.Text)
			if err == nil {
				logFileObserver.reqDetailChannel <- reqDetail
			}
		case <-logFileObserver.quitChannel:
			close(logFileObserver.quitChannel)
			return logFileObserver.Stop()
		}
	}
}

func (logFileObserver *LogFileObserver) Stop() error {
	logFileObserver.observer.Stop()
	logFileObserver.observer.Cleanup()
	return nil
}
