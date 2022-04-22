package rpc

import "fmt"

/**
input 传入参数
output 回传参数 引用传递不需要返回值
*/
func (*Server) Ping(input string, output *string) error {

	fmt.Println(input)
	*output = "收到了"
	return nil
}