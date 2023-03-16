package msg

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
	"sync"
)

var msgCodeAndMsgDescMap = make(map[int16]protoreflect.MessageDescriptor)

var msgNameAndMsgCodeMap = make(map[string]int16)
var locker = sync.Mutex{}

func getMsgDescByMsgCode(msgCode int16) (protoreflect.MessageDescriptor, error) {
	if msgCode < 0 {
		return nil, errors.New("消息代号无效")
	}
	return msgCodeAndMsgDescMap[msgCode], nil
}

func getMsgCodeByMsgName(msgName string) (int16, error) {
	if len(msgName) <= 0 {
		return -1, errors.New("消息名称为空")
	}
	lowerName := strings.ToLower(
		strings.Replace(msgName, "_", "", -1),
	)
	return msgNameAndMsgCodeMap[lowerName], nil

}
func init() {
	initMap()
}

// InitMap 初始化
func initMap() {

	locker.Lock()
	defer locker.Unlock()
	//名称-->d代号
	if len(msgNameAndMsgCodeMap) <= 0 {
		for k, v := range MsgCode_value {
			msgName := strings.ToLower(strings.Replace(k, "_", "", -1))
			msgNameAndMsgCodeMap[msgName] = int16(v)
		}

	}

	//代号-->描述
	if len(msgCodeAndMsgDescMap) <= 0 {
		msgDescList := File_GameMsgProtocol_proto.Messages()
		for i := 0; i < msgDescList.Len(); i++ {
			msgDesc := msgDescList.Get(i)

			name := strings.ToLower(strings.Replace(string(msgDesc.Name()), "_", "", -1))
			msgCode := msgNameAndMsgCodeMap[name]
			msgCodeAndMsgDescMap[msgCode] = msgDesc

		}
	}
}
