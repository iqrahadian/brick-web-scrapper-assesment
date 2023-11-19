package filerepo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
)

type FileRepoType string

const (
	FileRepoCsvType FileRepoType = "csv"
)

func NewFileRepo(repoType FileRepoType, rowLimit int) (FileRepo, error) {

	switch repoType {
	case FileRepoCsvType:
		return CsvRepo{rowLimit}, nil
	default:
		return nil, errors.New("Unsupported file type")
	}
}

type FileRepo interface {
	Save([]model.Product) error
}

type CsvRepo struct {
	rowLimit int
}

func (c CsvRepo) Save(products []model.Product) error {

	folderTarget := "./tmp"
	baseFileName := fmt.Sprintf("tokopedia-product-%d", time.Now().Nanosecond())

	maxChunk := math.Ceil(float64(len(products)) / float64(c.rowLimit))
	for i := float64(0); i < maxChunk; i++ {

		filePath := fmt.Sprintf("%s/%s-%f.csv", folderTarget, baseFileName, i)
		startRow := int(i * float64(c.rowLimit))
		endRow := c.rowLimit * (int(i) + 1)

		if startRow+c.rowLimit >= len(products) {
			endRow = len(products) - 1
		}

		fmt.Println("start : ", startRow, "end : ", endRow)
		err := c.executeWrite(filePath, products[startRow:endRow])
		if err != nil {
			return err
		}

	}

	return nil

}

func (c CsvRepo) executeWrite(filePath string, products []model.Product) error {

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	defer writer.Flush()

	headers := []string{"Title", "Rate", "Price", "Merchant", "ProductUrl", "ImageUrl", "Description"}

	writer.Write(headers)
	for _, product := range products {

		strData := []string{
			product.Name,
			strconv.FormatFloat(float64(product.Rating), 'f', -1, 64),
			strconv.FormatFloat(product.Price, 'f', -2, 64),
			product.Merchant,
			product.ProductUrl,
			product.ImageUrl,
			product.Description,
		}

		writer.Write(strData)
	}

	return nil

}
