package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	torrentscraper "github.com/tharindu96/torrentscraper-go"
	"github.com/tharindu96/torrentscraper-go/providers"

	"github.com/anacrolix/torrent"

	_ "github.com/gdm85/go-libdeluge"
)

var torrentClient *torrent.Client
var torrentScraper *torrentscraper.Scraper
var activeTorrents map[string]*torrent.Torrent
var showProgressFlag bool = true

func main() {

	activeTorrents = make(map[string]*torrent.Torrent)

	torrentClient, _ = torrent.NewClient(nil)
	defer torrentClient.Close()

	torrentScraper = torrentscraper.New()

	var name string
	var season uint
	var episode uint
	var match []string
	var exclude []string
	var raw_match string
	var raw_exclude string

	flag.StringVar(&name, "name", "", "Name of the show")
	flag.UintVar(&season, "season", 0, "Season number")
	flag.UintVar(&episode, "episode", 0, "Episode number")
	flag.StringVar(&raw_match, "match", "", "Match Keywords")
	flag.StringVar(&raw_exclude, "exclude", "", "Exclude Keywords")

	flag.StringVar(&name, "n", "", "Name of the show")
	flag.UintVar(&season, "s", 0, "Season number")
	flag.UintVar(&episode, "e", 0, "Episode number")
	flag.StringVar(&raw_match, "ma", "", "Match Keywords")
	flag.StringVar(&raw_exclude, "ex", "", "Exclude Keywords")

	flag.Parse()

	for _, k := range strings.Split(raw_match, ",") {
		x := strings.TrimSpace(k)
		if x != "" {
			match = append(match, x)
		}
	}

	for _, k := range strings.Split(raw_exclude, ",") {
		x := strings.TrimSpace(k)
		if x != "" {
			exclude = append(exclude, x)
		}
	}

	res := torrentScraper.SearchShow(name, season, episode)
	if len(match) > 0 {
		res = res.FilterMatchAny(match...)
	}
	if len(exclude) > 0 {
		res = res.FilterExcludeAny(exclude...)
	}

	selected := selectResult(res.Torrents)
	for uint(len(res.Torrents)) < selected {
		selected = selectResult(res.Torrents)
	}

	addTorrent(res.Torrents[selected].Name, res.Torrents[selected].Magnet)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _ := reader.ReadString('\n')
			if string([]byte(input)[0]) == "q" {
				showProgressFlag = false
				break
			}
		}
	}()

	showProgress()

}

func addTorrent(name string, magnet string) {
	t, err := torrentClient.AddMagnet(magnet)
	if err != nil {
		return
	}
	activeTorrents[name] = t
	<-t.GotInfo()
	go func(t *torrent.Torrent) {
		fmt.Println("Starting Download...")
		t.DownloadAll()
	}(t)
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

func showProgress() {
	for showProgressFlag {
		clearScreen()
		for n, t := range activeTorrents {
			pcount := t.NumPieces()
			completed := 0
			s := t.Stats()

			for _, x := range t.PieceStateRuns() {
				if x.Complete {
					completed += x.Length
				}
			}

			per := float32(completed) / float32(pcount)
			fmt.Printf("Status: %s\nProgress: %.2f\nSeeders: %d\nActive Peers: %d\nPeers: %d\n", n, per, s.ConnectedSeeders, s.ActivePeers, s.TotalPeers)
		}
		time.Sleep(time.Millisecond * 500)
	}
}
