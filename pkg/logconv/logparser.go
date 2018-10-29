package logconv

import (
	"fmt"
	"strconv"
	"strings"
)

type ServerType string

const (
	ServerTypeNginx = "Nginx"
	ServerTypeFake  = "Fake"
)

type LogParserConf struct {
	Type string
}

type LogParser interface {
	Parse(log string) (*ReqDetail, error)
}

func NewLogParser(config LogParserConf) (LogParser, error) {
	switch config.Type {
	case ServerTypeNginx:
		return &NginxLogParser{}, nil
	case ServerTypeFake:
		return &FakeLogParser{}, nil
	default:
		return nil, fmt.Errorf("%s LogParser not supported", config.Type)
	}
}

type NginxLogParser struct{}

// Todo (mk): Consider making this function more robust. At the moment, it is
// tightly coupled with the exact formatting of the logs.
func (nginxLogParser *NginxLogParser) Parse(log string) (*ReqDetail, error) {
	fields := strings.Split(log, " ")
	remoteAddrIndex, routeIndex, statusCodeIndex := 2, 14, 16
	if len(fields) < statusCodeIndex {
		return nil, fmt.Errorf("Unexpected log format")
	}

	statusCode, err := strconv.Atoi(fields[statusCodeIndex])
	if err != nil {
		return nil, fmt.Errorf("Unexpected value for status code, %s", fields[statusCodeIndex])
	}

	remoteAddr := strings.TrimSuffix(fields[remoteAddrIndex], ",")
	return &ReqDetail{
		remoteAddr: remoteAddr,
		statusCode: statusCode,
		route:      fields[routeIndex],
	}, nil
}

// FakeLogParser can be used to simplify unit and integration tests.
// It assumes logs are in the following format:
// <remote_address>,<route>,<status_code>
// example: 123.222.11.32,/order,402
type FakeLogParser struct{}

func (nginxLogParser *FakeLogParser) Parse(log string) (*ReqDetail, error) {
	fields := strings.Split(log, ",")
	remoteAddrIndex, routeIndex, statusCodeIndex := 0, 1, 2
	if len(fields) < statusCodeIndex {
		return nil, fmt.Errorf("Unexpected log format (%s)", log)
	}

	statusCode, err := strconv.Atoi(fields[statusCodeIndex])
	if err != nil {
		return nil, fmt.Errorf("Unexpected value for status code, %s", fields[statusCodeIndex])
	}
	return &ReqDetail{
		remoteAddr: fields[remoteAddrIndex],
		statusCode: statusCode,
		route:      fields[routeIndex],
	}, nil
}
