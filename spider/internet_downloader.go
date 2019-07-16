package spider

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/getumen/arachne"
)

// DownloadInternet is a sample spider that follows all link in the html.
func DownloadInternet(response *arachne.Response) ([]*arachne.Request, error) {
	requests := make([]*arachne.Request, 0)
	if strings.Contains(response.ContentType(), "text/html") {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Text()))
		if err != nil {
			log.Printf("fail to parse html in %s", response.Request.URL)
		}
		title := doc.Find("title").Text()
		fmt.Println(title)
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
