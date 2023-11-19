package scrapper

import (
	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
)

var productMap = make(map[string]model.Product)

func ExtractProductList(stringHtml string) []model.Product

func ExtractProductPage(product model.Product, stringHtml string) (model.Product, error)
