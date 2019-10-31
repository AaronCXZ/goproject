package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//拆包和封包的单元测试
func TestDataPack(t *testing.T) {
	/*
	   模拟的服务
	*/
	//创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println(err)
	}
	//创建一个go承载，负责处理客户端业务
	go func() {
		//从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
			}
			go func(conn net.Conn) {
				//处理客户端的请求,拆包
				dp := NewDataPack()
				for {
					//第一次从conn读，读出包的head信息
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println(err)
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println(err)
					}
					if msgHead.GetMsgLen() > 0 {
						//	第二次从conn读。读出包的msg信息
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(msg.Id)
						fmt.Println(msg.DataLen)
						fmt.Println(string(msg.Data))
					}
				}
			}(conn)
		}
	}()
	/*
	   模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println(err)
	}
	//	创建一个封包对象dp
	dp := NewDataPack()
	//模拟粘包过程，封装两个msg一起发送
	//封装msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'Z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println(err)
	}

	//封装msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'H', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println(err)
	}
	//将两个包粘在一起
	sendData := append(sendData1, sendData2...)
	//一次性发送两个包给服务端
	conn.Write(sendData)
	//	客户端阻塞
	select {}
}
