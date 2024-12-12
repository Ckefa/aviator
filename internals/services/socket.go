package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	Clients map[*websocket.Conn]bool
	Mu      sync.Mutex
}

type Info struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

type Message struct {
	Name string  `json:"name"`
	Msg  float32 `json:"msg"`
}

var Clients = ClientManager{
	Clients: make(map[*websocket.Conn]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleSocket(w http.ResponseWriter, r *http.Request) {
	// log.Println("Upgrading conn: ", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS conn Err: ", err)
		return
	}

	log.Println("CONNECTED ", conn.RemoteAddr())

	Clients.Mu.Lock()
	Clients.Clients[conn] = true
	Clients.Mu.Unlock()

	defer conn.Close()
	defer func() {
		Clients.Mu.Lock()
		delete(Clients.Clients, conn)
		Clients.Mu.Unlock()
		log.Println("DISCONNECTED: ", conn.RemoteAddr())
	}()

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("ERROR READING: ", err)
			break
		}

		log.Println("MSG :", fmt.Sprintf("%s", msg))

		resp := Info{
			Name: "info",
			Msg:  "RECEIVED",
		}

		respJson, _ := json.Marshal(resp)

		if err := conn.WriteMessage(websocket.TextMessage, respJson); err != nil {
			log.Fatal(err)
		}
	}
}

func BroadCastMsg(data map[string]interface{}) {
	Clients.Mu.Lock()
	defer Clients.Mu.Unlock()

	for conn := range Clients.Clients {
		if err := conn.WriteJSON(data); err != nil {
			log.Println("BroadCast ERR", err)
			conn.Close()
			delete(Clients.Clients, conn)
		}
	}
}

func RunSocketIO() {
	http.HandleFunc("/ws", handleSocket)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
