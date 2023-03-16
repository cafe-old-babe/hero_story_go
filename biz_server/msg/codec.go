package msg

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Decode 解码
func Decode(msgData []byte, msgCode int16) (*dynamicpb.Message, error) {
	if nil == msgData || len(msgData) <= 0 {
		return nil, errors.New("消息数据为空")
	}
	//拿到消息描述符
	code, err := getMsgDescByMsgCode(msgCode)
	if err != nil {
		return nil, err
	}
	message := dynamicpb.NewMessage(code)
	if err := proto.Unmarshal(msgData, message); err != nil {
		return nil, err
	}

	return message, nil
}

// Encode 编码
func Encode(msgObj protoreflect.ProtoMessage) ([]byte, error) {
	if msgObj == nil {
		return nil, errors.New("消息对象为空")
	}
	//消息代号
	msgCode, err := getMsgCodeByMsgName(string(msgObj.ProtoReflect().Descriptor().Name()))

	if err != nil {
		return nil, err
	}

	msgCodeByteArray := make([]byte, 2)

	binary.BigEndian.PutUint16(msgCodeByteArray, uint16(msgCode))
	msgSizeByteArray := make([]byte, 2)
	binary.BigEndian.PutUint16(msgSizeByteArray, 0)

	bytes, err := proto.Marshal(msgObj)
	if err != nil {
		return nil, err
	}
	completeMsg := append(msgSizeByteArray, msgCodeByteArray...)
	completeMsg = append(completeMsg, bytes...)
	return completeMsg, nil
}
