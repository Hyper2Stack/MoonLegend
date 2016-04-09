package handler

import (
    "net/http"
    "encoding/json"

    "controller/model"
)

// GET /api/v1/user/nodes
//
func ListNode(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(LoginUserVars[r].Nodes())
}

// POST /api/v1/user/nodes
//
func PostNode(w http.ResponseWriter, r *http.Request) {
    u := LoginUserVars[r]

    defer r.Body.Close()

    in := struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    if model.GetRepoByNameAndOwner(in.Name, u.Id) != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    n := new(model.Node)
    n.Name = in.Name
    n.Description = in.Description
    n.Owner = u.Id
    n.Save()

    w.WriteHeader(http.StatusCreated)
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
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    if re := model.GetNodeByNameAndOwner(in.Name, NodeVars[r].Owner); re != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    NodeVars[r].Name = in.Name
    NodeVars[r].Description = in.Description
    NodeVars[r].Update()
}

// DELETE /api/v1/user/nodes/{node_name}
//
func DeleteNode(w http.ResponseWriter, r *http.Request) {
    NodeVars[r].Delete()
}

// GET /api/v1/user/nodes/{node_name}/tags
//
func ListNodeTag(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(NodeVars[r].Tags())
}

// POST /api/v1/user/nodes/{node_name}/tags
//
func AddNodeTag(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        TagName string `json:"name"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.TagName) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    t := new(model.NodeTag)
    t.Name   = in.TagName

    NodeVars[r].AddTag(t)

    w.WriteHeader(http.StatusCreated)
}

// DELETE /api/v1/user/nodes/{node_name}/tags/{tag_name}
//
func DeleteNodeTag(w http.ResponseWriter, r *http.Request) {
    NodeVars[r].RemoveTag(TagVars[r])
}

// GET /api/v1/user/nodes/{node_name}/nics
//
func ListNic(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(NodeVars[r].Nics())
}

// POST /api/v1/user/nodes/{node_name}/nics
//
func PostNic(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Name        string `json:"name"`
        Ip4Addr     string `json:"ip"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    if n := NodeVars[r].GetNicByName(in.Name); n != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    n := new(model.Nic)
    n.Name    = in.Name
    n.Ip4Addr = in.Ip4Addr

    NodeVars[r].AddNic(n)

    w.WriteHeader(http.StatusCreated)
}

// GET /api/v1/user/nodes/{node_name}/nics/{nic_name}
//
func GetNic(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(NicVars[r])
}

// PUT /api/v1/user/nodes/{node_name}/nics/{nic_name}
//
func PutNic(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Name        string `json:"name"`
        Ip4Addr     string `json:"ip"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    if n := NodeVars[r].GetNicByName(in.Name); n != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    NicVars[r].Name = in.Name
    NicVars[r].Ip4Addr = in.Ip4Addr
    NodeVars[r].AddNic(NicVars[r])
}

// DELETE /api/v1/user/nodes/{node_name}/nics/{nic_name}
//
func DeleteNic(w http.ResponseWriter, r *http.Request) {
    NodeVars[r].RemoveNic(NicVars[r].Name)
}

// GET /api/v1/user/nodes/{node_name}/nics/{nic_name}/tags
//
func ListNicTag(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(NicVars[r].NicTags())
}

// POST /api/v1/user/nodes/{node_name}/nics/{nic_name}/tags
//
func AddNicTag(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        TagName string `json:"name"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

     if len(in.TagName) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

   t := new(model.NicTag)
    t.Name   = in.TagName

    NicVars[r].AddNicTag(t)

    w.WriteHeader(http.StatusCreated)
}

// DELETE /api/v1/user/nodes/{node_name}/nics/{nic_name}/tags/{tag_name}
//
func DeleteNicTag(w http.ResponseWriter, r *http.Request) {
    NicVars[r].RemoveNicTag(TagVars[r])
}