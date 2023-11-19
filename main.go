package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/headlessclient"
	"github.com/iqrahadian/brick-web-scrapper-assesment/scrapper"
)

func mainSingleProduct() {

	// init http client with rate limiter
	// rateLimit := rate.NewLimiter(rate.Every(10*time.Second), 50) // 50 request every 10 seconds
	// httpClient := httpclient.NewClient(rateLimit)

	// tokpedUrl := "https://www.tokopedia.com/duniagadgetku/xiaomi-redmi-10-4-64-gb-6-128-gb-garansi-resmi-6-128-promo?extParam=cmp%3D1%26ivf%3Dtrue"

	// stringHtml, err := scrapper.RetrieveProductDetailPage(tokpedUrl)
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Failed to retrieve product page", err))
	// }

	// fmt.Println("HTML : ", stringHtml)

	// product := model.Product{ProductUrl: tokpedUrl}
	// product, err := scrapper.ScrapProductDetailPage(&httpClient, product)
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Failed to retrieve product page", err))
	// }

}

// func mainProductList() {
func main() {

	var wg sync.WaitGroup

	headlessClient := headlessclient.NewClient(1, 5*time.Second)

	productList, _ := scrapper.ExtractProductList(page1)

	i := 0
	for _, product := range productList {

		if strings.Contains(product.ProductUrl, "https://ta.tokopedia.com") {
			continue
		}
		if i > 2 {
			break
		}

		go func(hc *headlessclient.RLHeadlessClient, product model.Product, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			detailedProduct, err := scrapper.ScrapProductDetailPage(hc, product)

			if err != nil {
				fmt.Println("Name : ", detailedProduct.Name)
				fmt.Println("ERROR : ", err)
			} else {
				fmt.Println("Name : ", detailedProduct.Name)
				// fmt.Println("Description : ", detailedProduct.Description)
				fmt.Println("Rate : ", detailedProduct.Rating)
			}
		}(&headlessClient, product, &wg)

		i += 1
		fmt.Println("===========================")
	}

	// fmt.Println(produ)

	wg.Wait()

}

func main2() {

	//init http client with rate limiter
	headlessClient := headlessclient.NewClient(10, 10*time.Second)

	var (
		wg         sync.WaitGroup
		jobChannel = make(chan model.Product)
	)

	startProductScrapper(&wg, &headlessClient, jobChannel)

	// init database

	tokpedUrl := []string{
		"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=1",
		"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=2",
	}

	productList := []model.Product{}

	productLen := 0
	for _, url := range tokpedUrl {
		if productLen > 100 {
			break
		}

		products, err := scrapper.ScrapProductListPage(&headlessClient, url)
		if err != nil {
			fmt.Errorf(
				fmt.Sprintf("Failed to retrieve Product list page : ", url),
				err,
			)
		}

		productList = append(productList, products...)
		productLen += 1
	}

	for _, product := range productList {
		fmt.Println(product.Name, product.Price)
		// jobChannel <- product
	}

	fmt.Println("PRODUCT LEN :", len(productList))

}

func startProductScrapper(wg *sync.WaitGroup, client *headlessclient.RLHeadlessClient, productChannel <-chan model.Product) {
	for w := 1; w <= 5; w++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, productChannel <-chan model.Product) {

			defer wg.Done()
			for {
				select {
				case product := <-productChannel:
					scrapper.ScrapProductDetailPage(client, product)
				}
			}

		}(wg, productChannel)
	}

}
