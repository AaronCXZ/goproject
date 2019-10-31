package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/*
	链接管理模块实现
*/

type ConnManager struct {
	//管理的链接集合
	connections map[uint32]ziface.IConnection
	//保护链接集合的读写锁
	connLock sync.RWMutex
}

//初始化
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//	保护共享资源，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//	将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Printf("Connection add to connManager succ, connID = %d， num = %d\n", conn.GetConnID(), connMgr.Len())
}

//删除
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//	保护共享资源，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除链接
	delete(connMgr.connections, conn.GetConnID())
	fmt.Printf("Connection remove to connManager succ, connID = %d, num = %d\n", conn.GetConnID(), connMgr.Len())
}

//查找
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//	保护共享资源，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("Connection not found!")
	}
}

//获取长度
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清空
func (connMgr *ConnManager) ClearConn() {
	//	保护共享资源，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//	删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear all connection to connManager succ.")
}
