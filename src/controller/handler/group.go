package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "strings"
    "text/template"

    "gopkg.in/yaml.v2"
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

    if model.GetGroupByNameAndOwnerId(in.Name, GroupVars[r].OwnerId) != nil {
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
        Tag  string `json:"tag"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if GroupVars[r].Status != model.StatusRaw {
        http.Error(w, GroupBusy, http.StatusBadRequest)
        return
    }

    ss := strings.Split(in.Repo, "/")
    if len(ss) != 2 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    namespace := ss[0]
    name := ss[1]

    // TBD, check if namespace/name exists
    // TBD, check if current user has deloy permission
    // TBD, check if tag exists

    repo := model.GetRepoByNameAndOwnerId(name, model.GetUserByName(namespace).Id)
    repoTag := repo.GetTag(name)

    d := construct([]byte(repoTag.Yml), GroupVars[r])
    rederedYml := render(d, repoTag.Yml)
    d = construct(rederedYml, GroupVars[r])

    GroupVars[r].Deployment = d
    GroupVars[r].Status = model.StatusCreated
    GroupVars[r].Update()

    w.WriteHeader(http.StatusCreated)
}

func construct(y []byte, g *model.Group) *model.Deployment {
    yml := new(Yml)
    if err := yaml.Unmarshal(y, yml); err != nil {
        panic(err)
        return nil
    }

    // TBD
    return nil
}

func render(d *model.Deployment, yml string) []byte {
    var b bytes.Buffer
    t, _ := template.New("").Parse(yml)
    t.Execute(&b, d)
    return b.Bytes()
}

// PUT /api/v1/user/groups/{group_name}/deployment/prepare
//
func Prepare(w http.ResponseWriter, r *http.Request) {
    if GroupVars[r].Status != model.StatusCreated {
        http.Error(w, InvalidOperation, http.StatusBadRequest)
        return
    }

    go prepare(GroupVars[r])
}

func prepare(group *model.Group) {
    group.Status = model.StatusPreparing
    group.Update()

    // TBD pull docker images
    // possible error

    group.Status = model.StatusPrepared
    group.Update()
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

    // TBD pull docker images
    // possible error, roll back

    group.Status = model.StatusDeployed
    group.Update()
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
