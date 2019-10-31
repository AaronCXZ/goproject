package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
存储一切有关Zinx框架的全局参数，供其它模块使用
一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //  当前Zinx全局的Server对象
	Host      string         //  当前服务监听的IP
	TcpPort   int            //  当前服务监听的端口
	Name      string         //  当前服务的名称
	/*
		Zinx
	*/
	Version          string //当前Zinx的版本
	MaxConn          int    //当前服务允许的最大连接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 //worker的数量
	MaxWorkerTaskLen uint32 //允许用户开辟多少个worker
}

/*
	定义一个全局的对外GlobalObj
*/

var GlobalObject *GlobalObj

//从zinx.json加载自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("src/myDome/ZinxV0.10/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//	将json文件数据解析到GlobalObj中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//初始化方法，初始化当前GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxSereverApp",
		Version:          "V0.4",
		TcpPort:          8080,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}
