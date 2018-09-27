package providers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const convertKB = 1024
const convertMB = 1024 * convertKB
const convertGB = 1024 * convertMB

// GetGoQueryDocument gets a goquery document provided the url
func GetGoQueryDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	return doc, err
}

// ConvertSize converts kb,mb,gb to bytes
func ConvertSize(ssize string) (uint, error) {
	ssize = strings.ToLower(strings.TrimSpace(ssize))
	parts := strings.Split(ssize, " ")

	fsize, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	switch parts[1] {
	case "kb":
		fsize *= convertKB
		break
	case "mb":
		fsize *= convertMB
		break
	case "gb":
		fsize *= convertGB
		break
	}

	return uint(fsize), nil
}

// GetTorrentXt returns the xt parameter from a magnet link
func GetTorrentXt(magnet string) (string, error) {
	rp, err := regexp.Compile("xt=([^&]+)")
	if err != nil {
		return "", err
	}
	t := rp.FindStringSubmatch(magnet)
	if len(t) == 0 {
		return "", nil
	}
	return t[1], nil
}
