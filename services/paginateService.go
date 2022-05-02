package services

import (
	"math"
	"strconv"

	"github.com/Yefhem/mongo/dictionary/models"
	"github.com/Yefhem/mongo/dictionary/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginateService interface {
	PaginateFunc(value, page, sort string) ([]models.Word, int64, int64, float64, error)
}

type paginateService struct {
	wordRepository repository.WordRepository
}

func NewPaginateService(wordRepo repository.WordRepository) PaginateService {
	return &paginateService{
		wordRepository: wordRepo,
	}
}

func (p *paginateService) PaginateFunc(value, page, sort string) ([]models.Word, int64, int64, float64, error) {
	var words []models.Word
	var DefaultLimit int64 = 2
	var lastpage float64

	convPage, _ := strconv.ParseInt(page, 10, 64)
	if convPage < 1 {
		convPage = 1
	}

	findOptions := options.Find()

	findOptions.SetLimit(DefaultLimit)
	findOptions.SetSkip((convPage - 1) * DefaultLimit)

	if sort == "asc" {
		findOptions.SetSort(bson.D{{"created_at", 1}})
	} else if sort == "desc" {
		findOptions.SetSort(bson.D{{"created_at", -1}})
	} else {
		findOptions.SetSort(bson.D{{"created_at", -1}})
	}

	words, total, err := p.wordRepository.Search(value, findOptions)
	if err != nil {
		return words, total, convPage, lastpage, err
	}

	if (total % DefaultLimit) == 0 {
		lastpage = math.Ceil(float64(total / DefaultLimit))
	} else {
		lastpage = (math.Ceil(float64(total/DefaultLimit)) + 1)
	}

	return words, total, convPage, lastpage, nil
}
