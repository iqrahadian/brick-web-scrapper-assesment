package scrapper

import (
	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/httpclient"
)

func ScrapProductListPage(client *httpclient.RLHTTPClient, url string) ([]model.Product, error) {
	stringHtml, err := RetrieveProductListPage(client, url)
	if err != nil {
		return []model.Product{}, err
	}

	productList := ExtractProductList(stringHtml)
	if err != nil {
		return []model.Product{}, err
	}

	return productList, nil
}

func ScrapProductDetailPage(client *httpclient.RLHTTPClient, product model.Product) {

	stringHtml, err := RetrieveProductDetailPage(client, product.ProductUrl)
	if err != nil {
		return
	}

	product, err = ExtractProductPage(product, stringHtml)
	if err != nil {
		return
	}

	// put product to saving queue
}
