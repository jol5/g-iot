package exsample
import (
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"time"
)


type echoServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}




func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {

	//gbk := mahonia.NewDecoder("gbk")
	bytes, _ := simplifiedchinese.GBK.NewDecoder().Bytes(frame)
	st := string(bytes)
	println(st)

	data := append([]byte{}, frame...)
	//go func() {
	//	time.Sleep(time.Second)
	//	c.AsyncWrite(data)
	//}()

	_ = es.pool.Submit(func() {
		time.Sleep(1 * time.Second)
		c.AsyncWrite(data)
	})

	return


	//out = frame
	//return
}

func MyMain() {
	println("服务开始启动")
	echo := new(echoServer)
	log.Fatal(gnet.Serve(echo, "tcp://:8000", gnet.WithMulticore(true)))
}