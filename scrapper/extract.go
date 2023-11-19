package scrapper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
)

func ExtractProductList(stringHtml string) ([]model.Product, error) {
	productList := []model.Product{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(stringHtml))
	if err != nil {
		return []model.Product{}, err
	}

	doc.Find("div.css-bk6tzz").Each(func(i int, s *goquery.Selection) {

		product, err := parseSingleProductList(s)
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to parse single product list", err))
		} else {
			productList = append(productList, product)
		}
	})

	return productList, nil
}

func parseSingleProductList(doc *goquery.Selection) (model.Product, error) {

	// For each item found, get the title
	title := doc.Find("div.css-11s9vse span.css-20kt3o").Text()
	if title == "" {
		return model.Product{}, errors.New("Failed to parse Product Title")
	}

	merChantName := doc.Find("a > div.css-16vw0vn > div.css-11s9vse > div.css-tpww51 > div.css-vbihp9 > span").Last().Text()
	if merChantName == "" {
		return model.Product{}, errors.New("Failed to parse Merchant Name")
	}

	price, err := parsePrice(doc.Find("div.css-pp6b3e span.css-o5uqvq").Text())
	if err != nil {
		return model.Product{}, err
	}

	productUrl, exist := doc.Find("a.css-54k5sq").Attr("href")
	if !exist {
		return model.Product{}, errors.New("Failed to parse Product Url")
	}

	imageUrl, exist := doc.Find("img").Attr("src")
	if !exist {
		return model.Product{}, errors.New("Failed to parse Image Url")
	}

	return model.Product{
		Name:       title,
		Price:      price,
		ProductUrl: productUrl,
		Merchant:   merChantName,
		ImageUrl:   imageUrl,
	}, nil

}

func parsePrice(price string) (float64, error) {

	price = strings.Replace(strings.ReplaceAll(price, ".", ""), "Rp", "", -1)
	if price == "" {
		return 0, errors.New("Failed to retrieve Price")
	}

	s, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0, err
	}

	return s, nil
}

func ExtractProductPage(product model.Product, stringHtml string) (model.Product, error) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(stringHtml))
	if err != nil {
		return product, err
	}

	description := doc.Find("span.css-11oczh8 > span.css-17zm3l > div").Text()
	description = strings.ReplaceAll(description, "	", " ")
	description = strings.ReplaceAll(description, "\t", " ")

	product.Description = description

	// we can assume all product will have at least 1 star rating, because we take top 100 product
	rate, err := parseRate(doc.Find("div.css-8atqhb > div.css-856ghu > div.css-1m5sihj > div.css-1fogemr > div.css-jmbq56 > div.css-bczdt6 > div.items > p.css-vni7t6-unf-heading > span > span.main").Text())
	if err != nil {
		return product, err
	}

	product.Rating = rate

	return product, nil
}

func parseRate(r string) (float64, error) {

	rate, err := strconv.ParseFloat(r, 64)
	if err != nil {
		return 0, err
	}

	return rate, nil

}
