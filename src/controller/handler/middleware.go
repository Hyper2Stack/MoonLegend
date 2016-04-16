package handler

import (
    "net/http"

    "github.com/gorilla/mux"

    "controller/model"
)

var LoginUserVars  map[*http.Request]*model.User
var RepoVars       map[*http.Request]*model.Repo
var NodeVars       map[*http.Request]*model.Node
var GroupVars      map[*http.Request]*model.Group
var GlobalRepoVars map[*http.Request]*model.Repo

func init() {
    LoginUserVars  = make(map[*http.Request]*model.User)
    RepoVars       = make(map[*http.Request]*model.Repo)
    NodeVars       = make(map[*http.Request]*model.Node)
    GroupVars      = make(map[*http.Request]*model.Group)
    GlobalRepoVars = make(map[*http.Request]*model.Repo)
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
        defer delete(LoginUserVars, r)

        inner.ServeHTTP(w, r)
    })
}

func repoWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        re := model.GetRepoByNameAndOwnerId(mux.Vars(r)["repo_name"], LoginUserVars[r].Id)
        if re == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        RepoVars[r] = re
        defer delete(RepoVars, r)

        inner.ServeHTTP(w, r)
    })
}

func nodeWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        n := model.GetNodeByNameAndOwnerId(mux.Vars(r)["node_name"], LoginUserVars[r].Id)
        if n == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        NodeVars[r] = n
        defer delete(NodeVars, r)

        inner.ServeHTTP(w, r)
    })
}

func groupWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        g := model.GetGroupByNameAndOwnerId(mux.Vars(r)["group_name"], LoginUserVars[r].Id)
        if g == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        GroupVars[r] = g
        defer delete(NodeVars, r)

        inner.ServeHTTP(w, r)
    })
}

func globalRepoWrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        namespace := model.GetUserByName(mux.Vars(r)["namespace"])
        if namespace == nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        repo := model.GetRepoByNameAndOwnerId(mux.Vars(r)["name"], namespace.Id)
        if repo == nil || !repo.IsPublic {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        GlobalRepoVars[r] = repo
        defer delete(GlobalRepoVars, r)

        inner.ServeHTTP(w, r)
    })
}
