package agent

import (
    "encoding/json"
    "fmt"
    "strings"
    "sync"

    "github.com/gorilla/websocket"
)

const (
    ActionPing          = "ping"
    ActionNodeInfo      = "get-node-info"
    ActionAgentInfo     = "get-agent-info"
    ActionExecShell     = "exec-shell-script"
    StatusOK            = "ok"
    StatusBadRequest    = "bad-request"
    StatusError         = "error"
    StatusInternalError = "internal-error"
)

type Request struct {
    Action  string
    Content []byte
}

type Response struct {
    Status  string
    Content []byte
}

func DecodeRequest(msg []byte) *Request {
    ss := strings.SplitN(string(msg), "\r\n", 2)
    return &Request{Action: ss[0], Content: []byte(ss[1])}
}

func EncodeRequest(req *Request) []byte {
    return []byte(fmt.Sprintf("%s\r\n%s", req.Action, req.Content))
}

func DecodeResponse(msg []byte) *Response {
    ss := strings.SplitN(string(msg), "\r\n", 2)
    return &Response{Status: ss[0], Content: []byte(ss[1])}
}

func EncodeResponse(res *Response) []byte {
    return []byte(fmt.Sprintf("%s\r\n%s", res.Status, res.Content))
}

////////////////////////////////////////////////////////////////////////////////

type NodeInfo struct {
    Hostname string `json:"hostname"`
    Nics     []*Nic `json:"nics"`
}

type Nic struct {
    Name    string   `json:"name"`
    Ip4Addr string   `json:"ip4addr"`
    Tags    []string `json:"tags"`
}

type AgentInfo struct {
    Version string `json:"version"`
}

type ShellCommand struct {
    Command  string   `json:"command"`
    Args     []string `json:"args"`
    Restrict bool     `json:"restrict"`
}

type ScriptJob struct {
    Commands []*ShellCommand `json:"commands"`
}

type ScriptJobResult struct {
    ErrCommand *ShellCommand `json:"err_command"`
}

////////////////////////////////////////////////////////////////////////////////

type AgentConnection struct {
    Ws         *websocket.Conn
    Lock       *sync.Mutex
}

var (
    wsConns = make(map[string]*AgentConnection)
)

func IsConnected(uuid string) bool {
    return wsConns[uuid] != nil
}

func Register(uuid string, conn *websocket.Conn) {
    wsConns[uuid] = &AgentConnection{Ws: conn, Lock: new(sync.Mutex)}
}

func Unregister(uuid string, conn *websocket.Conn) {
    if wsConns[uuid] != nil && wsConns[uuid].Ws == conn {
        delete(wsConns, uuid)
    }
}

func Ping(uuid string) error {
    // TBD
    return nil
}

func GetNodeInfo(uuid string) (*NodeInfo, error) {
    req := new(Request)
    req.Action = ActionNodeInfo
    res, err := wsConns[uuid].do(req)
    if err != nil {
        return nil, err
    }

    if res.Status != StatusOK {
        return nil, fmt.Errorf("bad response %s", res.Status)
    }

    nodeinfo := new(NodeInfo)
    if err := json.Unmarshal(res.Content, nodeinfo); err != nil {
        return nil, err
    }

    return nodeinfo, nil
}

func GetAgentInfo(uuid string) (*AgentInfo, error) {
    req := new(Request)
    req.Action = ActionAgentInfo
    res, err := wsConns[uuid].do(req)
    if err != nil {
        return nil, err
    }

    if res.Status != StatusOK {
        return nil, fmt.Errorf("bad response %s", res.Status)
    }

    agentinfo := new(AgentInfo)
    if err := json.Unmarshal(res.Content, agentinfo); err != nil {
        return nil, err
    }

    return agentinfo, nil
}

func ExecScript(uuid string, script *ScriptJob) (*ScriptJobResult, error) {
    req := new(Request)
    req.Action = ActionExecShell
    req.Content, _ = json.Marshal(script)
    res, err := wsConns[uuid].do(req)
    if err != nil {
        return nil, err
    }

    if res.Status != StatusOK {
        scriptResult := new(ScriptJobResult)
        json.Unmarshal(res.Content, scriptResult)
        return scriptResult, fmt.Errorf("bad response %s", res.Status)
    }

    return nil, nil
}

func (ac *AgentConnection) do(req *Request) (*Response, error) {
    ac.Lock.Lock()
    defer ac.Lock.Unlock()

    if err := ac.Ws.WriteMessage(websocket.TextMessage, EncodeRequest(req)); err != nil {
        return nil, fmt.Errorf("websocket write %s", err.Error())
    }

    _, m, err := ac.Ws.ReadMessage()
    if err != nil {
        return nil, fmt.Errorf("websocket read %s", err.Error())
    }

    return DecodeResponse(m), nil
}
