package scrapper

import (
	"fmt"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/headlessclient"
)

func ScrapTopProductList(client *headlessclient.RLHeadlessClient) ([]model.Product, error) {

	// not using pagination, 2 page of tokopedia product list will retrieve 150 product
	// query string ob=5 means sort by review, which I assume it's the top product
	tokpedUrl := []string{
		"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=1",
		// "https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=2",
	}

	maxProduct := 100
	productList := []model.Product{}
	for _, url := range tokpedUrl {
		if len(productList) > maxProduct {
			break
		}
		fmt.Println("Start Retrieving Top product Page : ", url)

		products, err := ScrapProductListPage(client, url)
		if err != nil {
			fmt.Errorf(
				fmt.Sprintf("Failed to retrieve Product list page : ", url),
				err,
			)

			continue
		}

		productRange := (maxProduct - len(productList)) - 1
		productList = append(productList, products[0:productRange]...)
		fmt.Println("Done")
	}

	return productList, nil
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

	// startTime := time.Now()

	stringHtml, err := RetrieveProductDetailPage(client, product.ProductUrl)
	if err != nil {
		return product, err
	}

	product, err = ExtractProductPage(product, stringHtml)
	if err != nil {
		return product, err
	}

	// timeDiff := time.Now().Sub(startTime)
	// fmt.Println("ScrapProductDetailPage finished in : ", timeDiff)

	return product, nil
}
