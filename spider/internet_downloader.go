package spider

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/getumen/lucy"
)

// DownloadInternet is a sample spider that follows all link in the html.
func DownloadInternet(response *lucy.Response) ([]*lucy.Request, error) {
	requests := make([]*lucy.Request, 0)
	if strings.Contains(response.ContentType(), "text/html") {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Text()))
		if err != nil {
			log.Printf("fail to parse html in %s", response.Request.URL)
		}
		title := doc.Find("title").Text()
		fmt.Printf(title)
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			if exists {
				request, err := response.FollowRequest(link)
				if err == nil {
					requests = append(requests, request)
				}
			}
		})
	}
	return requests, nil
}
