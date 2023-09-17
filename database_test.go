package mediabot

import (
	"testing"
)

func TestUpdateXkcd(t *testing.T) {
	for i := 0; i < 3; i++ {
		var item SavedNews = DB.UpdateXkcd()
		t.Logf("item: %+v\n", item)
	}
}

func TestReturnAllRecords(t *testing.T) {
	var items []SavedNews = DB.ReturnAllRecords("xkcd")
	for _, item := range items {
		t.Logf("item: %+v\n", item)
	}
}
