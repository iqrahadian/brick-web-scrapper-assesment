package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/filerepo"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/headlessclient"
	"github.com/iqrahadian/brick-web-scrapper-assesment/repo/sqlite"
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

// func mainTest() {

// 	var wg sync.WaitGroup

// 	headlessClient := headlessclient.NewClient(1, 5*time.Second)

// 	productList, _ := scrapper.ExtractProductList(page1)

// 	i := 0
// 	for _, product := range productList {

// 		if strings.Contains(product.ProductUrl, "https://ta.tokopedia.com") {
// 			continue
// 		}
// 		if i > 2 {
// 			break
// 		}

// 		go func(hc *headlessclient.RLHeadlessClient, product model.Product, wg *sync.WaitGroup) {
// 			wg.Add(1)
// 			defer wg.Done()
// 			detailedProduct, err := scrapper.ScrapProductDetailPage(hc, product)

// 			if err != nil {
// 				fmt.Println("Name : ", detailedProduct.Name)
// 				fmt.Println("ERROR : ", err)
// 			} else {
// 				fmt.Println("Name : ", detailedProduct.Name)
// 				// fmt.Println("Description : ", detailedProduct.Description)
// 				fmt.Println("Rate : ", detailedProduct.Rating)
// 			}
// 		}(&headlessClient, product, &wg)

// 		i += 1
// 		fmt.Println("===========================")
// 	}

// 	wg.Wait()

// }

func main() {

	//init http client with rate limiter
	headlessClient := headlessclient.NewClient(10, 5*time.Second)   // to scrap product detail
	headlessClient2 := headlessclient.NewClient(10, 10*time.Second) // to scrap product list

	var (
		wg            sync.WaitGroup
		jobChannel    = make(chan model.Product)
		reportChannel = make(chan model.Product)
	)

	startProductScrapper(&wg, &headlessClient, jobChannel)

	tokpedUrl := []string{
		"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=1",
		"https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=2",
	}

	productList := []model.Product{}

	productLen := 0
	for _, url := range tokpedUrl {
		if productLen > 10 {
			break
		}

		products, err := scrapper.ScrapProductListPage(&headlessClient2, url)
		if err != nil {
			fmt.Errorf(
				fmt.Sprintf("Failed to retrieve Product list page : ", url),
				err,
			)
		}

		productList = append(productList, products...)
		productLen += 1
	}

	// submit all retrieved product to complete necessary data in product detail page
	// in this case we are looking for product rating & description
	for _, product := range productList {
		jobChannel <- product
	}

	// this line meant to waiting for all data completion finished
	// in real case, usually it run on async, not waiting for all the data fo be finished
	close(jobChannel)
	wg.Wait()
	close(reportChannel)

	// retrieving all finalized productdata
	finalProductList := []model.Product{}
	for product := range reportChannel {
		finalProductList = append(finalProductList, product)
	}

	// storing to csv file, this will store to /tmp folder in the project
	// inside the implementation we will chunk the row based on config submitted
	csvRepo, _ := filerepo.NewFileRepo(filerepo.FileRepoCsvType, 2)

	err := csvRepo.Save(finalProductList)
	if err != nil {
		fmt.Println("ERR", err)
	}

	// storing to sqlite database
	// TODO : move to postgres database
	// p.s. do not have enough time for seting up postgres, potentially run on docker, etc. (it's sunday)
	db := sqlite.NewClient()

	// gorm already have batch system implemented, so we will use it
	// in more advance situation, this batching process might be written manually-
	// to make sure throughput to database & incoming data are processed savely
	result := db.CreateInBatches(finalProductList, 50)
	if result.Error != nil {
		fmt.Println("FAILED to store data to persitence ", err)
	}

	// var users []model.Product
	// // Get all records
	// result = db.Find(&users)

	// _, err = result.Rows()
	// if err != nil {
	// 	fmt.Println("Failed to select from database", err)
	// }

	// for _, prod := range users {
	// 	fmt.Println(prod.ID, prod.Name)
	// }

	fmt.Println("DONE")

}

func startProductScrapper(
	wg *sync.WaitGroup,
	client *headlessclient.RLHeadlessClient,
	productChannel <-chan model.Product,

) {
	for worker := 1; worker <= 5; worker++ {
		wg.Add(1)
		go func() {

			defer wg.Done()

			for product := range productChannel {

				// product url that start with https://ta.tokopedia.com tends to have security issue
				// opening this webpage in some browser resulting in security concern
				// for now, we skip product with this url
				if !strings.Contains(product.ProductUrl, "https://ta.tokopedia.com") {
					product, err := scrapper.ScrapProductDetailPage(client, product)
					if err != nil {
						fmt.Println("=================================")
						fmt.Println("Name : ", product.Name)
						fmt.Println("Failed To Scrap Product Detail : ", err)
						fmt.Println("=================================")
					}

				}

				scrapper.AppendTo(product)
			}

		}()
	}

}
