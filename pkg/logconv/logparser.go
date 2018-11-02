package logconv

import (
	"fmt"
	"regexp"
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
	// Todo (mk): Not totally sure how to do posative lookaheads in Go...
	// just going to remove the comma for now
	remoteAddressRegex := regexp.MustCompile(`- ([0-9]{1,3}\.){3}[0-9]{0,3}`)
	remoteAddress := remoteAddressRegex.FindString(log)
	remoteAddress = strings.TrimPrefix(remoteAddress, "- ")
	if remoteAddress == "" {
		return nil, fmt.Errorf("Could not find remote address (%s)", log)
	}

	requestRegex := regexp.MustCompile(`\"(GET|POST|PUT|PATCH|DELETE|HEAD|CONNECT|OPTIONS|TRACE) /.* \w+/\d\.\d\" \d{3}`)
	request := requestRegex.FindString(log)
	if request == "" {
		return nil, fmt.Errorf("Could not find request information (%s)", log)
	}

	requestFields := strings.Split(request, " ")
	if len(requestFields) != 4 {
		return nil, fmt.Errorf("unexpected length of request fields (%s)", requestFields)
	}

	route := requestFields[1]
	statusCode, err := strconv.Atoi(requestFields[3])

	if err != nil {
		return nil, fmt.Errorf("Unexpected value for status code, %s", requestFields[3])
	}

	return &ReqDetail{
		remoteAddr: remoteAddress,
		statusCode: statusCode,
		route:      route,
	}, nil
}

// FakeLogParser can be used to simplify unit and integration tests.
// It assumes logs are in the following format:
// <remote_address>,<route>,<status_code>
// example: 123.222.11.32,/order,402
type FakeLogParser struct{}

func (nginxLogParser *FakeLogParser) Parse(log string) (*ReqDetail, error) {
	fields := strings.Split(log, ",")
	_, routeIndex, statusCodeIndex := 0, 1, 2
	if len(fields) < statusCodeIndex {
		return nil, fmt.Errorf("Unexpected log format (%s)", log)
	}

	// Todo(mk): Not totally sure how to do posative lookaheads in Go...
	// just going to remove the comma for now
	remoteAddressRegex := regexp.MustCompile(`([0-9]{1,3}\.){3}[0-9]{1,3},`)
	// remoteAddressRegex := regexp.MustCompile(`([0-9]{1,3}\.){3}[0-9]{0,1}-([0-9]{1,3}\.){3}[0-9]`)

	remoteAddress := remoteAddressRegex.FindString(log)
	remoteAddress = strings.TrimSuffix(remoteAddress, ",")
	if remoteAddress == "" {
		return nil, fmt.Errorf("Could not find remote address (%s)", log)
	}

	statusCode, err := strconv.Atoi(fields[statusCodeIndex])
	if err != nil {
		return nil, fmt.Errorf("Unexpected value for status code, %s", fields[statusCodeIndex])
	}
	return &ReqDetail{
		remoteAddr: remoteAddress,
		statusCode: statusCode,
		route:      fields[routeIndex],
	}, nil
}
