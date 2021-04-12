package myrpc

import (
	"errors"
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	addr  string                   // 服务器监听地址
	funcs map[string]reflect.Value //函数名到函数的映射表
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

// 注册服务
func (s *Server) Register(name string, f interface{}) error {
	if _, exist := s.funcs[name]; exist {
		return errors.New("register fail, service already exist, name:" + name)
	}
	s.funcs[name] = reflect.ValueOf(f)
	return nil
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		//连接会话处理
		go func() {
			session := NewSession(conn)
			//接受消息
			message, err := session.RecvMsg()
			if err != nil {
				fmt.Println("session RecvMsg fail,", err)
				return
			}
			//找到函数名对应的函数
			f, ok := s.funcs[message.FuncName]
			if !ok {
				errMsg := fmt.Sprintf("rpc server func:%s not exist", message.FuncName)
				err = session.SendMsg(RpcMessage{FuncName: message.FuncName, Args: message.Args, ErrMsg: errMsg})
				if err != nil {
					fmt.Println("session SendMsg fail, err:", err)
					return
				}
			}

			inArgs := make([]reflect.Value, len(message.Args))
			for i := range message.Args {
				inArgs[i] = reflect.ValueOf(message.Args[i])
			}
			// 调用函数
			out := f.Call(inArgs)
			outArgs := make([]interface{}, len(out))
			for i := 0; i < len(out); i++ {
				outArgs[i] = out[i].Interface()
			}

			// 返回结果
			err = session.SendMsg(RpcMessage{FuncName: message.FuncName, Args: outArgs, ErrMsg: ""})
			if err != nil {
				fmt.Println("session SendMsg fail, err:", err)
				return
			}
		}()
	}
}
