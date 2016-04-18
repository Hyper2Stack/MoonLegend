package handler

import (
    "fmt"
    "golang.org/x/net/websocket"
    "github.com/gorilla/mux"
)

func WsWrapper(h func (*websocket.Conn)) websocket.Handler {
    return websocket.Handler(h)
}

func AgentHandler(ws *websocket.Conn) {
    userKey := mux.Vars(ws.Request())["user_key"]
    fmt.Println(userKey)
    // TODO agent operations

    /*
    e.g. echo message
    msg := make([]byte, 512)
    n, err := ws.Read(msg)
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Printf("Receive: %s\n", msg[:n])

    m, err := ws.Write(msg[:n])
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Printf("Send: %s\n", msg[:m])
    */
}
