package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

//Configuration - Global configuration
var Configuration = Config{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}

		for _, originAllowed := range Configuration.OriginsAllowed {
			if strings.Contains(origin, originAllowed) {
				return true
			}
		}
		return false
	},
}

// SetupAndHandleConnections is used to set up listeners and socket connections
func SetupAndHandleConnections() {

	http.HandleFunc("/ws", httpSocketHandler)

	configReadErr := ReadJSONConfig(&Configuration)
	if configReadErr != nil {
		return
	}

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		return
	}
}

func httpSocketHandler(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("handler")
	sockconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error in connecting to socket", 500)
	}
	socketConnection(sockconn)

	//defer sockconn.Close()
}

func readMessages(sockconn *websocket.Conn, clientinfo *ClientInfo) {
	for {
		var msg SurgeConnect

		err := sockconn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(msg)
		switch msg.Type {
		case "sub":
			SubscribeToChannel(msg.ChannelName, clientinfo)
		case "unsub":
			UnsubscribeFromChannel(msg.ChannelName, clientinfo)
		}

	}
}

func socketConnection(conn *websocket.Conn) {
	var clientinfo ClientInfo
	ircMessageChannel := make(chan SurgeMessage)

	//Point the clients channel to the one that will listen for surges
	clientinfo.surgeChannel = ircMessageChannel

	go readMessages(conn, &clientinfo)
	go listenForIrcChannelSurges(conn, &clientinfo, ircMessageChannel)

}

//SendIrcChannelSurges - Write surges back to client
func listenForIrcChannelSurges(conn *websocket.Conn, client *ClientInfo, ircMessageChannel <-chan SurgeMessage) {
	var msg SurgeMessage

	msg.ChannelName = "test"
	//ircMessageChannel <- msg
	for {
		select {
		case msg := <-ircMessageChannel:
			fmt.Println("message received")
			sendMessage(conn, msg)
		}
	}
}

func sendMessage(conn *websocket.Conn, message SurgeMessage) {
	conn.WriteJSON(message)
}
