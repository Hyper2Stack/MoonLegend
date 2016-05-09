package handler

import (
    "errors"
    "net/http"
    "time"
    "strings"
    "fmt"
    "encoding/json"

    "github.com/gorilla/websocket"

    "controller/model"
)

type Request struct {
    Action  string
    Content []byte
}

type Response struct {
    Status  string
    Content []byte
}

type NodeInfo struct {
    Hostname string       `json:"hostname"`
    Nics     []*model.Nic `json:"nics"`
}

type AgentConnection struct {
    Ws    *websocket.Conn
    Send  chan []byte
    Done  chan string
}

var (
    wsUpgrader = websocket.Upgrader{}
    wsConns = make(map[string]*AgentConnection)
)

const (
    RetryCount          = 3
    ActionNodeInfo      = "get-node-info"
    ActionAgentInfo     = "get-agent-info"
    ActionExecShell     = "exec-shell-script"
    StatusOK            = "ok"
    StatusBadRequest    = "bad-request"
    StatusError         = "error"
    StatusInternalError = "internal-error"
)

func decodeResponse(msg []byte) *Response {
    ss := strings.SplitN(string(msg), "\r\n", 2)
    return &Response{Status: ss[0], Content: []byte(ss[1])}
}

func encodeRequest(req *Request) []byte {
    return []byte(fmt.Sprintf("%s\r\n%s", req.Action, req.Content))
}

/*
  now store (connId=mix(userKey,nodeUUid), AgentConnection) in a map,
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
        WebsocketDisconnect(c)
        d := time.Now().Sub(s)
        log.Infof("WEBSOCKET %s(lasted=%s)", "closed", d.String())
    }()
    log.Infof("WEBSOCKET %s", "connected")

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

    // Step 2: update node info
    req := new(Request)
    req.Action = ActionNodeInfo
    res, err := WebsocketSendRecv(c, req)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if res.Status != StatusOK {
        log.Errorf("fail to get node-info from agent")
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    nodeinfo := new(NodeInfo)
    err = json.Unmarshal(res.Content, nodeinfo)
    if err != nil {
        log.Errorf("fail to parse node-info from agent: %s", err.Error())
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    if n := model.GetNodeByUuid(auth[1]); n == nil {
        n := new(model.Node)
        n.Uuid = auth[1]
        n.Name = nodeinfo.Hostname
        n.Description = u.Name + "'s host: " + nodeinfo.Hostname
        n.OwnerId = u.Id
        n.Nics = nodeinfo.Nics
        n.Save()
    } else {
        n.Name = nodeinfo.Hostname
        n.Nics = nodeinfo.Nics
        n.Update()
    }

    // Step 3: register wsConns
    ac := &AgentConnection{Ws: c, Send: make(chan []byte), Done: make(chan string)}
	wsConns[r.Header.Get("Moon-Authentication")] = ac

    // Step 4: wait and exec deploy commands
    for {
        req.Action = ActionExecShell
        req.Content = <- ac.Send

        res, err = WebsocketSendRecv(c, req)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            break
        }

        if res.Status != StatusOK {
            log.Errorf("Agent %s failed to exec cmds %s: %s", r.Header.Get("Moon-Authentication"), req.Content, string(res.Content))
        }
        ac.Done <- res.Status
    }
}

func WebsocketSendRecv (c *websocket.Conn, req *Request) (*Response, error) {
    res := new(Response)
    for retry := RetryCount; retry > 0; retry-- {
        if err := WebsocketSendMessage(c, websocket.TextMessage, encodeRequest(req)); err != nil {
            log.Errorf("WEBSOCKET WRITE %s", err.Error())
            return nil, err
        }

        _, m, err := c.ReadMessage()
        if err != nil {
            log.Errorf("WEBSOCKET READ %s", err.Error())
            return nil, err
        }

        res = decodeResponse(m)
        if res.Status == StatusOK {
            return res, nil
        }
    }
    return res, nil
}

func WebsocketGetConnection (connId string) *websocket.Conn {
    c := wsConns[connId].Ws
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
        if val.Ws == c {
            delete(wsConns, key)
            break
        }
    }
}

func makeConnId (userKey string, nodeUuid string) string {
    return userKey + "," + nodeUuid
}

func digest (mT int, m []byte) ([]byte, error) {
    // TODO agent body
    return m, nil
}
