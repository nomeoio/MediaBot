package mediabot

import (
	"encoding/json"
	"fmt"

	"github.com/naughtymonsta/utilities"
)

// XKCD
type XKCDClient struct{}

type xkcdJSON struct {
	Month      string `json:"month,omitempty"`
	Num        int    `json:"num,omitempty"`
	Link       string `json:"link,omitempty"`
	Year       string `json:"year,omitempty"`
	News       string `json:"news,omitempty"`
	SafeTitle  string `json:"safe_title,omitempty"`
	Transcript string `json:"transcript,omitempty"`
	Alt        string `json:"alt,omitempty"`
	Img        string `json:"img,omitempty"`
	Title      string `json:"title,omitempty"`
	Day        string `json:"day,omitempty"`
}

// func init() {
// 	db.UpdateXkcd()
// }

func (xk XKCDClient) RetrieveJsonById(id string) (xkj xkcdJSON, err error) {
	var fStr string = "https://xkcd.com/%s/info.0.json"
	var body []byte
	if body, err = utilities.HttpRequest("GET", nil, fmt.Sprintf(fStr, id), nil); err != nil {
		return
	}
	if err = json.Unmarshal(body, &xkj); err != nil {
		return
	}
	return
}

func (xk XKCDClient) GetStoryById(id string) (mbs utilities.MessageBlocks, err error) {
	var xkj xkcdJSON
	xkj, err = xk.RetrieveJsonById(id)
	mbs = utilities.MessageBlocks{
		Blocks: []utilities.MessageBlock{
			{Type: "divider"},
			{
				Type: "section",
				Text: &utilities.ElementText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*%s* [No.%s] <%s-%s-%s>\n%s", xkj.SafeTitle, id, xkj.Year, xkj.Month, xkj.Day, xkj.Transcript),
				},
			},
			{
				Type:     "image",
				ImageUrl: xkj.Img,
				AltText:  xkj.Alt,
			},
		},
	}
	return
}
