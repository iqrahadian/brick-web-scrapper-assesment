package scrapper

import (
	"github.com/go-rod/rod"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/httpclient"
)

func retrieveHtmlPage(url string) (string, error) {

	page := rod.New().MustConnect().MustPage(url).MustWaitStable()
	defer page.MustClose()

	stringHtml, err := page.HTML()
	return stringHtml, err

}

func RetrieveProductListPage(client *httpclient.RLHTTPClient, url string) (string, error)

func RetrieveProductDetailPage(client *httpclient.RLHTTPClient, url string) (string, error)
