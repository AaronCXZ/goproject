package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//Server的接口实例，定义一个server的服务器模块
type Server struct {
	//  服务器名称
	Name string
	//  服务器IP版本
	IPVersion string
	//  服务器绑定的IP地址
	IP string
	//  服务器监听的端口
	Port int
	//	当前Server的消息管理模块，用来绑定MsgID的对应的处理业务API的关系
	MsgHandler ziface.IMsgHandle
	//	该Server的链接管理器
	ConnManager ziface.IConnManager
	//  该server创建链接之后自动调用的Hook函数
	OnConnStart func(conn ziface.IConnection)
	//  该server销毁链接之前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

////调用当前客户端连接所绑定的handleAPI
//func CallBackTOClinet(conn *net.TCPConn, data []byte, cnt int) error {
//	//	回显业务
//	fmt.Println("[Conn Handle] CallBackTOClient....")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Printf("Write back buf err %s\n", err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d is starting....\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("[Start] Server Listenner at IP: %s Port: %d, is starting\n", s.IP, s.Port)
	go func() {
		//开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()
		//获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println(err)
			return
		}
		//监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Start Zinx server succ, %s succ, Listenning... ", s.Name)
		var cid uint32
		cid = 0
		//阻塞的等待客户端连接，处理客户端连接业务
		for {
			//如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				continue
			}
			//设置最大链接个数的判断，如果超过最大连接数，那么关闭此新链接
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应连接数过大
				fmt.Printf("Too many conns MaxConn = %d\n", utils.GlobalObject.MaxConn)
				_ = conn.Close()
				continue
			}
			//将处理新链接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//  启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	//将一些服务器的资源、状态或一些已经开辟的连接信息进行停止会回收
	fmt.Printf("Stop Zinx server name %s\n", s.Name)
	s.ConnManager.ClearConn()
}
func (s *Server) Server() {
	//启动server的服务功能
	s.Start()
	//TODO 做一些启动服务器额外的业务
	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

//注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart...")
		s.OnConnStart(conn)
	}
}

//调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop...")
		s.OnConnStop(conn)
	}
}

//初始化server模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}
