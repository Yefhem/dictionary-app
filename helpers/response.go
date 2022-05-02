package helpers

// Response Model -------->
type Response struct {
	Status    bool        `json:"status"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Errors    interface{} `json:"errors"`
	PageModel interface{} `json:"page_model"`
	Data      interface{} `json:"data"`
}

type PageModel struct {
	Page     int64   `json:"page"`
	LastPage float64 `json:"last_page"`
	Total    int64   `json:"total"`
}

// BuildSuccessResponse -------->
func BuildSuccessResponse(status bool, code int, message string, total int64, page int64, lastPage float64, data interface{}) Response {
	pageModel := PageModel{
		Page:     page,
		LastPage: lastPage,
		Total:    total,
	}
	res := Response{
		Status:    status,
		Code:      code,
		Message:   message,
		Errors:    nil,
		PageModel: pageModel,
		Data:      data,
	}
	return res
}

// BuildErrorResponse -------->
func BuildErrorResponse(message string, code int, err string) Response {
	res := Response{
		Status:    false,
		Code:      code,
		Message:   message,
		Errors:    err,
		PageModel: nil,
		Data:      nil,
	}
	return res
}
