package handler

import (
    "encoding/json"
    "net/http"

    "controller/agent"

    "github.com/gorilla/mux"
    "github.com/hyper2stack/mooncommon/protocol"
)

// PUT /api/v1/debug-agent/{uuid}/{action}
//
func DebugAgent(w http.ResponseWriter, r *http.Request) {
    uuid := mux.Vars(r)["uuid"]
    action := mux.Vars(r)["action"]

    if !agent.IsConnected(uuid) {
        http.Error(w, "node is not connected", http.StatusBadRequest)
        return
    }

    switch action {
    case protocol.MethodNodeInfo:
        info, err := agent.GetNodeInfo(uuid)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(info)
    case protocol.MethodAgentInfo:
        info, err := agent.GetAgentInfo(uuid)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(info)
    case protocol.MethodExecScript:
        script := new(protocol.Script)
        if err := json.NewDecoder(r.Body).Decode(script); err != nil {
            http.Error(w, "Decode json error", http.StatusBadRequest)
            return
        }
        result, err := agent.ExecScript(uuid, script)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        json.NewEncoder(w).Encode(result)
    case protocol.MethodCreateFile:
        file := new(protocol.File)
        if err := json.NewDecoder(r.Body).Decode(file); err != nil {
            http.Error(w, "Decode json error", http.StatusBadRequest)
            return
        }
        err := agent.CreateFile(uuid, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    default:
        http.Error(w, "action not supported", http.StatusBadRequest)
        return
    }
}
