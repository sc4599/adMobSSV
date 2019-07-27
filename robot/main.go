package main

import (
	"fmt"
	"os"
	"golang.org/x/net/websocket"
	"time"
	"awesomeProject/robot/protocol"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
)

func main() {
	origin := "http://127.0.0.1:7777/"
	addr := "ws://127.0.0.1:7777/ws"

	//origin := "http://192.168.0.210:7777/"
	//addr := "ws://192.168.0.210:7777/ws"
	ws, err := websocket.Dial(addr, "", origin)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := ws.Read(buffer)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			msg := &protocol.MsgHead{}
			proto.Unmarshal(buffer[:n], msg)

			switch msg.CmdId {
			case protocol.Cmd_Login:
				login := &protocol.LoginResp{}
				if err := proto.Unmarshal(msg.Body, login); err != nil {
					fmt.Println(err)
				} else {
					if login.Code != protocol.RespCode_SUCCESS {
						fmt.Printf("response error: code[%d], error[%v]\n", login.Code, login.Error)
					} else {
						fmt.Printf("response success: %v\n", login.Error)
					}
				}
			case protocol.Cmd_Create:
				create := &protocol.CreateResp{}
				if err := proto.Unmarshal(msg.Body, create); err != nil {
					fmt.Println(err)
				} else {
					if create.Code != protocol.RespCode_SUCCESS {
						fmt.Printf("response error: code[%d], error[%v]\n", create.Code, create.Error)
					} else {
						fmt.Printf("response success: %v\n", create.Error)
					}
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()

	// 登录
	if n, err := ws.Write(login()); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	} else {
		fmt.Printf("write[%d] success.\n", n)
	}
	time.Sleep(time.Second*5)
	// 创建
	if n, err := ws.Write(create()); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	} else {
		fmt.Printf("write[%d] success.\n", n)
	}
	ioutil.ReadAll(os.Stdin)
}

func read(conn *websocket.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println(buffer[:n])
		time.Sleep(time.Second)
	}
}

func login() []byte {
	token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjQwNjA4ODMsImlhdCI6MTU2NDA1NzI4MywiaXAiOiIxOTIuMTY4LjAuMjExIiwidXNlcl9pZCI6Mjd9.9D3KKxHLfVv95yuhRaslkac8RlLw3QZeplVmOWHT_SY`
	login := &protocol.LoginReq{
		Token: *proto.String(token),
	}
	return genMsg(protocol.App_System, protocol.Cmd_Login, login)
}

func create() []byte {
	create := &protocol.CreateReq{
		Options: *proto.Int64(123456),
	}
	return genMsg(protocol.App_NiuNiu, protocol.Cmd_Create, create)
}

func genMsg(app protocol.App, cmd protocol.Cmd, reqData proto.Message) []byte {
	body, _ := proto.Marshal(reqData)
	msg := &protocol.MsgHead{
		CmdId: cmd,
		AppId: app,
		Body:  body,
	}
	data, _ := proto.Marshal(msg)
	return data
}

func readMsg(ws *websocket.Conn){
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := ws.Read(buffer)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			msg := &protocol.MsgHead{}
			proto.Unmarshal(buffer[:n], msg)

			switch msg.CmdId {
			case protocol.Cmd_Login:
				login := &protocol.LoginResp{}
				if err := proto.Unmarshal(msg.Body, login); err != nil {
					fmt.Println(err)
				} else {
					if login.Code != protocol.RespCode_SUCCESS {
						fmt.Printf("response error: code[%d], error[%v]\n", login.Code, login.Error)
					} else {
						fmt.Printf("response success: %v\n", login.Error)
					}
				}
			case protocol.Cmd_Create:
				create := &protocol.CreateResp{}
				if err := proto.Unmarshal(msg.Body, create); err != nil {
					fmt.Println(err)
				} else {
					if create.Code != protocol.RespCode_SUCCESS {
						fmt.Printf("response error: code[%d], error[%v]\n", create.Code, create.Error)
					} else {
						fmt.Printf("response success: %v\n", create.Error)
					}
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()
}