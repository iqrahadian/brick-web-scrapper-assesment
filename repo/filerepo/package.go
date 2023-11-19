package filerepo

import "errors"

type FileRepoType string

const (
	FileRepoCsvType FileRepoType = "csv"
)

func NewFileRepo(repoType FileRepoType) (FileRepo, error) {

	switch repoType {
	case FileRepoCsvType:
		return CsvRepo{}, nil
	default:
		return nil, errors.New("Unsupported file type")
	}
}

type FileRepo interface {
	Save(path string) error
}

type CsvRepo struct{}

func (c CsvRepo) Save(path string) error {
	return nil
}
