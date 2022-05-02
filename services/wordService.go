package services

import (
	"fmt"
	"time"

	"github.com/Yefhem/mongo/dictionary/models"
	"github.com/Yefhem/mongo/dictionary/repository"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WordService interface {
	Insert(wordDTO dto.WordDTO) error
	GetWord(id string) (models.Word, error)
	GetAllWords() ([]models.Word, error)
	UpdateWord(id string, wordDTO dto.WordDTO) error
	DeleteWord(id string) error
}

type wordService struct {
	wordRepository repository.WordRepository
}

func NewWordService(wordRepo repository.WordRepository) WordService {
	return &wordService{
		wordRepository: wordRepo,
	}
}

// ----------------> Methods
// --------> Add Book
func (w *wordService) Insert(wordDTO dto.WordDTO) error {

	var word models.Word

	word.ID = primitive.NewObjectID()
	word.CreatedAt = time.Now()
	word.Slug = slug.Make(word.ID.String()[10:34])
	word.English = wordDTO.English
	word.Turkish = wordDTO.Turkish
	word.Abbreviation = wordDTO.Abbreviation
	word.Description = wordDTO.Description

	if err := w.wordRepository.Insert(word); err != nil {
		return err
	}
	return nil
}

// --------> Get Word
func (w *wordService) GetWord(id string) (models.Word, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Word{}, err
	}

	word, err := w.wordRepository.Get(objID)
	fmt.Println(word)
	if err != nil {
		return word, err
	}
	return word, nil
}

// --------> Get All Words
func (w *wordService) GetAllWords() ([]models.Word, error) {
	words, err := w.wordRepository.GetAll()
	if err != nil {
		return words, err
	}
	return words, nil
}

// --------> Update Word
func (w *wordService) UpdateWord(id string, wordDTO dto.WordDTO) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := w.GetWord(id)
	if err != nil {
		return err
	}

	UpdateWord := models.Word{
		ID:           result.ID,
		CreatedAt:    result.CreatedAt,
		Slug:         result.Slug,
		English:      wordDTO.English,
		Turkish:      wordDTO.Turkish,
		Abbreviation: wordDTO.Abbreviation,
		Description:  wordDTO.Description,
	}

	if UpdateWord == result {
		return apperrors.ErrSameObj
	}

	if err := w.wordRepository.Update(objID, UpdateWord); err != nil {
		return err
	}
	return nil
}

// --------> Delete Word
func (w *wordService) DeleteWord(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := w.wordRepository.Delete(objID); err != nil {
		return err
	}
	return nil
}
