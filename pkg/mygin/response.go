package mygin

import (
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}

func (g *Gin) ResponseSuc(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  e.GetMsg(SUCCESS),
		Data: data,
	})
	return
}

func (g *Gin) ResponseError(errCode int, data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}