package ziface

/*
	定义一个解决TCp粘包问题的封包和拆包的模块
    直接面向TCP连接中的数据流
*/

type IDataPack interface {
	GetHeadLen() uint32              //	获取包的头的长度方法
	Pack(IMessage) ([]byte, error)   //   封包方法
	UnPack([]byte) (IMessage, error) //   拆包方法
}
