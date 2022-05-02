package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Yefhem/mongo/dictionary/dto"
	"github.com/Yefhem/mongo/dictionary/helpers"
	apperrors "github.com/Yefhem/mongo/dictionary/models/app-errors"
	"github.com/Yefhem/mongo/dictionary/services"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type WordController interface {
	// ---------------->
	DashboardIndex(c echo.Context) error
	DashboardNewWord(c echo.Context) error
	DashboardEditWord(c echo.Context) error
	// ---------------->
	AddWord(c echo.Context) error
	UpdateWord(c echo.Context) error
	DeleteWord(c echo.Context) error
}

type wordController struct {
	wordService    services.WordService
	userService    services.UserService
	sessionService services.SessionService
	alertService   services.AlertService
}

func NewWordController(wordServ services.WordService, userServ services.UserService, sessionServ services.SessionService, alertServ services.AlertService) WordController {
	return &wordController{
		wordService:    wordServ,
		userService:    userServ,
		sessionService: sessionServ,
		alertService:   alertServ,
	}
}

// ----------------> Pages
// --------> Index Page
func (cont *wordController) DashboardIndex(c echo.Context) error {
	if !cont.sessionService.CheckUser(c) {
		cont.alertService.SetAlert(c, "Lütfen Giriş Yapınız!")
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	results, err := helpers.Include("dashboard/dashboard_main")
	if err != nil {
		log.Fatal(err)
	}
	view, err := template.ParseFiles(results...)
	if err != nil {
		fmt.Println(err)
	}

	words, err := cont.wordService.GetAllWords()
	if err != mongo.ErrNoDocuments && err != nil {
		return err
	}

	session, _ := cont.sessionService.InitSession(c)

	data := make(map[string]interface{})

	data["Words"] = words
	data["Alert"] = cont.alertService.GetAlert(c)
	data["SessionUsername"] = session.Values["username"]
	data["SessionPicture"] = session.Values["picture"]

	return view.ExecuteTemplate(c.Response().Writer, "index", data)
}

// --------> New Word Page
func (cont *wordController) DashboardNewWord(c echo.Context) error {
	if !cont.sessionService.CheckUser(c) {
		cont.alertService.SetAlert(c, "Lütfen Giriş Yapınız!")
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	results, err := helpers.Include("dashboard/dashboard_add")
	if err != nil {
		log.Fatal(err)
	}
	view, err := template.ParseFiles(results...)
	if err != nil {
		fmt.Println(err)
	}

	session, _ := cont.sessionService.InitSession(c)
	usernameSession := session.Values["username"]

	data := make(map[string]interface{})
	data["Alert"] = cont.alertService.GetAlert(c)
	data["Session"] = usernameSession

	return view.ExecuteTemplate(c.Response().Writer, "index", data)
}

// --------> Edit Word Page
func (cont *wordController) DashboardEditWord(c echo.Context) error {
	if !cont.sessionService.CheckUser(c) {
		cont.alertService.SetAlert(c, "Lütfen Giriş Yapınız!")
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}
	results, err := helpers.Include("dashboard/dashboard_edit")
	if err != nil {
		log.Fatal(err)
	}
	view, err := template.ParseFiles(results...)
	if err != nil {
		fmt.Println(err)
	}

	result, err := cont.wordService.GetWord(c.Param("id"))
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	data["Word"] = result
	data["Alert"] = cont.alertService.GetAlert(c)

	return view.ExecuteTemplate(c.Response().Writer, "index", data)
}

// ----------------> Methods
// --------> Add Word
func (cont *wordController) AddWord(c echo.Context) error {
	if !cont.sessionService.CheckUser(c) {
		return nil
	}
	word := dto.WordDTO{
		English:      c.FormValue("word-name-en"),
		Turkish:      c.FormValue("word-name-tr"),
		Abbreviation: c.FormValue("word-abb"),
		Description:  c.FormValue("word-desc"),
	}

	result, str := helpers.NewWordValidate(word.English, word.Turkish, word.Abbreviation, word.Description)
	if !result {
		cont.alertService.SetAlert(c, str)
		return c.Redirect(http.StatusSeeOther, "/admin/new-word")
	}

	if err := cont.wordService.Insert(word); err != nil {
		log.Println(err)
		cont.alertService.SetAlert(c, "İşlem Başarısız!")
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}
	cont.alertService.SetAlert(c, "Kayıt Ekleme İşlemi Başarılı.")
	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")

}

// --------> Update Word
func (cont *wordController) UpdateWord(c echo.Context) error {
	if !cont.sessionService.CheckUser(c) {
		return nil
	}
	id := c.Param("id")

	word := dto.WordDTO{
		English:      c.FormValue("word-name-en"),
		Turkish:      c.FormValue("word-name-tr"),
		Abbreviation: c.FormValue("word-abb"),
		Description:  c.FormValue("word-desc"),
	}

	result, str := helpers.NewWordValidate(word.English, word.Turkish, word.Abbreviation, word.Description)
	if !result {
		cont.alertService.SetAlert(c, str)
		return c.Redirect(http.StatusSeeOther, "/admin/word-edit/"+id)
	}

	if err := cont.wordService.UpdateWord(id, word); err != nil {
		if err == apperrors.ErrSameObj {
			cont.alertService.SetAlert(c, err.Error())
			return c.Redirect(http.StatusSeeOther, "/admin/word-edit/"+id)
		}
		cont.alertService.SetAlert(c, "Güncelleme İşlemi Başarısız.")
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	} else {
		cont.alertService.SetAlert(c, "Güncelleme İşlemi Başarılı.")
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}
}

// --------> Delete Word
func (cont *wordController) DeleteWord(c echo.Context) error {
	if !cont.sessionService.CheckUser(c) {
		return nil
	}
	id := c.Param("id")

	if err := cont.wordService.DeleteWord(id); err != nil {
		log.Println(err)
		cont.alertService.SetAlert(c, "İşlem Başarısız!")
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}
	cont.alertService.SetAlert(c, "Kayıt Silme İşlemi Başarılı.")
	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}
