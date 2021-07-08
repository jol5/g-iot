package routers

import (
	"g-iot/pkg/mygin"
	"github.com/gin-gonic/gin"
	"net/http"

)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/index", http.Dir('/'))
	r.GET("/get", func(ctx *gin.Context) {
		//ctx.Writer.WriteString("ok")

		r := mygin.Gin{C: ctx}

		r.Response(http.StatusOK,mygin.SUCCESS,"ok")
	})

	return r
}
