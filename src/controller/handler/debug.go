package handler

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"

    "controller/agent"
    "controller/model"
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
    case agent.ActionPing:
        if err := agent.Ping(uuid); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    case agent.ActionNodeInfo:
        info, err := agent.GetNodeInfo(uuid)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(info)
    case agent.ActionAgentInfo:
        info, err := agent.GetAgentInfo(uuid)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(info)
    case agent.ActionExecShell:
        script := new(model.ScriptJob)
        if err := json.NewDecoder(r.Body).Decode(script); err != nil {
            http.Error(w, "Decode json error", http.StatusBadRequest)
            return
        }
        if _, err := agent.ExecScript(uuid, script); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    default:
        http.Error(w, "action not supported", http.StatusBadRequest)
        return
    }
}
