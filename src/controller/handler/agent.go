package handler

import (
    "errors"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{}
var wsConns = make(map[string]*websocket.Conn)
/*
  now store (connId=mix(userKey, nodeUUid), websocket.Conn) in a map,
  may store paris in redis in future for large scale
 */

func ConnectAgent(w http.ResponseWriter, r *http.Request) {
    c, err := wsUpgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Errorf("WEBSOCKET %s", "upgrade failed")
        return
    }
    s := time.Now()
    defer func() {
        WebsocketDisconnect(c);
        d := time.Now().Sub(s)
        log.Infof("WEBSOCKET %s(lasted=%s)", "closed", d.String())
    }()
    log.Infof("WEBSOCKET %s", "connected")

    for {
        mT, m, err := c.ReadMessage()
        if err != nil {
            log.Errorf("WEBSOCKET READ %s", err.Error())
            break
        }
        log.Debugf("WEBSOCKET %d %s", mT, m)
        m, err = digest(mT, m)
        if err != nil {
            break
        }
        err = WebsocketSendMessage(c, mT, m)
        if err != nil {
            log.Errorf("WEBSOCKET WRITE %s", err.Error())
            break
        }
    }
}

func WebsocketGetConnection (connId string) *websocket.Conn {
    c := wsConns[connId]
    return c
}

func WebsocketSendMessage (c *websocket.Conn, mT int, m []byte) error {
    // mT: websocket.(TextMessage|BinaryMessage|CloseMessage)
    if c == nil {
        return errors.New("connection is null")
    }
    return c.WriteMessage(mT, m)
}

func WebsocketDisconnect (c *websocket.Conn) {
    c.Close()
    for key,val := range(wsConns) {
        if val == c {
            delete(wsConns, key)
            break
        }
    }
}

func makeConnId (userKey string, nodeUuid string) string {
    return userKey + "." + nodeUuid
}

func digest (mT int, m []byte) ([]byte, error) {
    // TODO agent body
    return m, nil
}
