package apicontrollers

import (
	"log"
	"net/http"

	"github.com/Yefhem/mongo/dictionary/helpers"
	"github.com/Yefhem/mongo/dictionary/services"
	"github.com/labstack/echo/v4"
)

type SearchController interface {
	SearchingWord(c echo.Context) error
}

type searchController struct {
	paginateService services.PaginateService
}

func NewSearchController(pagiServ services.PaginateService) SearchController {
	return &searchController{
		paginateService: pagiServ,
	}
}

func (cont *searchController) SearchingWord(c echo.Context) error {
	value := c.QueryParam("q")
	pageNumber := c.QueryParam("page_number")
	orderBy := c.QueryParam("order_by")

	if value == "" {
		response := helpers.BuildErrorResponse("value boş olamaz!", http.StatusBadRequest, "BAD_REQUEST")
		return c.JSON(http.StatusBadRequest, response)
	}

	words, total, page, lastPage, err := cont.paginateService.PaginateFunc(value, pageNumber, orderBy)
	if err != nil {
		log.Println(err)
		response := helpers.BuildErrorResponse(err.Error(), http.StatusBadRequest, "BAD_REQUEST")
		return c.JSON(http.StatusBadRequest, response)
	}
	response := helpers.BuildSuccessResponse(true, http.StatusOK, "OK", total, page, lastPage, words)
	return c.JSON(http.StatusOK, response)
}
