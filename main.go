package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/httpclient"
	"github.com/iqrahadian/brick-web-scrapper-assesment/scrapper"
	"golang.org/x/time/rate"
)

func main() {

	//init http client with rate limiter
	rateLimit := rate.NewLimiter(rate.Every(10*time.Second), 50) // 50 request every 10 seconds
	httpClient := httpclient.NewClient(rateLimit)

	var (
		wg         sync.WaitGroup
		jobChannel = make(chan model.Product)
	)

	startProductScrapper(&wg, &httpClient, jobChannel)

	// init database

	tokpedUrl := []string{
		"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=1",
		"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=2",
	}

	productList := []model.Product{}

	for _, url := range tokpedUrl {
		products, err := scrapper.ScrapProductListPage(&httpClient, url)
		if err != nil {
			fmt.Errorf(
				fmt.Sprintf("Failed to retrieve Product list page : ", url),
				err,
			)
		}
		productList = append(productList, products...)
	}

	for _, product := range productList {
		jobChannel <- product
	}

}

func startProductScrapper(wg *sync.WaitGroup, client *httpclient.RLHTTPClient, productChannel <-chan model.Product) {
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
