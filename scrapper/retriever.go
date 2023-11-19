package scrapper

import (
	"time"

	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/headlessclient"
)

func RetrieveProductListPage(hc *headlessclient.RLHeadlessClient, url string) (string, error) {
	return hc.RetrieveHtml(url, 2*time.Second)
}

func RetrieveProductDetailPage(hc *headlessclient.RLHeadlessClient, url string) (string, error) {
	return hc.RetrieveHtml(url, 1*time.Second)
}
