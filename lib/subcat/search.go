package subcat

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

type SearchResult struct {
	Title     string
	Path      string
	Size      string
	Downloads int
}

func Search(keyword string) []SearchResult {
	u := buildURL(baseURL, "index.php", map[string]string{"search": keyword})
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	res := []SearchResult{}

	doc.Find("table.sub-table tbody tr").Each(func(i int, s *goquery.Selection) {
		it := SearchResult{}

		it.Path, _ = url.QueryUnescape(s.Find("td a").AttrOr("href", ""))

		// path maybe urlencoded

		for i, n := range s.Find("td").EachIter() {
			text := strings.TrimSpace(n.Text())
			switch i {
			case 0:
				it.Title = text
			case 2:
				it.Size = strings.TrimPrefix(text, "SIZE")
			case 3:
				parts := strings.Fields(text)
				if len(parts) > 0 {
					downloads, err := strconv.Atoi(parts[0])
					if err != nil {
						log.Warnf("failed to parse downloads from %q: %v", text, err)
					}
					it.Downloads = downloads
				}
			}
		}
		res = append(res, it)
	})

	return res
}
