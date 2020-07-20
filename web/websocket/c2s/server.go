package main

import (
	"bufio"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", "localhost:9000", "http service address")

var upgrader = websocket.Upgrader{}

func main() {
	flag.Parse()
	http.HandleFunc("/echo", echo)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	log.Println("Connected by" + r.RemoteAddr)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	stdReader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			input, err := stdReader.ReadString('\n')
			if err != nil {
				log.Println("input:", err)
			}
			err = c.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Println("write:", err)
			}
		}
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
	}
}
