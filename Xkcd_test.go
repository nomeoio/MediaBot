package mediabot

import (
	"os"
	"testing"

	"github.com/naughtymonsta/utilities"
)

func TestXKSend(t *testing.T) {
	var mbs utilities.MessageBlocks
	var err error
	mbs, err = xk.GetStoryById("614")
	if err != nil {
		t.Fatal(err)
	}
	err = SC.SendBlocks(mbs, os.Getenv("SlackWebHookUrlTest"))
	if err != nil {
		t.Fatal(err)
	}
}

// func TestJsonInt(t *testing.T) {
// 	var lastID int
// 	_ = json.Unmarshal(utilities.ReadFile(xkcdFilename), &lastID)
// 	t.Log("M:", fmt.Sprintf("%d", lastID))
// 	j, _ := json.Marshal(lastID)
// 	utilities.WriteFile(j, xkcdFilename)
// }
