package znet

import (
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"

	"github.com/pkg/errors"
)

type Connection struct {
	//  当前conn属于哪个server
	TCPServer ziface.IServer
	//	当前链接的socket TCP套结字
	Conn *net.TCPConn
	//  链接的ID
	ConnID uint32
	//  当前链接状态
	isClosed bool
	////  当前链接所绑定的处理业务方法API
	//handleAPI ziface.HandleFunc
	//	告知当前链接已经退出的channel,由Reader告知Write退出的信号
	ExitChan chan bool
	//  无缓冲的通道，写Goroutine之间的消息通道
	msgChan chan []byte
	//	当前Server的消息管理模块，用来绑定MsgID的对应的处理业务API的关系
	MsgHandle ziface.IMsgHandle

	//	连接属性集合
	property map[string]interface{}
	//  保护链接属性的锁
	propertyLock sync.RWMutex
}

//初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, MsgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer: server,
		Conn:      conn,
		ConnID:    connID,
		//handleAPI: callBackApi,
		MsgHandle: MsgHandle,
		isClosed:  false,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1), //
		property:  make(map[string]interface{}),
	}
	// 将conn加入到connManager中
	c.TCPServer.GetConnManager().Add(c)

	return c
}

//链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader] Goroutine is running....")
	defer fmt.Printf("connID = %d , Reader is exit, remote addr is %s\n", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()
	for {
		//	读取客户端的数据到buffer中
		/*
			buf := make([]byte, utils.GlobalObject.MaxPackageSize)
			_, err := c.Conn.Read(buf)
			if err != nil {
				fmt.Printf("recv buf err: %s\n", err)
				continue
			}
		*/
		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端消息的head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read msg head err: ", err)
			break
		}
		//拆包，得到MsgID和MsgDataLen
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err: ", err)
			break
		}

		//根据MsgDataLen再次读取Data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err: ", err)
				break
			}
		}
		msg.SetData(data)
		//  得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//使用工作池处理请求
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//  从路由中找到注册绑定的Conn对应的router调用
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

//写消息的Goroutine，用于给客户端发送消息
func (c *Connection) StartWrite() {
	fmt.Println("[Write] Goroutine is running... ")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Write exit!")
	//	不断的阻塞的等待channel的消息，写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//	有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			//	代表Reader已经退出，此时需要退出Write
			return
		}
	}
}

//提供一个SendMsg方法，将需要发送的数据先进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when read msg...")
	}
	//将data进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		return errors.New("Pack msg err")
	}
	//将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) Start() {
	fmt.Printf("Conn start().... ConnID = %d", c.ConnID)
	//	启动从当前链接读取数据的业务
	go c.StartReader()
	//  启动向当前链接写数据的业务
	go c.StartWrite()
	//	按照开发者传递进来的创建链接之后需要调用的处理业务，执行对应的hook函数
	c.TCPServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop()... ConnID = ", c.ConnID)
	//如果当前链接已关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//	按照开发者传递进来的销毁链接之前需要调用的处理业务，执行对应的hook函数
	c.TCPServer.CallOnConnStop(c)
	//关闭socket链接
	_ = c.Conn.Close()
	//关闭write
	c.ExitChan <- true
	//将当前链接从connManager删除
	c.TCPServer.GetConnManager().Remove(c)
	//关闭channel
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//添加属性
	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	//读取属性
	if value, ok := c.property[key]; ok {
		return value, nil
		fmt.Println("=====>", value)
	}
	return nil, errors.New("no property found")
}

//移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
