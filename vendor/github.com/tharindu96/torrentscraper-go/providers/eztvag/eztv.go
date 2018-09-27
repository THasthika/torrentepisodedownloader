package eztvag

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/tharindu96/torrentscraper-go/providers"
)

const id = "eztvag"
const name = "eztv.ag"
const ttype = providers.TorrentTypeTV

const urlPlaceholder = "https://eztv.ag/search/%s"

const colName = 1
const colLink = 2
const colSize = 3
const colSeeds = 5

// Init func
func Init() *providers.Provider {
	return &providers.Provider{
		ID:             id,
		SupportedTypes: ttype,
		Search:         Search,
		SearchShow:     SearchShow,
	}
}

// Search func
func Search(query string, t providers.TorrentType, out chan []*providers.TorrentMeta) {
	if t&ttype != t {
		out <- nil
		return
	}
	ret := search(query)
	out <- ret
}

// SearchShow func
func SearchShow(name string, season uint, episode uint, out chan []*providers.TorrentMeta) {
	var query string
	if episode > 0 {
		query = fmt.Sprintf("%s-s%02de%02d", name, season, episode)
	} else {
		query = fmt.Sprintf("%s-s%02d", name, season)
	}
	ret := search(query)
	newret := make([]*providers.TorrentMeta, 0)
	for _, t := range ret {
		if strings.Contains(t.Name, name) {
			newret = append(newret, t)
		}
	}
	out <- newret
}

func search(query string) []*providers.TorrentMeta {
	ret := make([]*providers.TorrentMeta, 0)

	url := fmt.Sprintf(urlPlaceholder, query)

	doc, err := providers.GetGoQueryDocument(url)
	if err != nil {
		return nil
	}

	table := doc.Find("table.forum_header_border").Last()
	if len(table.Nodes) == 0 {
		return nil
	}
	table.Find("tr.forum_header_border").Each(func(i int, row *goquery.Selection) {
		t := providers.TorrentMeta{}
		row.Find("td.forum_thread_post").Each(func(j int, col *goquery.Selection) {
			switch j {
			case colName:
				t.Name = col.Find("a").Text()
				break
			case colLink:
				t.Magnet, _ = col.Find("a.magnet").Attr("href")
				break
			case colSeeds:
				s, err := strconv.ParseUint(col.Text(), 10, 32)
				if err == nil {
					t.Seeds = uint(s)
				}
				break
			case colSize:
				size, err := providers.ConvertSize(col.Text())
				if err == nil {
					t.Size = size
				}
				break
			}
		})
		ret = append(ret, &t)
	})

	return ret
}
