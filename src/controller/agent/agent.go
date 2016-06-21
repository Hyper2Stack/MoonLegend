package agent

import (
    "encoding/json"
    "errors"
    "fmt"
    "sync"
    "time"

    "github.com/op/go-logging"
    "github.com/gorilla/websocket"
    "github.com/satori/go.uuid"
    "github.com/hyper2stack/mooncommon/protocol"
)

////////////////////////////////////////////////////////////////////////////////

type ConnWrapper struct {
    Conn *websocket.Conn
    Lock *sync.Mutex
    Jobs map[string]chan *Result
}

type Result struct {
    Status  string
    Payload []byte
}

var log, _  = logging.GetLogger("moonlegend")
var wsConns = make(map[string]*ConnWrapper)

var ErrNotConnected = errors.New("not connected")
var ErrTimeout      = errors.New("request timeout")
var DefaultTimeout  = 30
var ScriptTimeout   = 900

func IsConnected(uuid string) bool {
    return wsConns[uuid] != nil && wsConns[uuid].Conn != nil
}

func Listen(uuid string, conn *websocket.Conn) error {
    if wrapper, ok := wsConns[uuid]; ok {
        wrapper.Lock.Lock()
        wrapper.Conn = conn
        wrapper.Lock.Unlock()
    } else {
        wsConns[uuid] = &ConnWrapper{
            Conn: conn,
            Lock: new(sync.Mutex),
            Jobs: make(map[string]chan *Result),
        }
    }

    for {
        _, msg, err := wsConns[uuid].Conn.ReadMessage()
        if err != nil {
            wsConns[uuid].Conn = nil
            return err
        }

        go wsConns[uuid].receive(msg)
    }
}

func GetNodeInfo(uuid string) (*protocol.Node, error) {
    if !IsConnected(uuid) {
        return nil, ErrNotConnected
    }

    status, body, err := wsConns[uuid].send(protocol.MethodNodeInfo, nil, DefaultTimeout)
    if err != nil {
        return nil, err
    }

    if status != protocol.StatusOK {
        return nil, parseErrorBody(body)
    }

    info := new(protocol.Node)
    if err := json.Unmarshal(body, info); err != nil {
        return nil, err
    }

    return info, nil
}

func GetAgentInfo(uuid string) (*protocol.Agent, error) {
    if !IsConnected(uuid) {
        return nil, ErrNotConnected
    }

    status, body, err := wsConns[uuid].send(protocol.MethodAgentInfo, nil, DefaultTimeout)
    if err != nil {
        return nil, err
    }

    if status != protocol.StatusOK {
        return nil, parseErrorBody(body)
    }

    info := new(protocol.Agent)
    if err := json.Unmarshal(body, info); err != nil {
        return nil, err
    }

    return info, nil
}

func CreateFile(uuid string, file *protocol.File) error {
    if !IsConnected(uuid) {
        return ErrNotConnected
    }

    payload, _ := json.Marshal(file)
    status, body, err := wsConns[uuid].send(protocol.MethodCreateFile, payload, DefaultTimeout)
    if err != nil {
        return err
    }

    if status != protocol.StatusOK {
        return parseErrorBody(body)
    }

    return nil
}

func ExecScript(uuid string, script *protocol.Script) (*protocol.ScriptResult, error) {
    if !IsConnected(uuid) {
        return nil, ErrNotConnected
    }

    payload, _ := json.Marshal(script)
    status, body, err := wsConns[uuid].send(protocol.MethodExecScript, payload, ScriptTimeout)
    if err != nil {
        return nil, err
    }

    if status != protocol.StatusOK {
        return nil, parseErrorBody(body)
    }

    result := new(protocol.ScriptResult)
    if err := json.Unmarshal(body, result); err != nil {
        return nil, err
    }

    return result, nil
}


func parseErrorBody(body []byte) error {
    message := struct {
        Err string `json:"error"`
    }{}
    if err := json.Unmarshal(body, &message); err != nil {
        return err
    }

    return fmt.Errorf("client error, %s", message.Err)
}

func (cw *ConnWrapper) receive(body []byte) {
    msg, err := protocol.Decode(body)
    if err != nil {
        log.Errorf("Message decode error")
        return
    }

    if msg.Type != protocol.Res {
        return
    }

    if ch, ok := cw.Jobs[msg.Uuid]; ok {
        ch <- &Result{Status: msg.Method, Payload: msg.Payload}
        return
    }

    log.Infof("Job %s not found", msg.Uuid)
}

func (cw *ConnWrapper) send(method string, payload []byte, timeout int) (string, []byte, error) {
    uuid := uuid.NewV4().String()

    msg := new(protocol.Msg)
    msg.Type = protocol.Req
    msg.Uuid = uuid
    msg.Method = method
    msg.Payload = payload

    body, err := protocol.Encode(msg)
    if err != nil {
        log.Errorf("Message encode error")
        return "", nil, err
    }

    ch := make(chan *Result)
    cw.Jobs[uuid] = ch
    defer delete(cw.Jobs, uuid)

    cw.Lock.Lock()
    if err := cw.Conn.WriteMessage(websocket.TextMessage, body); err != nil {
        cw.Lock.Unlock()
        return "", nil, err
    }
    cw.Lock.Unlock()

    select {
    case res := <- ch:
        return res.Status, res.Payload, nil
    case <-time.After(time.Duration(timeout) * time.Second):
        return "", nil, ErrTimeout
    }

    return "", nil, nil
}
