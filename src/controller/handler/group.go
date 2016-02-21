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

    if model.GetGroupByNameAndOwner(in.Name, u.Id) != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    g := new(model.Group)
    g.Name = in.Name
    g.Description = in.Description
    g.Owner = u.Id
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

   if re := model.GetGroupByNameAndOwner(in.Name, GroupVars[r].Owner); re != nil {
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
    json.NewEncoder(w).Encode(GroupVars[r].Nodes)
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
    
    n := model.GetNodeByNameAndOwner(in.Name, LoginUserVars[r].Id)
    if n == nil {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    GroupVars[r].AddMembership(n.Id)

    w.WriteHeader(http.StatusCreated)
}

// DELETE /api/v1/user/groups/{group_name}/nodes/{node_name}
//
func DeleteGroupNode(w http.ResponseWriter, r *http.Request) {
    GroupVars[r].DeleteMembership(NodeVars[r].Id)
}

// GET /api/v1/user/groups/{group_name}/deployment
//
func GetDeployment(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(GroupVars[r].Deployment)
}

// POST /api/v1/user/groups/{group_name}/deployment
//
func AddDeployment(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        RepoName string `json:"name"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if GroupVars[r].Deployment != nil {
        GroupVars[r].DeleteDeployment()
    }
    
    if len(in.RepoName) != 0 {
        re := model.GetRepoByNameAndOwner(in.RepoName, LoginUserVars[r].Id)
        if re == nil {
            http.Error(w, RequestBodyError, http.StatusBadRequest)
            return
        }
        if err := GroupVars[r].ParseYml(re.YmlPath); err == nil {
            panic("Invalid yml config")
        }
        GroupVars[r].Deployment.RepoName = in.RepoName
    }
}

// PUT /api/v1/user/groups/{group_name}/deployment/execute
//
func ExecuteDeployment(w http.ResponseWriter, r *http.Request) {
    if GroupVars[r].Deployment == nil {
        http.Error(w, InvalidOperation, http.StatusBadRequest)
        return
    }

    re := model.GetRepoByNameAndOwner(GroupVars[r].Deployment.RepoName, LoginUserVars[r].Id)
    if re == nil {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    GroupVars[r].Execute(re.Id)
}

// DELETE /api/v1/user/groups/{group_name}/deployment
//
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
    GroupVars[r].DeleteDeployment()
}
