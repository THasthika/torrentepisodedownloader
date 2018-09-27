package torrentscraper

import (
	"strings"

	"github.com/tharindu96/torrentscraper-go/providers"
)

// Result result from the scraper
type Result struct {
	Torrents []*providers.TorrentMeta
}

// FilterMatchAll filters torrents from the result
func (result *Result) FilterMatchAll(keywords ...string) *Result {
	res := make([]*providers.TorrentMeta, 0)
	old := result.Torrents
	for _, t := range old {
		found := true
		for _, k := range keywords {
			if strings.Contains(t.Name, k) != true {
				found = false
				break
			}
		}
		if found {
			res = append(res, t)
		}
	}
	result.Torrents = res
	return result
}

// FilterMatchAny filters torrents from the result
func (result *Result) FilterMatchAny(keywords ...string) *Result {
	res := make([]*providers.TorrentMeta, 0)
	old := result.Torrents
	for _, t := range old {
		for _, k := range keywords {
			if strings.Contains(t.Name, k) {
				res = append(res, t)
				break
			}
		}
	}
	result.Torrents = res
	return result
}

// FilterExcludeAll filters torrents from the result
func (result *Result) FilterExcludeAll(keywords ...string) *Result {
	res := make([]*providers.TorrentMeta, 0)
	old := result.Torrents
	for _, t := range old {
		found := 0
		for _, k := range keywords {
			if strings.Contains(t.Name, k) {
				found++
			}
		}
		if found != len(keywords) {
			res = append(res, t)
		}
	}
	result.Torrents = res
	return result
}

// FilterExcludeAny filters torrents from the result
func (result *Result) FilterExcludeAny(keywords ...string) *Result {
	res := make([]*providers.TorrentMeta, 0)
	old := result.Torrents
	for _, t := range old {
		found := true
		for _, k := range keywords {
			if strings.Contains(t.Name, k) {
				found = false
				break
			}
		}
		if found {
			res = append(res, t)
		}
	}
	result.Torrents = res
	return result
}
