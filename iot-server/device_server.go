package iot_server

import (
	"fmt"
	log "g-iot/pkg/log"
	"g-iot/protocol"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"go.uber.org/zap"
	"time"

	"github.com/orcaman/concurrent-map"
)
var logger = func() *zap.Logger{
	return log.GetLogger()
}

var DeviceSessionMg IotDeviceSessionMg
var unAuthSessionMap cmap.ConcurrentMap

type iotDeviceServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}


func init(){
	DeviceSessionMg = IotDeviceSessionMg{
		unAuthSessionMap: cmap.New(),
		SessionMap: cmap.New(),
	}

	unAuthSessionMap =  DeviceSessionMg.unAuthSessionMap
}

func InitAuthConnChecker() {
	log.Logger.Info("启动授权连接校验器")
	for {
		curSec := time.Now().Unix()
		for iter := range unAuthSessionMap.IterBuffered() {

			val := iter.Val
			session := val.(*IotDeviceSession)

			if session.IsAuth {
				unAuthSessionMap.Remove(iter.Key)
				return
			}

			if curSec >= session.createTimestamp {
				err := (*session.Conn).Close()
				if err != nil {
					log.Logger.Error(fmt.Sprintf("关闭设备连接失败 =>%s",err))
					return
				}
				log.Logger.Info(fmt.Sprintf("关闭未认证的超时设备连接 => %s",iter.Key))
				unAuthSessionMap.Remove(iter.Key)
			}
		}

	}
}


func (es *iotDeviceServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Logger.Info( fmt.Sprintf("设备服务器 监听地址为 %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop))
	return
}

func (es *iotDeviceServer) OnClosed(c gnet.Conn, err error) (action gnet.Action)  {

	log.Logger.Info("连接关闭")
	return
}

func (es iotDeviceServer)  OnOpened(c gnet.Conn) (out []byte, action gnet.Action)  {

	s := IotDeviceSession{
		Conn: &c,
		createTimestamp: time.Now().Unix() + 5,
	}
	key := c.RemoteAddr().String()
	unAuthSessionMap.Set(key,&s)

	logger().Info(fmt.Sprintf("与[%s]建立了通讯通道",c.RemoteAddr().String()))
	return
}


func (es *iotDeviceServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {

	if !DeviceSessionMg.IsAuth(c) {
		go func() {
			auth := protocol.DeviceConnAuth(string(frame), c)
			if auth {
				DeviceSessionMg.RegistryAuthedConn(c)
			}else {
				DeviceSessionMg.RemoveErrAuthConn(c)
			}
		}()

	}else {
		// 业务处理
		logger().Info(fmt.Sprintf("[%s] 接收到设备数据 => %s", c.RemoteAddr().String() ,string(frame)))
	}


	return
}


func Start() {
	log.Logger.Info("开始运行 设备TCP服务")

	go InitAuthConnChecker()

	port := 9610

	server := new(iotDeviceServer)
	err := gnet.Serve(server, fmt.Sprintf("tcp://:%d", port))
	if err != nil {
		log.Logger.Error(fmt.Sprintf("启动 IOT Server 失败！！！\n=>%s",err.Error()))
		return
	}
	log.Logger.Fatal(fmt.Sprintf("启动 IOT Server 成功！=> [ tcp://:%d ]",port))

}