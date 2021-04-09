package myrpc

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSession_Read(t *testing.T) {
	addr := "127.0.0.1:9000"
	testData := "hello shopee"
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatal(err)
			}
			srv := NewSession(conn)
			err = srv.write([]byte(testData))
			if err != nil {
				t.Fatal(err)
			}
			break
		}
	}()

	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		session := NewSession(conn)
		data, err := session.read()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(data))

		if string(data) != testData {
			t.Error("data not equal")
		}
	}()

	wg.Wait()

}
