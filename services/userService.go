package services

import (
	"log"

	"github.com/Yefhem/mongo/dictionary/models"
	"github.com/Yefhem/mongo/dictionary/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Insert(user models.User) error
	FindByKeyValue(key string, value interface{}) (models.User, error)
	CheckEmailPass(email, pass interface{}) (models.User, error)
	VerifyPassword(hashedPassword, password string) bool
}

type userService struct {
	userRepository repository.UserRepository
}

func MewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (a *userService) Insert(user models.User) error {
	if err := a.userRepository.Insert(user); err != nil {
		return err
	}
	return nil
}

func (a *userService) FindByKeyValue(key string, value interface{}) (models.User, error) {
	user, err := a.userRepository.FindByKeyValue(key, value)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *userService) CheckEmailPass(email, pass interface{}) (models.User, error) {
	user, err := a.userRepository.CheckEmailPass(email, pass)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *userService) VerifyPassword(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Println(err)
		return false
	}
	return true
}
