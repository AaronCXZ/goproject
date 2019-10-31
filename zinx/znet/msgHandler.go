package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
  消息管理实现
*/
type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//	负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//  worker的数量
	WorkerPoolSize uint32
}

//初始化方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//  调度执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//	从request中找到MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("api msgID = %d is NOT FOUND! Need Register!", request.GetMsgID())
	}
	//根据MsgID调度对应的Router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//	为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前msg绑定的API处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//	添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Printf("Add api MsgID = %d succ...", msgID)
}

//启动一个worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//  根据WorkerPoolSize开启Worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//  一个worker被启动
		//当前的worker对应的channel的消息队列
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前的worker，阻塞等待消息从channel传递过来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Printf("Worker ID = %d is started\n", workerID)
	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，取出一个客户端的request，执行当前的request绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给taskQueue处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的worker
	//根据客户端监控的ConnID进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Printf("Add ConnID = %d, request MsgId = %d\n", request.GetConnection().GetConnID(), workerID)
	//将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
