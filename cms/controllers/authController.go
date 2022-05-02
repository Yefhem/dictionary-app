package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Yefhem/mongo/dictionary/helpers"
	"github.com/Yefhem/mongo/dictionary/services"
	"github.com/labstack/echo/v4"
)

type AuthController interface {
	LoginIndex(c echo.Context) error

	Login(c echo.Context) error
	Logout(c echo.Context) error
}

type authController struct {
	userService    services.UserService
	sessionService services.SessionService
	alertService   services.AlertService
}

func NewAuthController(userServ services.UserService, sessionServ services.SessionService, alertServ services.AlertService) AuthController {
	return &authController{
		userService:    userServ,
		sessionService: sessionServ,
		alertService:   alertServ,
	}
}

// ----------------> Pages
// --------> Login Index Page
func (cont *authController) LoginIndex(c echo.Context) error {
	results, err := helpers.Include("userops/login")
	if err != nil {
		log.Fatal(err)
	}
	view, err := template.ParseFiles(results...)
	if err != nil {
		fmt.Println(err)
	}

	data := make(map[string]interface{})

	data["Alert"] = cont.alertService.GetAlert(c)

	return view.ExecuteTemplate(c.Response().Writer, "index", data)
}

// ----------------> Methods
// --------> Login
func (cont *authController) Login(c echo.Context) error {
	email := c.FormValue("email")
	pass := c.FormValue("password")

	result, ans := helpers.LoginValidate(email, pass)
	if !result {
		cont.alertService.SetAlert(c, ans)
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	user, err := cont.userService.FindByKeyValue("email", email)
	if err != nil {
		log.Println(err)
		cont.alertService.SetAlert(c, "Kullanıcı adı veya Parola Hatalı!")
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}
	storedUserPass := user.Password

	if cont.userService.VerifyPassword(storedUserPass, pass) {
		// login
		cont.sessionService.SetUser(c, user.Username, email, storedUserPass, user.Picture)
		cont.alertService.SetAlert(c, "Hoşgeldiniz")
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	} else {
		// Denied
		cont.alertService.SetAlert(c, "Yanlış Kullanıcı Adı veya Şifre!")
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

}

// --------> Logout
func (cont *authController) Logout(c echo.Context) error {
	cont.sessionService.RemoveUser(c)
	cont.alertService.SetAlert(c, "Güle Güle")
	return c.Redirect(http.StatusSeeOther, "/admin/login")
}
