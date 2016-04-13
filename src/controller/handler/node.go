package handler

import (
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "controller/model"
)

// GET /api/v1/user/nodes
//
func ListNode(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(LoginUserVars[r].Nodes())
}

// GET /api/v1/user/nodes/{node_name}
//
func GetNode(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(NodeVars[r])
}

// PUT /api/v1/user/nodes/{node_name}
//
func PutNode(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        Status      string `json:"status"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    if re := model.GetNodeByNameAndOwnerId(in.Name, NodeVars[r].OwnerId); re != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    NodeVars[r].Name = in.Name
    NodeVars[r].Description = in.Description
    NodeVars[r].Status = in.Status
    NodeVars[r].Update()
}

// DELETE /api/v1/user/nodes/{node_name}
//
func DeleteNode(w http.ResponseWriter, r *http.Request) {
    NodeVars[r].Delete()
}

// POST /api/v1/user/nodes/{node_name}/tags
//
func AddNodeTag(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Name string `json:"name"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    NodeVars[r].AddTag(in.Name)
    w.WriteHeader(http.StatusCreated)
}

// DELETE /api/v1/user/nodes/{node_name}/tags/{tag_name}
//
func DeleteNodeTag(w http.ResponseWriter, r *http.Request) {
    NodeVars[r].RemoveTag(mux.Vars(r)["tag_name"])
}

// POST /api/v1/user/nodes/{node_name}/nics/{nic_name}/tags
//
func AddNicTag(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Name string `json:"name"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

     if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    NodeVars[r].AddNicTag(mux.Vars(r)["nic_name"], in.Name)
    w.WriteHeader(http.StatusCreated)
}

// DELETE /api/v1/user/nodes/{node_name}/nics/{nic_name}/tags/{tag_name}
//
func DeleteNicTag(w http.ResponseWriter, r *http.Request) {
    NodeVars[r].RemoveNicTag(mux.Vars(r)["nic_name"], mux.Vars(r)["tag_name"])
}
