package handler

import (
    "net/http"
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
    GroupVars     = make(map[*http.Request]*model.Group)
    UserVars      = make(map[*http.Request]*model.User)
}

func authWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer delete(LoginUserVars, r)

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

        inner.ServeHTTP(w, r)
    })
}

func repoWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer delete(RepoVars, r)

        if LoginUserVars[r] == nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        name := mux.Vars(r)["repo_name"]
        re := model.GetRepoByNameAndOwnerId(name, LoginUserVars[r].Id)
        if re == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        RepoVars[r] = re

        inner.ServeHTTP(w, r)
    })
}

func nodeWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer delete(NodeVars, r)

        if LoginUserVars[r] == nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        name := mux.Vars(r)["node_name"]
        n := model.GetNodeByNameAndOwnerId(name, LoginUserVars[r].Id)
        if n == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        NodeVars[r] = n

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
        g := model.GetGroupByNameAndOwnerId(name, LoginUserVars[r].Id)
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
        if LoginUserVars[r] == nil || !model.IsAdmin(LoginUserVars[r].Id ) {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        inner.ServeHTTP(w, r)
    })
}

func userWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer delete(UserVars, r)

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
