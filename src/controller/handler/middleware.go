package handler

import (
    "net/http"
    "fmt"
    "os"

    "github.com/gorilla/mux"

    "controller/model"
)

var YmlDirBase    string
var LoginUserVars map[*http.Request]*model.User
var RepoVars      map[*http.Request]*model.Repo
var NodeVars      map[*http.Request]*model.Node
var NicVars       map[*http.Request]*model.Nic
var GroupVars     map[*http.Request]*model.Group
var UserVars      map[*http.Request]*model.User
var TagVars       map[*http.Request]string

func init() {
    YmlDirBase = absCleanPath("./ymls")
    if os.MkdirAll(YmlDirBase, 0770) != nil {
        panic("fail to create dir for yml files")
    }
    LoginUserVars = make(map[*http.Request]*model.User)
    RepoVars      = make(map[*http.Request]*model.Repo)
    NodeVars      = make(map[*http.Request]*model.Node)
    NicVars       = make(map[*http.Request]*model.Nic)
    GroupVars     = make(map[*http.Request]*model.Group)
    UserVars      = make(map[*http.Request]*model.User)
    TagVars       = make(map[*http.Request]string)
}

func authWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        name, valid := decodeUserToken(r.Header.Get("Authorization"))
        if !valid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        u := model.GetUserByName(name)
        if u == nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        LoginUserVars[r] = u

        fmt.Printf("authWrapper started\n")

        defer func() {
            delete(LoginUserVars, r)
            fmt.Printf("authWrapper completed\n")
        }()

        inner.ServeHTTP(w, r)
    })
}

func repoWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if LoginUserVars[r] == nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        name := mux.Vars(r)["repo_name"]
        re := model.GetRepoByNameAndOwner(name, LoginUserVars[r].Id)
        if re == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        RepoVars[r] = re
        fmt.Printf("repoWrapper received message %s\n", name)

        defer func() {
            delete(RepoVars, r)
            fmt.Printf("repoWrapper completed\n")
        }()

        inner.ServeHTTP(w, r)
    })
}

func nodeWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if LoginUserVars[r] == nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        name := mux.Vars(r)["node_name"]
        n := model.GetNodeByNameAndOwner(name, LoginUserVars[r].Id)
        if n == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        NodeVars[r] = n
        fmt.Printf("nodeWrapper received message %s\n", name)

        defer func() {
            delete(NodeVars, r)
            fmt.Printf("nodeWrapper completed\n")
        }()

        inner.ServeHTTP(w, r)
    })
}

func nicWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if LoginUserVars[r] == nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        
        if NodeVars[r] == nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        name := mux.Vars(r)["nic_name"]
        n := NodeVars[r].GetNicByName(name)
        if n == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        NicVars[r] = n

        defer delete(NicVars, r)

        inner.ServeHTTP(w, r)
    })
}

func groupWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer delete(NodeVars, r)

        if LoginUserVars[r] == nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        name := mux.Vars(r)["group_name"]
        g := model.GetGroupByNameAndOwner(name, LoginUserVars[r].Id)
        if g == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        GroupVars[r] = g

        inner.ServeHTTP(w, r)
    })
}

func adminWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer delete(UserVars, r)

        if LoginUserVars[r] == nil || !model.IsAdmin(LoginUserVars[r].Id ) {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        name := mux.Vars(r)["user_name"]
        u := model.GetUserByName(name)
        if u == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        UserVars[r] = u

        inner.ServeHTTP(w, r)
    })
}

func tagWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        TagVars[r] = mux.Vars(r)["tag_name"]
 
        defer delete(TagVars, r)

        inner.ServeHTTP(w, r)
    })
}
