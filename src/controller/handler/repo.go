package handler

import (
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    "controller/model"
)

// GET /api/v1/user/repos
//
func ListRepo(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(LoginUserVars[r].Repos())
}

// POST /api/v1/user/repos
//
func PostRepo(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    u := LoginUserVars[r]
    in := struct {
        Name        string `json:"name"`
        IsPublic    bool   `json:"is_public"`
        Description string `json:"description"`
        Readme      string `json:"readme"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    if model.GetRepoByNameAndOwnerId(in.Name, u.Id) != nil {
        w.WriteHeader(http.StatusConflict)
        return
    }

    re := new(model.Repo)
    re.Name = in.Name
    re.Description = in.Description
    re.IsPublic = in.IsPublic
    re.Readme = in.Readme
    re.OwnerId = u.Id
    re.Save()

    w.WriteHeader(http.StatusCreated)
}

// GET /api/v1/user/repos/{repo_name}
//
func GetRepo(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(RepoVars[r])
}

// PUT /api/v1/user/repos/{repo_name}
//
func PutRepo(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        IsPublic    bool   `json:"is_public"`
        Description string `json:"description"`
        Readme      string `json:"readme"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    RepoVars[r].IsPublic = in.IsPublic
    RepoVars[r].Description = in.Description
    RepoVars[r].Readme = in.Readme
    RepoVars[r].Update()
}

// DELETE /api/v1/user/repos/{repo_name}
//
func DeleteRepo(w http.ResponseWriter, r *http.Request) {
    RepoVars[r].Delete()
}

// GET /api/v1/user/repos/{repo_name}/tags
//
func ListRepoTag(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(RepoVars[r].Tags())
}

// POST /api/v1/user/repos/{repo_name}/tags
//
func AddRepoTag(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Name string `json:"name"`
        Yml  string `json:"yml"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 || len(in.Yml) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    if RepoVars[r].GetTag(in.Name) != nil {
        w.WriteHeader(http.StatusConflict)
        return
    }

    t := new(model.RepoTag)
    t.Name = in.Name
    t.Yml = in.Yml
    RepoVars[r].AddTag(t)

    w.WriteHeader(http.StatusCreated)
}

// GET /api/v1/user/repos/{repo_name}/tags/{tag_name}
//
func GetRepoTag(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(RepoVars[r].GetTag(mux.Vars(r)["tag_name"]))
}

// DELETE /api/v1/user/repos/{repo_name}/tags/{tag_name}
//
func DeleteRepoTag(w http.ResponseWriter, r *http.Request) {
    RepoVars[r].RemoveTag(mux.Vars(r)["tag_name"])
}

// GET /api/v1/repos/{namespace}/{name}
//
func GetGlobalRepo(w http.ResponseWriter, r *http.Request) {
    // TBD
}

// GET /api/v1/repos/{namespace}/{name}/tags
//
func ListGlobalRepoTag(w http.ResponseWriter, r *http.Request) {
    // TBD
}

// GET /api/v1/repos/{namespace}/{name}/tags/{tag}
//
func GetGlobalRepoTag(w http.ResponseWriter, r *http.Request) {
    // TBD
}
