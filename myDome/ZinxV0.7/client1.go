package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

//模拟客户端
func main() {
	fmt.Println("Client1 start...")
	time.Sleep(1 * time.Second)
	//连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		//发送封包的msg
		data := "Zinx v0.6 client test message..."
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(&znet.Message{
			Id:      1,
			DataLen: uint32(len(data)),
			Data:    []byte(data),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println(err)
			return
		}
		//读取服务端发送的消息
		//第一次从conn读，读出包的head信息
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println(err)
			break
		}
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println(err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//	第二次从conn读。读出包的msg信息
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println(err)
				break
			}
			fmt.Printf("recv from server: ID = %d, Len = %d, Msg: %s\n", msg.Id, msg.DataLen, string(msg.Data))
		}
		//	阻塞
		time.Sleep(1 * time.Second)
	}

}
