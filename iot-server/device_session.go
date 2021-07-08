package iot_server

import (
	cmap "github.com/orcaman/concurrent-map"
	"github.com/panjf2000/gnet"
)

type IotDeviceSessionMg struct {
	SessionMap cmap.ConcurrentMap
	unAuthSessionMap cmap.ConcurrentMap
}

func (m IotDeviceSessionMg) IsAuth(c gnet.Conn) bool {

	key := c.RemoteAddr().String()

	val, b := m.unAuthSessionMap.Get(key)
	if b {
		isAuth := val.(*IotDeviceSession).IsAuth
		return isAuth
	}

	return m.SessionMap.Has(key)
}

func (m IotDeviceSessionMg) RegistryAuthedConn(c gnet.Conn) {
	key := c.RemoteAddr().String()
	val, b := m.unAuthSessionMap.Get(key)
	if b {
		m.SessionMap.Set(key, val)
		m.unAuthSessionMap.Remove(key)
	}
}

func (m IotDeviceSessionMg) RemoveErrAuthConn(c gnet.Conn)  {
	unAuthSessionMap.Remove(c.RemoteAddr().String())
	_ = c.Close()
}


type IotDeviceSession struct {
	Conn *gnet.Conn
	IsAuth bool

	createTimestamp int64
	DeviceName string
	ProductKey string
}

func (s *IotDeviceSession) send(msg string) error {

	err := (*s.Conn).AsyncWrite([]byte(msg))
	if err != nil {
		return err
	}

	return nil
}
