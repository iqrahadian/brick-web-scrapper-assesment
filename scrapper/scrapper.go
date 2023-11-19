package scrapper

import (
	"fmt"
	"sync"
	"time"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/headlessclient"
)

var Mutex = &sync.RWMutex{}
var ProductArr = []model.Product{}

func AppendTo(product model.Product) {
	Mutex.Lock()
	ProductArr = append(ProductArr, product)
	Mutex.Unlock()
}

func ScrapProductListPage(client *headlessclient.RLHeadlessClient, url string) ([]model.Product, error) {
	stringHtml, err := RetrieveProductListPage(client, url)
	if err != nil {
		return []model.Product{}, err
	}

	productList, err := ExtractProductList(stringHtml)
	if err != nil {
		return []model.Product{}, err
	}

	return productList, nil
}

func ScrapProductDetailPage(client *headlessclient.RLHeadlessClient, product model.Product) (model.Product, error) {

	startTime := time.Now()

	stringHtml, err := RetrieveProductDetailPage(client, product.ProductUrl)
	if err != nil {
		return product, err
	}

	product, err = ExtractProductPage(product, stringHtml)
	if err != nil {
		return product, err
	}

	endTime := time.Now()
	timeDiff := endTime.Sub(startTime)

	fmt.Println("ScrapProductDetailPage finished in : ", timeDiff)

	return product, err
}
