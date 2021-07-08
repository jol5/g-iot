package main

import (
	"fmt"
	iot_server "g-iot/iot-server"
	log "g-iot/pkg/log"
	"g-iot/routers"
	"net/http"
	"time"
)
import "github.com/gin-gonic/gin"

func main() {
	log.InitLogger("./all.log", "debug", 50, 10, 0, false, false, true, true)
	logger := log.Logger

	logger.Info("初始化日志完成")

	go iot_server.Start()

	gin.SetMode(gin.DebugMode)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           "localhost:8000",
		Handler:        routers.InitRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logger.Info(fmt.Sprintf("服务器启动，地址为 %s",server.Addr))

	err := server.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Sprintf("监听服务器失败=>%s",err.Error()))
		return
	}


}
