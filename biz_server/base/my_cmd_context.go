package base

import "google.golang.org/protobuf/reflect/protoreflect"

type MyCmdContext interface {
	// BindUserId 绑定用户
	BindUserId(val int64)
	// GetUserId 获取用户
	GetUserId() int64
	// GetClientIpAddr 获取客户端ip
	GetClientIpAddr() string
	// Write 写入消息
	Write(msgObj protoreflect.ProtoMessage)
	// SendError 错误消息
	SendError(errorCode int, errorInfo string)
	// Disconnect 断开连接
	Disconnect()
}
