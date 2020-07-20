package main

import (
	"bufio"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

var addr = flag.String("addr", "localhost:9000", "http service address")

func main() {
	flag.Parse()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatal("read:", err)
				return
			}
			log.Printf("recv: %s", msg)
		}
	}()

	stdReader := bufio.NewReader(os.Stdin)
	msg := make(chan string)

	go func() {
		for {
			input, err := stdReader.ReadString('\n')
			if err != nil {
				log.Println("input:", err)
				msg <- ""
				continue
			}
			msg <- input
		}
	}()

	for {
		select {
		case r := <-msg:
			err := conn.WriteMessage(websocket.TextMessage, []byte(r))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
