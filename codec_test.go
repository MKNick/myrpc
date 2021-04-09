package myrpc

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	m1 := RpcMessage{
		FuncName: "Login",
		Args: []interface{}{
			"nickxiong",
			"123456",
			30}}

	buf, err := Encode(m1)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := Decode(buf)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v\n", m2)
}
