package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type R struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
}

type HandleResult = func(c *gin.Context) *R

// Ok 返回http状态码为200数据
func Ok(data interface{}) *R {
	return &R{StatusCode: http.StatusOK, Data: data, Message: "OK"}
}

// OkMessage 返回http状态码为200数据与message信息
func OkMessage(data interface{}, message string) *R {
	return &R{StatusCode: http.StatusOK, Data: data, Message: message}
}

// InternalError 当程序执行过程中发生内部错误进行返回
func InternalError(err error) *R {
	return Error(http.StatusInternalServerError, err)
}

// AuthenticationError 细化错误，表示身份验证错误
func AuthenticationError(err error) *R {
	return Error(http.StatusUnauthorized, err)
}

// Error return error
func Error(statusCode int, err error) *R {
	return &R{StatusCode: statusCode, Message: err.Error()}
}

// Handle 处理请求，根据handler参数获取返回结果集
func Handle(handler HandleResult) gin.HandlerFunc {
	if handler == nil {
		panic("handler is empty")
	}
	return func(c *gin.Context) {
		r := handler(c)
		if r != nil {
			Write(c, r)
		}
	}
}

func Write(c *gin.Context, r *R) {
	c.JSON(r.StatusCode, r)
}
