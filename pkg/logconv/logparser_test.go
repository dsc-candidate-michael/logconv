package logconv

import (
	"reflect"
	"testing"
)

func TestValidNginxLogParse(t *testing.T) {
	config := LogParserConf{
		Type: ServerTypeNginx,
	}
	nginxLogParser, err := NewLogParser(config)
	if err != nil {
		t.Fatalf("Error creating nginx log parser (%#v)", err)
	}
	expectedReqDetail := &ReqDetail{
		remoteAddr: "72.34.110.66",
		route:      "/",
		statusCode: 200,
	}
	validNginxLog := "10.10.180.161 - 72.34.110.66, 192.33.28.238 - - - [03/Aug/2015:15:50:06 +0000]  https https https \"GET / HTTP/1.1\" 200 20027 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.107 Safari/537.36\""
	reqDetail, err := nginxLogParser.Parse(validNginxLog)
	if err != nil {
		t.Errorf("Error parsing nginx log (%#v)", err)
	}
	equal := reflect.DeepEqual(expectedReqDetail, reqDetail)
	if !equal {
		t.Errorf("Got unexpected request detail, %v, expected %v", reqDetail, expectedReqDetail)
	}
}

func TestInvalidNginxLogParse(t *testing.T) {
	config := LogParserConf{
		Type: ServerTypeNginx,
	}
	nginxLogParser, err := NewLogParser(config)
	if err != nil {
		t.Fatalf("Error creating nginx log parser (%#v)", err)
	}
	_, err = nginxLogParser.Parse("Some invalid log")
	if err == nil {
		t.Errorf("Expected an error for an invalid log (%s)", err)
	}
}

func TestFakeLogParse(t *testing.T) {
	config := LogParserConf{
		Type: ServerTypeFake,
	}
	fakeLogParser, err := NewLogParser(config)
	if err != nil {
		t.Fatalf("Error creating fake log parser (%#v)", err)
	}

	expectedReqDetail := &ReqDetail{
		remoteAddr: "72.34.110.66",
		route:      "/",
		statusCode: 200,
	}

	reqDetail, err := fakeLogParser.Parse("72.34.110.66,/,200")
	if err != nil {
		t.Errorf("Error parsing fake log (%#v)", err)
	}
	equal := reflect.DeepEqual(expectedReqDetail, reqDetail)
	if !equal {
		t.Errorf("Got unexpected request detail, %v, expected %v", reqDetail, expectedReqDetail)
	}
}
