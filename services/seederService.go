package services

import (
	"log"
	"os"

	"github.com/Yefhem/mongo/dictionary/models"
	"github.com/Yefhem/mongo/dictionary/repository"
)

type SeederService interface {
	UploadUser()
}

type seederService struct {
	userRepository repository.UserRepository
}

func NewSeederService(userRepo repository.UserRepository) SeederService {
	return &seederService{
		userRepository: userRepo,
	}
}

func (s *seederService) UploadUser() {
	if err := s.userRepository.DeleteAllUser(); err != nil {
		log.Printf("Kullanıcılar silinirken bir hata oluştu: %v", err)
		return
	}

	username := os.Getenv("ADMIN_USERNAME")
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASS")

	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
		Picture:  "uploads/cat.jpg",
	}

	if err := s.userRepository.Insert(user); err != nil {
		log.Printf("[-] %v", err)
		panic(err)
	}
	log.Printf("[+] %v", username)
}
