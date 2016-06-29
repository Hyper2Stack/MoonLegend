package handler

import (
    "encoding/json"
    "net/http"
    "time"

    "controller/agent"
    "controller/model"
    "controller/model/yml"
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

    u := LoginUserVars[r]
    if model.GetGroupByNameAndOwnerId(in.Name, u.Id) != nil {
        w.WriteHeader(http.StatusConflict)
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

    if model.GetGroupByNameAndOwnerId(in.Name, GroupVars[r].OwnerId) != nil {
        w.WriteHeader(http.StatusConflict)
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

    if GroupVars[r].HasNode(n) {
        w.WriteHeader(http.StatusConflict)
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
        Repo    string       `json:"repo"`
        Runtime *yml.Runtime `json:"runtime"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    repo, repoTag := model.ParseRepoString(in.Repo)
    if repo == nil || repoTag == nil || !LoginUserVars[r].CanDeploy(repo) {
        w.WriteHeader(http.StatusForbidden)
        return
    }

    if err := GroupVars[r].InitDeployment(repo, repoTag, in.Runtime); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// PUT /api/v1/user/groups/{group_name}/deployment/prepare
//
func Prepare(w http.ResponseWriter, r *http.Request) {
    if GroupVars[r].Status != model.StatusCreated        &&
       GroupVars[r].Status != model.StatusPrepareTimeout &&
       GroupVars[r].Status != model.StatusPrepareError   &&
       GroupVars[r].Status != model.StatusPrepared {
        http.Error(w, InvalidOperation, http.StatusBadRequest)
        return
    }

    GroupVars[r].InitDeployProcess()

    for uuid, _ := range GroupVars[r].Process {
        if n := model.GetNodeByUuid(uuid); n == nil || !agent.IsConnected(uuid) {
            http.Error(w, "invalid node " + uuid, http.StatusBadRequest)
            return
        }
    }

    go prepare(GroupVars[r])
}

func prepare(group *model.Group) {
    group.Status = model.StatusPreparing
    group.Update()

    for uuid, _ := range group.Process {
        go prepareInstancesOnNode(uuid, group)
    }

    for {
        time.Sleep(5 * time.Second)
        // TBD, fix, group may be updated during prepare process
        group.Update()
        if status, done := group.ParsePrepareProcessStatus(); done {
            group.Status = status
            group.Update()
            break
        }
    }
}

func prepareInstancesOnNode(uuid string, group *model.Group) {
    for _, is := range group.Process[uuid] {
        is.Status = model.StatusPreparing
        i := group.Deployment.FindInstanceByName(is.Name)
        if i.PrepareFile != nil && agent.CreateFile(uuid, i.PrepareFile) != nil {
            is.Status = model.StatusPrepareError
            continue
        }

        if result, err := agent.ExecScript(uuid, i.PrepareCommand); err != nil || !result.Ok {
            is.Status = model.StatusPrepareError
            continue
        }
        is.Status = model.StatusPrepared
    }
}

// PUT /api/v1/user/groups/{group_name}/deployment/execute
//
func Deploy(w http.ResponseWriter, r *http.Request) {
    if GroupVars[r].Status != model.StatusPrepared {
        http.Error(w, InvalidOperation, http.StatusBadRequest)
        return
    }

    go deploy(GroupVars[r])
}

func deploy(group *model.Group) {
    group.Status = model.StatusDeploying
    group.Update()

    services, err := group.Deployment.TopSortServices()
    if err != nil {
        log.Info(err.Error())
        group.Status = model.StatusDeployError
        group.Update()
        return
    }

    for _, service := range services {
        process := group.GetProcessOfServiceMap()
        for _, is := range process[service] {
            is.Status = model.StatusDeploying
            group.Update()
            i := group.Deployment.FindInstanceByName(is.Name)
            result, err := agent.ExecScript(is.NodeUuid, i.RunCommand)
            if err != nil || !result.Ok {
                is.Status = model.StatusDeployError
                group.Status = model.StatusDeployError
                group.Update()
                return
            }
            is.Status = model.StatusDeployed
            group.Update()
        }
    }

    group.Status = model.StatusDeployed
    group.Update()
}

// PUT /api/v1/user/groups/{group_name}/deployment/clear
//
func Clear(w http.ResponseWriter, r *http.Request) {
    if GroupVars[r].Status != model.StatusDeployed      &&
       GroupVars[r].Status != model.StatusDeployTimeout &&
       GroupVars[r].Status != model.StatusDeployError {
        http.Error(w, InvalidOperation, http.StatusBadRequest)
        return
    }

    go clear(GroupVars[r])
}

func clear(group *model.Group) {
    group.Status = model.StatusClearing
    group.Update()
    for _, isList := range group.Process {
        for _, is := range isList {
            if is.Status == model.StatusDeployed {
                is.Status = model.StatusClearing
                group.Update()
                i := group.Deployment.FindInstanceByName(is.Name)
                // TBD, fix, should consider error and timeout
                agent.ExecScript(is.NodeUuid, i.RmCommand)
            }
            is.Status = model.StatusPrepared
            group.Update()
        }
    }

    group.Status = model.StatusPrepared
    group.Update()
}

// GET /api/v1/user/groups/{group_name}/deployment/process
//
func GetProcess(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(GroupVars[r].Process)
}

// DELETE /api/v1/user/groups/{group_name}/deployment
//
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
    clear(GroupVars[r])

    GroupVars[r].Deployment = nil
    GroupVars[r].Process = nil
    GroupVars[r].Status = model.StatusRaw
    GroupVars[r].Update()
}
