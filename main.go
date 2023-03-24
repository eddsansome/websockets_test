package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var connections = []*websocket.Conn{}

func handleWeatherUpdates(ws *websocket.Conn) {

	weather(ws)
}

func weather(ws *websocket.Conn) {
	bodyCh := make(chan []byte)

	for {
		time.Sleep(time.Second * 5)
		go getWeather(bodyCh)

		select {
		case w := <-bodyCh:
			ws.Write(w)
		}
	}
}

func getWeather(bodyCh chan []byte) {

	req, err := http.NewRequest(http.MethodGet,
		"https://api.open-meteo.com/v1/forecast?latitude=51.52&longitude=-0.34&current_weather=true",
		nil)

	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	bodyCh <- b

}

func handleWSConnection(ws *websocket.Conn) {

	connections = append(connections, ws)

	chatter(ws)
}

func broadcast(msg []byte) {
	for _, conn := range connections {
		conn.Write(msg)
	}
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
		msg := buf[:n]
		go broadcast(msg)
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
	http.Handle("/weather_updates", websocket.Handler(handleWeatherUpdates))

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", nil)

}
