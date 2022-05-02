package services

import (
	"net/http"

	"github.com/Yefhem/mongo/dictionary/repository"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type SessionService interface {
	InitSession(c echo.Context) (*sessions.Session, error)
	SetUser(c echo.Context, username, email, pass, Picture string) error
	CheckUser(c echo.Context) bool
	RemoveUser(c echo.Context) error
}

type sessionService struct {
	userRepository repository.UserRepository
	store          *sessions.CookieStore
}

func NewSessionService(userRepo repository.UserRepository, store *sessions.CookieStore) SessionService {
	return &sessionService{
		userRepository: userRepo,
		store:          store,
	}
}

func (s *sessionService) InitSession(c echo.Context) (*sessions.Session, error) {
	session, err := s.store.Get(c.Request(), "dictionary")
	if err != nil {
		return nil, err
	} else {
		if session.IsNew {
			session.Options.MaxAge = 86400 * 3 // 3 day
			session.Options.HttpOnly = true
			session.Options.SameSite = http.SameSiteStrictMode
		}
		return session, err
	}
}

func (s *sessionService) SetUser(c echo.Context, username, email, pass, picture string) error {
	session, err := s.InitSession(c)
	if err != nil {
		return err
	}

	session.Values["username"] = username
	session.Values["email"] = email
	session.Values["password"] = pass
	session.Values["picture"] = picture

	return sessions.Save(c.Request(), c.Response())
}

func (s *sessionService) CheckUser(c echo.Context) bool {
	session, err := s.InitSession(c)
	if err != nil {
		return false
	}

	email := session.Values["email"]
	pass := session.Values["password"]

	user, err := s.userRepository.CheckEmailPass(email, pass)
	if err != nil {
		return false
	}

	if user.Email == email && user.Password == pass {
		return true
	} else {
		return false
	}
}

func (s *sessionService) RemoveUser(c echo.Context) error {
	session, err := s.InitSession(c)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1

	return session.Save(c.Request(), c.Response())
}
