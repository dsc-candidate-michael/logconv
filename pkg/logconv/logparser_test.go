package logconv

import (
	"bufio"
	"log"
	"os"
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
		remoteAddr: "182.118.53.206",
		route:      "/onion/cheese",
		statusCode: 200,
	}
	validNginxLog := "10.10.180.40 - 182.118.53.206 - - - - - - [03/Aug/2015:15:50:06 +0000]  https https https \"GET /onion/cheese HTTP/1.1\" 200 20027 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.107 Safari/537.36\""
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

// This test just ensures that we do not return any errors when
// parsing the sample logs that were provided.
func TestSampleNginxLogParse(t *testing.T) {
	config := LogParserConf{
		Type: ServerTypeNginx,
	}
	nginxLogParser, _ := NewLogParser(config)
	file, err := os.Open("../../test-data/sample.log")
	if err != nil {
		t.Errorf("Couldn't open test file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := nginxLogParser.Parse(scanner.Text())
		if err != nil {

			t.Errorf("Generated an error running through sample log (%v)", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
