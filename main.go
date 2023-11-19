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

func main() {

	//init http client with rate limiter
	headlessClient := headlessclient.NewClient(10, 1*time.Second)   // to scrap product detail
	headlessClient2 := headlessclient.NewClient(10, 10*time.Second) // to scrap product list

	var (
		wg             sync.WaitGroup
		productChannel = make(chan model.Product)
		reportChannel  = make(chan model.Product, 150)
	)

	// init multithread to scrape product detail
	startProductScrapper(&wg, &headlessClient, productChannel, reportChannel)

	// retrieve and submit top product to complete necessary data in product detail page
	// in this case we are looking for product rating & description
	baseProductList, err := scrapper.ScrapTopProductList(&headlessClient2)
	if err != nil {
		panic(err)
	}
	fmt.Println("BASE PRODUCT COUNT : ", len(baseProductList))

	fmt.Println("Start submitting product to channel, to complete the data")
	for idx, baseProduct := range baseProductList {
		if idx > 10 {
			break
		}
		productChannel <- baseProduct
	}

	// this line meant to waiting for all data completion finished
	// in real case, usually it run on async, not waiting for all the data fo be finished
	close(productChannel)
	wg.Wait()
	close(reportChannel)

	fmt.Println("Finished completing product data")

	// retrieving all finalized productdata
	finalProductList := []model.Product{}
	for product := range reportChannel {
		finalProductList = append(finalProductList, product)
	}

	fmt.Println("Start saving product data to csv file under ./tmp")
	// storing to csv file, this will store to /tmp folder in the project
	// inside the implementation we will chunk the row based on config submitted
	csvRepo, _ := filerepo.NewFileRepo(filerepo.FileRepoCsvType, 50)

	err = csvRepo.Save(finalProductList)
	if err != nil {
		fmt.Println("Failed to save CSV File", err)
	}

	fmt.Println("Start saving product data to sqlite")
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

	// Re-retrieving data that has been inserted, validating the process
	var products []model.Product
	result = db.Find(&products)

	_, err = result.Rows()
	if err != nil {
		fmt.Println("Failed to select from database", err)
	} else {
		fmt.Println("============================================")
		fmt.Println("Start Re-retrieving product data from sqlite")
		fmt.Println("============================================")
		for _, prod := range products {
			fmt.Println("ID :", prod.ID, "Name :", prod.Name, "Rating :", prod.Rating)
		}
	}

	fmt.Println("============================================")
	fmt.Println("DONE, csv file can be found under ./tmp written using tab as Comma")

}

func startProductScrapper(
	wg *sync.WaitGroup,
	client *headlessclient.RLHeadlessClient,
	productChannel <-chan model.Product,
	reportChan chan model.Product,
) {
	for worker := 1; worker <= 5; worker++ {
		wg.Add(1)
		go func(w *sync.WaitGroup, pc <-chan model.Product, rc chan model.Product) {

			defer w.Done()

			for product := range pc {

				fmt.Println("Start retrieving Detail Product Page for : ", product.Name)

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
					rc <- product

				} else {
					rc <- product
				}

			}

		}(wg, productChannel, reportChan)
	}

}
