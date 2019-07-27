package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"testing"
	"time"
)

func BenchmarkLogin(b *testing.B) {
	origin := "http://192.168.0.211:7777/"
	addr := "ws://192.168.0.211:7777/ws"
	ws, err := websocket.Dial(addr, "", origin)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.N = 1000
	for i := 0; i < b.N; i++ {
		// 登录
		if _, err := ws.Write(login()); err != nil {
			b.Fatal(err)
		}
	}
	readMsg(ws)
	time.Sleep(1000)
	ws.Close()
	fmt.Println("test end")

}
