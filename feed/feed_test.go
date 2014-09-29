package feed

import (
	"testing"
)

func TestParseDate(t *testing.T) {
	testFormat := "Thu, 27 Feb 2014 18:46:18 +0100"
	var timestamp int64 = 1393523178

	date, err := ParseDate(testFormat)
	if err != nil {
		t.Fatal(err)
	}

	if date.Unix() != timestamp {
		t.Fatalf("Expected %d - got %d", timestamp, date.Unix())
	}
}
