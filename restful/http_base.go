package restful

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ERROR      = 1000
	SUCCESS    = 2000
	UN_SIGN_IN = 3000
	FAILURE    = 5000
)

type RestFulResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Error(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, RestFulResponse{Code: ERROR, Message: "服务器繁忙,请稍后再试!", Data: nil})
}

func Success(ctx *gin.Context, message *string, data any) {
	if message == nil {
		ctx.JSON(http.StatusOK, RestFulResponse{Code: SUCCESS, Message: "操作成功!", Data: data})
	} else {
		ctx.JSON(http.StatusOK, RestFulResponse{Code: SUCCESS, Message: *message, Data: data})
	}
}

func Failure(ctx *gin.Context, message *string, data any) {
	if message == nil {
		ctx.JSON(http.StatusOK, RestFulResponse{Code: FAILURE, Message: "操作失败!", Data: data})
	} else {
		ctx.JSON(http.StatusOK, RestFulResponse{Code: FAILURE, Message: *message, Data: data})
	}
}

func UnSignIn(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, RestFulResponse{Code: UN_SIGN_IN, Message: "请先进行登录!", Data: data})
}
