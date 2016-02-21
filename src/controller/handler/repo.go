package handler

import (
    "net/http"
    "os"
    "path/filepath"
    "encoding/json"

    "controller/model"
)

// GET /api/v1/user/repos
//
func ListRepo(w http.ResponseWriter, r *http.Request) {
    // TODO: convert ymlPath to yml
    json.NewEncoder(w).Encode(LoginUserVars[r].Repos())
}

// POST /api/v1/user/repos
//
func PostRepo(w http.ResponseWriter, r *http.Request) {
    u := LoginUserVars[r]

    in := struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        Yml         string `json:"yml"`
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

    y := filepath.Join(YmlDirBase, u.Name + "_" + in.Name)
    file, err := os.OpenFile(y, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
    if err != nil {
        panic(err)
    }

    defer func() {
        r.Body.Close()
        file.Close()
    }()

    if _, err = file.WriteString(in.Yml); err != nil {
        panic(err)
    }

    re := new(model.Repo)
    re.Name = in.Name
    re.Description = in.Description
    re.Owner = u.Id
    re.YmlPath = y
    re.Save()

    w.WriteHeader(http.StatusCreated)
}

// GET /api/v1/user/repos/{repo_name}
//
func GetRepo(w http.ResponseWriter, r *http.Request) {
    // TODO: convert ymlPath to yml
    json.NewEncoder(w).Encode(RepoVars[r])
}

// PUT /api/v1/user/repos/{repo_name}
//
func PutRepo(w http.ResponseWriter, r *http.Request) {    
    in := struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        Yml         string `json:"yml"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }
    
    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

   if re := model.GetRepoByNameAndOwner(in.Name, RepoVars[r].Owner); re != nil {
        http.Error(w, DuplicateResource, http.StatusBadRequest)
        return
    }

    file, err := os.OpenFile(RepoVars[r].YmlPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
    if err != nil {
        panic(err)
    }
    
    defer func() {
        r.Body.Close()
        file.Close()
    }()

    if _, err = file.WriteString(in.Yml); err != nil {
        panic(err)
    }

    RepoVars[r].Name = in.Name
    RepoVars[r].Description = in.Description
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

    t := new(model.RepoTag)
    t.Name   = in.TagName

    RepoVars[r].AddTag(t)

    w.WriteHeader(http.StatusCreated)
}

// DELETE /api/v1/user/repos/{repo_name}/tags/{tag_name}
//
func DeleteRepoTag(w http.ResponseWriter, r *http.Request) {
    RepoVars[r].RemoveTag(TagVars[r])
}

// GET /api/v1/users/{user_name}/repos
//
func ListUserRepo(w http.ResponseWriter, r *http.Request) {
    // TODO: convert ymlPath to yml
    json.NewEncoder(w).Encode(UserVars[r].Repos())
}

// GET /api/v1/users/{user_name}/repos/{repo_name}
//
func GetUserRepo(w http.ResponseWriter, r *http.Request) {
    // TODO: convert ymlPath to yml
    json.NewEncoder(w).Encode(RepoVars[r])
}