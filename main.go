package main

import (
	"fmt"
	"strconv"
	"strings"

	torrentscraper "github.com/tharindu96/torrentscraper-go"
	"github.com/tharindu96/torrentscraper-go/providers"

	"github.com/tharindu96/go-deluge"
)

func main() {

	torrentScraper := torrentscraper.New()

	deluge, err := getDelugeConnection()

	if err != nil {
		panic("Could not connect to Deluge")
	}

	name := setShowName()
	season := setSeasonNumber(name)
	episodes := setEpisodeNumbers(name, season)
	match := setMatchKeywords()
	exclude := setExcludeKeywords()

	if len(match) == 0 {
		match = []string{
			"HDTV",
		}
	}

	if len(exclude) == 0 {
		exclude = []string{
			"1080",
			"720",
		}
	}

	for _, episode := range episodes {
		res := torrentScraper.SearchShow(name, season, episode)
		if len(match) > 0 {
			res = res.FilterMatchAny(match...)
		}
		if len(exclude) > 0 {
			res = res.FilterExcludeAny(exclude...)
		}

		if len(res.Torrents) == 0 {
			fmt.Printf("Warning: no torrent found!\n")
			continue
		}

		torr := res.Torrents[0]

		if torr.Seeds == 0 {
			fmt.Printf("Warning: %s has no seeds\n", torr.Name)
		}

		_, err := deluge.CoreAddTorrentMagnet(torr.Magnet, map[string]interface{}{})
		if err != nil {
			fmt.Printf("Could not add torrent: %s\n", err.Error())
		}

		fmt.Printf("%s has been added\n", res.Torrents[0].Name)
	}

}

func getDelugeConnection() (*deluge.Deluge, error) {
	var dhost = "localhost"
	var dport uint = 8112
	var dpass string

	fmt.Printf("Deluge Host (localhost): ")
	fmt.Scanln(&dhost)

	fmt.Printf("Deluge Port (8112): ")
	fmt.Scanln(&dport)

	fmt.Printf("Deluge Password: ")
	fmt.Printf("\033[8m")
	fmt.Scanln(&dpass)
	fmt.Printf("\033[28m")

	d, err := deluge.New(fmt.Sprintf("http://%s:%d//json", dhost, dport), dpass)

	return d, err
}

func setShowName() string {
	var name string
	fmt.Printf("Show Name: ")
	fmt.Scanln(&name)
	return name
}

func setSeasonNumber(name string) uint {
	var season uint
	fmt.Printf("%s Season: ", name)
	fmt.Scanln(&season)
	return season
}

func setEpisodeNumbers(name string, season uint) []uint {
	var episodes = make([]uint, 0)
	var l string
	var t []string
	fmt.Printf("%s %d Episodes: ", name, season)
	fmt.Scanln(&l)
	t = strings.Split(l, "-")
	if len(t) > 1 {
		min32, err1 := strconv.ParseUint(strings.TrimSpace(t[0]), 10, 32)
		max32, err2 := strconv.ParseUint(strings.TrimSpace(t[1]), 10, 32)
		if err1 != nil || err2 != nil {
			return episodes
		}
		min := uint(min32)
		max := uint(max32)
		episodes = makeRange(min, max)
		return episodes
	}
	t = strings.Split(l, ",")
	if len(t) > 0 {
		for _, i := range t {
			u32, err := strconv.ParseUint(strings.TrimSpace(i), 10, 32)
			if err != nil {
				continue
			}
			episodes = append(episodes, uint(u32))
		}
		return episodes
	}
	return episodes
}

func setMatchKeywords() []string {
	fmt.Printf("Match Keywords: ")
	return getKeywords()
}

func setExcludeKeywords() []string {
	fmt.Printf("Exclude Keywords: ")
	return getKeywords()
}

func getKeywords() []string {
	var l string
	var keywords = make([]string, 0)
	fmt.Scanln(&l)
	for _, k := range strings.Split(l, ",") {
		x := strings.TrimSpace(k)
		if x != "" {
			keywords = append(keywords, x)
		}
	}
	return keywords
}

func selectResult(res []*providers.TorrentMeta) uint {

	fmt.Printf("%2s | %-30s | %-10s | %-5s\n", "id", "name", "size", "seeds")

	for i, t := range res {
		fmt.Printf("%2d | %.30s | %10.2f | %5d\n", i, t.Name, float32(t.Size/(1024*1024)), t.Seeds)
	}

	fmt.Printf("Select Torrent: ")

	var ret uint

	fmt.Scan(&ret)

	return ret
}

func makeRange(min, max uint) []uint {
	a := make([]uint, max-min+1)
	for i := range a {
		a[i] = min + uint(i)
	}
	return a
}
