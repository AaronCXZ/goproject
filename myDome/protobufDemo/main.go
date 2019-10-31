package main

import (
	"fmt"
	"myDome/protobufDemo/pd"

	"github.com/golang/protobuf/proto"
)

func main() {
	//	定义一个Person结构对象
	person := &pd.Person{
		Name:   "Chen",
		Age:    22,
		Emails: []string{"a", "b", "c", "d", "e"},
		Phones: []*pd.PhoneNumber{
			&pd.PhoneNumber{
				Number: "17601008872",
				Type:   0,
			},
			&pd.PhoneNumber{
				Number: "18328491696",
				Type:   1,
			},
		},
	}

	//将person对象进行序列化，得到一个二进制数据data
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err: ", err)
	}
	//	data就是要进行网络传输的漱口，对端需要按照Message Person格式进行反序列化
	newPerson := &pd.Person{}
	err = proto.Unmarshal(data, newPerson)
	if err != nil {
		fmt.Println("unmarshal err: ", err)
	}
	fmt.Println(newPerson)
}
