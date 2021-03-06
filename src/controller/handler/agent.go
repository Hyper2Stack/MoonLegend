package handler

import (
    "fmt"
    "net/http"
    "time"
    "strings"

    "github.com/gorilla/websocket"

    "controller/model"
    "controller/agent"
)

var (
    wsUpgrader = websocket.Upgrader{}
)

func ConnectAgent(w http.ResponseWriter, r *http.Request) {
    c, err := wsUpgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Errorf("WEBSOCKET upgrade failed, %v", err)
        return
    }
    defer c.Close()

    // Moon-Authentication: key,uuid; auth[0]=key, auth[1]=uuid
    auth := strings.Split(r.Header.Get("Moon-Authentication"), ",")
    if len(auth) != 2 {
        log.Errorf("Invalid authentication")
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // Step 1: verify user
    u := model.GetUserByAgentkey(auth[0])
    if u == nil {
        log.Errorf("Unauthorized key %s", auth[0])
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // Step 2: register connection and read loop
    disconnected := make(chan struct{})
    go func() {
        err := agent.Listen(auth[1], c)
        log.Infof("Node %s disconnected, %v", auth[1], err)
        disconnected <- struct{}{}
    }()

    time.Sleep(1 * time.Second)

    // Step 3: update node info
    node := model.GetNodeByUuidAndOwnerId(auth[1], u.Id)
    if node == nil {
        ni, err := agent.GetNodeInfo(auth[1])
        if err != nil {
            log.Errorf("Get node info of %s failed, %v", auth[1], err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        node = new(model.Node)
        node.Status = model.StatusRaw
        node.Uuid = auth[1]
        node.OwnerId = u.Id
        node.Name = ni.Hostname
        if model.GetNodeByNameAndOwnerId(ni.Hostname, u.Id) != nil {
            node.Name = fmt.Sprintf("%s_%s", ni.Hostname, auth[1])
        }
        node.Description = fmt.Sprintf("%s's host: %s", u.Name, ni.Hostname)
        node.ImportAgentNics(ni.Nics)
        node.Save()
    }

    <- disconnected
}
