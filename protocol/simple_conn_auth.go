package protocol

import (
	"fmt"
	"g-iot/pkg/log"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"strings"
)


func DeviceConnAuth(body string, c gnet.Conn) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error(fmt.Sprintf("捕获到的错误：body=%s,err=%s \n",body, r))
		}
	}()

	split := strings.Split(body, "|")

	if split == nil ||len(split)  != 2 {
		_ = c.AsyncWrite([]byte("auth.empty_identify"))
	}
	
	userName := split[0]
	pwd := split[1]

	if "admin" != userName || "admin" != pwd {
		_ = c.AsyncWrite([]byte("auth.fail"))
		return false
	}

	addr := c.RemoteAddr().String()
	_ = c.AsyncWrite([]byte("auth.pass"))
	log.Logger.Info("设备接入认证成功=>",zap.Stringp("deviceAddr",&addr))

	return true
}