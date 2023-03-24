package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func handleWSConnection(ws *websocket.Conn) {

	ws.Write([]byte("connected"))

	chatter(ws)
}

func chatter(ws *websocket.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := ws.Read(buf)
		if err != nil {
			// client closed connection
			if err == io.EOF {
				break
			}
			continue
		}
		ws.Write(buf[:n])
	}
}

func pinger(ws *websocket.Conn) {

	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			// client closed connection
			if err == io.EOF {
				break
			}
			continue
		}

		if string(buf[:n]) == "ping" {
			ws.Write([]byte("pong"))
		}
	}
}

// just test writing
func tick(ws *websocket.Conn) {
	for {
		time.Sleep(time.Second * 1)
		ws.Write([]byte("tick"))
	}
}

func main() {

	fmt.Println("vim-go")

	http.Handle("/chat", websocket.Handler(handleWSConnection))

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", nil)

}
