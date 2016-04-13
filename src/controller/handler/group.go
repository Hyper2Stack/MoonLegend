package handler

import (
    "encoding/json"
    "net/http"

    "controller/model"
)

// GET /api/v1/user/groups
//
func ListGroup(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(LoginUserVars[r].Groups())
}

// POST /api/v1/user/groups
//
func PostGroup(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    u := LoginUserVars[r]

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

    if model.GetGroupByNameAndOwnerId(in.Name, u.Id) != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    g := new(model.Group)
    g.Name = in.Name
    g.Description = in.Description
    g.OwnerId = u.Id
    g.Status = model.StatusRaw
    g.Save()

    w.WriteHeader(http.StatusCreated)
}

// GET /api/v1/user/groups/{group_name}
//
func GetGroup(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(GroupVars[r])
}

// PUT /api/v1/user/groups/{group_name}
//
func PutGroup(w http.ResponseWriter, r *http.Request) {
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

   if re := model.GetGroupByNameAndOwnerId(in.Name, GroupVars[r].OwnerId); re != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    GroupVars[r].Name = in.Name
    GroupVars[r].Description = in.Description
    GroupVars[r].Update()

}

// DELETE /api/v1/user/groups/{group_name}
//
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
    GroupVars[r].Delete()
}

// GET /api/v1/user/groups/{group_name}/nodes
//
func ListGroupNode(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(GroupVars[r].Nodes())
}

// POST /api/v1/user/groups/{group_name}/nodes
//
func AddGroupNode(w http.ResponseWriter, r *http.Request) {
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

    n := model.GetNodeByNameAndOwnerId(in.Name, LoginUserVars[r].Id)
    if n == nil {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    GroupVars[r].AddNode(n)
    w.WriteHeader(http.StatusCreated)
}

// DELETE /api/v1/user/groups/{group_name}/nodes/{node_name}
//
func DeleteGroupNode(w http.ResponseWriter, r *http.Request) {
    GroupVars[r].DeleteNode(NodeVars[r])
}

// GET /api/v1/user/groups/{group_name}/deployment
//
func GetDeployment(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(GroupVars[r].Deployment)
}

// POST /api/v1/user/groups/{group_name}/deployment
//
func PostDeployment(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Repo string `json:"repo"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if GroupVars[r].Status != model.StatusRaw {
        http.Error(w, GroupBusy, http.StatusBadRequest)
        return
    }

    // construct Deployment
    // TBD
    w.WriteHeader(http.StatusCreated)
}

// PUT /api/v1/user/groups/{group_name}/deployment/prepare
//
func Prepare(w http.ResponseWriter, r *http.Request) {
    if GroupVars[r].Status != model.StatusCreated {
        http.Error(w, InvalidOperation, http.StatusBadRequest)
        return
    }

    // TBD
}

// PUT /api/v1/user/groups/{group_name}/deployment/execute
//
func Deploy(w http.ResponseWriter, r *http.Request) {
    if GroupVars[r].Status == model.StatusPrepared {
        http.Error(w, InvalidOperation, http.StatusBadRequest)
        return
    }

    // TBD
}

// GET /api/v1/user/groups/{group_name}/deployment/process
//
func GetProcess(w http.ResponseWriter, r *http.Request) {
    // TBD
}

// DELETE /api/v1/user/groups/{group_name}/deployment
//
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
    // TBD
}
