package main

import (
    "fmt"
    "code.google.com/p/go.net/websocket"
    "net/http"
    "log"
)

var (
	JSON          = websocket.JSON           // codec for JSON
	Message       = websocket.Message        // codec for string, []byte
	ActiveClients = make(map[ClientConn]int) // map containing clients
)

// Client connection consists of the websocket and the client ip
type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

// Initialize handlers and websocket handlers
func init() {
	http.Handle("/", http.FileServer(http.Dir("../app")))
	http.Handle("/sock", websocket.Handler(SockServer))
}

// WebSocket server to handle chat between clients
func SockServer(ws *websocket.Conn) {
    fmt.Println("This is this socket server.")

    var err error
    var clientMessage string
    // use []byte if websocket binary type is blob or arraybuffer
    // var clientMessage []byte

    // cleanup on server side
    defer func() {
        if err = ws.Close(); err != nil {
            log.Println("Websocket could not be closed", err.Error())
        }
    }()

    client := ws.Request().RemoteAddr
    log.Println("Client connected:", client)
    sockCli := ClientConn{ws, client}
    ActiveClients[sockCli] = 0
    log.Println("Number of clients connected ...", len(ActiveClients))

    // for loop so the websocket stays open otherwise
    // it'll close after one Receieve and Send
    for {
        if err = Message.Receive(ws, &clientMessage); err != nil {
            // If we cannot Read then the connection is closed
            log.Println("Websocket Disconnected waiting", err.Error())
            // remove the ws client conn from our active clients
            delete(ActiveClients, sockCli)
            log.Println("Number of clients still connected ...", len(ActiveClients))
            return
        }

        clientMessage = sockCli.clientIP + " Said: " + clientMessage
        for cs, _ := range ActiveClients {
            if err = Message.Send(cs.websocket, clientMessage); err != nil {
                // we could not send the message to a peer
                log.Println("Could not send message to ", cs.clientIP, err.Error())
            }
        }
    }
}

func main() {
    err := http.ListenAndServe(":3000", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}