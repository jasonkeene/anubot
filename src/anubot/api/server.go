package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

func EchoServer(ws *websocket.Conn) {
	fmt.Println("got ws connection")
	io.Copy(ws, ws)
}

func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	fmt.Println("listening on port 12345")
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
