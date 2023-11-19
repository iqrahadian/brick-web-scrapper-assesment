package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/filerepo"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/headlessclient"
	"github.com/iqrahadian/brick-web-scrapper-assesment/scrapper"
)

func mainCsv() {
	// func main() {

	csvRepo, _ := filerepo.NewFileRepo(filerepo.FileRepoCsvType, 2)

	err := csvRepo.Save([]model.Product{})
	if err != nil {
		fmt.Println("ERR", err)
	}

}

func mainTest() {

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

	wg.Wait()

}

func main() {

	//init http client with rate limiter
	headlessClient := headlessclient.NewClient(10, 5*time.Second) // to scrap product detail
	// headlessClient2 := headlessclient.NewClient(10, 10*time.Second) // to scrap product list

	var (
		wg         sync.WaitGroup
		jobChannel = make(chan model.Product)
	)

	startProductScrapper(&wg, &headlessClient, jobChannel)

	// init database

	// tokpedUrl := []string{
	// 	"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=1",
	// 	"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=2",
	// }

	// productList := []model.Product{}

	// productLen := 0
	// for _, url := range tokpedUrl {
	// 	if productLen > 10 {
	// 		break
	// 	}

	// 	products, err := scrapper.ScrapProductListPage(&headlessClient2, url)
	// 	if err != nil {
	// 		fmt.Errorf(
	// 			fmt.Sprintf("Failed to retrieve Product list page : ", url),
	// 			err,
	// 		)
	// 	}

	// 	productList = append(productList, products...)
	// 	productLen += 1
	// }
	productList, _ := scrapper.ExtractProductList(page1)

	for idx, product := range productList {
		// fmt.Println(product.Name, product.Price)
		if idx > 5 {
			break
		}
		jobChannel <- product
	}

	close(jobChannel)
	wg.Wait()

	csvRepo, _ := filerepo.NewFileRepo(filerepo.FileRepoCsvType, 2)

	err := csvRepo.Save(scrapper.ProductArr)
	if err != nil {
		fmt.Println("ERR", err)
	}

	fmt.Println("GOT HERE DONE")

}

func startProductScrapper(wg *sync.WaitGroup, client *headlessclient.RLHeadlessClient, productChannel <-chan model.Product) {
	for worker := 1; worker <= 5; worker++ {
		wg.Add(1)
		go func() {

			defer wg.Done()

			for product := range productChannel {
				if !strings.Contains(product.ProductUrl, "https://ta.tokopedia.com") {
					product, err := scrapper.ScrapProductDetailPage(client, product)
					if err != nil {
						fmt.Println("Name : ", product.Name)
						fmt.Println("ERROR : ", err)
					} else {
						fmt.Println("Name : ", product.Name)
						fmt.Println("Price : ", product.Price)
						fmt.Println("Rate : ", product.Rating)
					}

				} else {
					fmt.Println("Name : ", product.Name)
					fmt.Println("Skipping because ta.tokopedia")

				}

				scrapper.AppendTo(product)
				fmt.Println("================================")
			}

		}()
	}

}
