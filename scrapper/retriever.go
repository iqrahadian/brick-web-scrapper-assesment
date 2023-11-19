package scrapper

import (
	"github.com/go-rod/rod"
)

func retrieveHtmlPage(url string) (string, error) {

	page := rod.New().MustConnect().MustPage(url).MustWaitStable()
	defer page.MustClose()
	page.Mouse.MustScroll(0, 300)

	stringHtml, err := page.HTML()
	return stringHtml, err

}

func RetrieveProductListPage(url string) (string, error) {
	return retrieveHtmlPage(url)
}

func RetrieveProductDetailPage(url string) (string, error) {
	return retrieveHtmlPage(url)
}
