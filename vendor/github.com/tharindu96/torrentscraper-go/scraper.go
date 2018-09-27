package torrentscraper

import (
	"sort"

	"github.com/tharindu96/torrentscraper-go/providers"
	"github.com/tharindu96/torrentscraper-go/providers/eztvag"
	"github.com/tharindu96/torrentscraper-go/providers/zooqlecom"
)

// Scraper struct
type Scraper struct {
	Providers map[string]*providers.Provider
}

// New creates a new scraper
func New(excludeProviders ...string) *Scraper {
	s := &Scraper{
		Providers: make(map[string]*providers.Provider),
	}

	registerScraper(s, eztvag.Init(), &excludeProviders)
	registerScraper(s, zooqlecom.Init(), &excludeProviders)

	return s
}

// Search func
func (scraper *Scraper) Search(query string) *Result {
	return scraper.searchType(query, providers.TorrentTypeUnspecified)
}

// SearchMovie func
func (scraper *Scraper) SearchMovie(query string) *Result {
	return scraper.searchType(query, providers.TorrentTypeMovie)
}

// SearchShow func
func (scraper *Scraper) SearchShow(name string, season uint, episode uint) *Result {
	out := make(chan []*providers.TorrentMeta)

	count := len(scraper.Providers)
	for _, p := range scraper.Providers {
		if p.SearchShow == nil {
			count--
		} else {
			go p.SearchShow(name, season, episode, out)
		}
	}

	res := getTorrentList(out, count)

	close(out)

	ret := &Result{
		Torrents: res,
	}

	return ret
}

func (scraper *Scraper) searchType(query string, ttype providers.TorrentType) *Result {
	out := make(chan []*providers.TorrentMeta)

	count := len(scraper.Providers)
	for _, p := range scraper.Providers {
		if p.Search == nil {
			count--
		} else {
			go p.Search(query, ttype, out)
		}
	}

	res := getTorrentList(out, count)

	close(out)

	ret := &Result{
		Torrents: res,
	}

	return ret
}

func getTorrentList(out chan []*providers.TorrentMeta, count int) []*providers.TorrentMeta {
	tmap := make(map[string]*providers.TorrentMeta)

	for i := 0; i < count; i++ {
		result := <-out
		if result != nil {
			mergeTorrents(result, &tmap)
		}
	}

	res := make([]*providers.TorrentMeta, len(tmap))
	i := 0
	for _, v := range tmap {
		res[i] = v
		i++
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Seeds < res[j].Seeds {
			return false
		}
		return true
	})
	return res
}

func mergeTorrents(torrents []*providers.TorrentMeta, tmap *map[string]*providers.TorrentMeta) {
	for _, t := range torrents {
		hash, err := providers.GetTorrentXt(t.Magnet)
		if err != nil {
			continue
		}
		if x, ok := (*tmap)[hash]; ok {
			if x.Seeds < t.Seeds {
				(*tmap)[hash] = t
			}
		} else {
			(*tmap)[hash] = t
		}
	}
}

func registerScraper(scraper *Scraper, provider *providers.Provider, excludeProviders *[]string) {
	found := false
	for _, p := range *excludeProviders {
		if provider.ID == p {
			found = true
		}
	}
	if !found {
		scraper.Providers[provider.ID] = provider
	}
}
