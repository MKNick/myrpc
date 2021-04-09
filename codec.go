package myrpc

import (
	"bytes"
	"encoding/gob"
)

// 函数名，函数参数进行编解码
type RpcMessage struct {
	FuncName string
	Args     []interface{}
	ErrMsg   string
}

func Encode(m RpcMessage) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (RpcMessage, error) {
	var m RpcMessage
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&m)
	if err != nil {
		return m, err
	}
	return m, nil
}
