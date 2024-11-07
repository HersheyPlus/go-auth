package dto


import (
	"github.com/gin-gonic/gin"
)

const (
    StatusSuccess = "success"
    StatusError   = "error"
    StatusFail    = "fail"
)

type StandardResponse struct {
    Status  string      `json:"status"`
	StatusCode int 		`json:"status_code"`
    Message string      `json:"message,omitempty"`
    Error   *ErrorDetail `json:"error,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

type ErrorDetail struct {
    Message string `json:"message"`
    Note *string `json:"note,omitempty"`
}

type ResponseBuilder struct {
    ctx *gin.Context
}

func NewResponse(c *gin.Context) *ResponseBuilder {
    return &ResponseBuilder{ctx: c}
}

func (rb *ResponseBuilder) Success(httpStatus int ,data interface{}, message string) {
    response := StandardResponse{
        Status:  StatusSuccess,
		StatusCode: httpStatus,
        Message: message,
        Data:    data,
    }
    rb.ctx.JSON(httpStatus, response)
}

func (rb *ResponseBuilder) Error(httpStatus int, message string) {
    response := StandardResponse{
        Status: StatusError,
		StatusCode: httpStatus,
        Error: &ErrorDetail{
            Message: message,
        },
    }
    rb.ctx.JSON(httpStatus, response)
}

func (rb *ResponseBuilder) ValidationError(httpStatus int, message string, note string) {
    response := StandardResponse{
        Status: StatusFail,
        StatusCode: httpStatus,
        Error: &ErrorDetail{
            Message: message,
            Note: &note,
        },
    }
    rb.ctx.JSON(httpStatus, response)
}