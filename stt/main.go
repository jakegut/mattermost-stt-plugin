package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asticode/go-asticoqui"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var model *asticoqui.Model

func byteaToInt16a(bytes []byte) []int16 {
	res := make([]int16, len(bytes)/2)
	for i, j := 0, 0; i < len(bytes)-1; i += 2 {
		var val int16
		val |= int16(bytes[i]) << 8
		val |= int16(bytes[i+1])
		res[j] = val
		j += 1
		// res = append(res, int16(bytes[i]))
	}
	return res
}

func message(conn *websocket.Conn, msg chan []byte, disconnect chan bool) {
	for {
		mt, message, err := conn.ReadMessage()

		if mt == websocket.CloseMessage {
			disconnect <- true
			return
		}

		if _, ok := err.(*websocket.CloseError); ok {
			disconnect <- true
			return
		}
		if err != nil {
			log.Println("warning: read message error:", err)
			continue
		}
		if mt == websocket.BinaryMessage {
			msg <- message
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	stream, err := model.NewStream()
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer stream.Discard()
	defer c.Close()
	disconnect := make(chan bool)
	msg := make(chan []byte)
	go message(c, msg, disconnect)
	for {
		select {
		case incoming := <-msg:
			log.Println("before:", len(incoming))
			log.Printf("%v", incoming)
			bin := byteaToInt16a(incoming)
			log.Printf("%v", bin)
			log.Println("after: ", len(bin))
			stream.FeedAudioContent(bin)
			if err != nil {
				log.Println("read:", err)
				break
			}
		case <-ticker.C:
			res, err := stream.IntermediateDecode()
			log.Println("prediction:", res)
			if err != nil {
				log.Println("write:", err)
			}
			err = c.WriteMessage(websocket.TextMessage, []byte(res))
			if err != nil {
				log.Println("write:", err)
				break
			}
		case <-disconnect:
			return
		}
	}
}

func main() {
	flag.Parse()
	var err error
	model, err = asticoqui.New(os.Getenv("TRANSCRIBER_MODEL_PATH"))
	if err != nil {
		panic(err)
	}
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
