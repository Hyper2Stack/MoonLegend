package handler

import (
    "fmt"
    "net/http"
    "runtime/debug"
    "time"

    "github.com/gorilla/mux"
)

type Route struct {
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

var routes = []Route{
    // ping
    Route{"GET",    "/api/v1/ping", Ping},

    // login
    Route{"POST",   "/api/v1/login",               wrapper(Login)},
    Route{"POST",   "/api/v1/signup",              wrapper(Signup)},

    // user
    Route{"GET",    "/api/v1/user",                wrapper(authWrapper(GetMyProfile))},
    Route{"PUT",    "/api/v1/user/reset-password", wrapper(authWrapper(ResetPassword))},
    Route{"PUT",    "/api/v1/user/reset-key",      wrapper(authWrapper(ResetKey))},

    // repo
    Route{"GET",    "/api/v1/user/repos",             wrapper(authWrapper(ListRepo))},
    Route{"POST",   "/api/v1/user/repos",             wrapper(authWrapper(PostRepo))},
    Route{"GET",    "/api/v1/user/repos/{repo_name}", wrapper(authWrapper(repoWrapper(GetRepo)))},
    Route{"PUT",    "/api/v1/user/repos/{repo_name}", wrapper(authWrapper(repoWrapper(PutRepo)))},
    Route{"DELETE", "/api/v1/user/repos/{repo_name}", wrapper(authWrapper(repoWrapper(DeleteRepo)))},

    Route{"GET",    "/api/v1/user/repos/{repo_name}/tags",            wrapper(authWrapper(repoWrapper(ListRepoTag)))},
    Route{"POST",   "/api/v1/user/repos/{repo_name}/tags",            wrapper(authWrapper(repoWrapper(AddRepoTag)))},
    Route{"DELETE", "/api/v1/user/repos/{repo_name}/tags/{tag_name}", wrapper(authWrapper(repoWrapper(DeleteRepoTag)))},

    // node
    Route{"GET",    "/api/v1/user/nodes",             wrapper(authWrapper(ListNode))},
    Route{"GET",    "/api/v1/user/nodes/{node_name}", wrapper(authWrapper(nodeWrapper(GetNode)))},
    Route{"PUT",    "/api/v1/user/nodes/{node_name}", wrapper(authWrapper(nodeWrapper(PutNode)))},
    Route{"DELETE", "/api/v1/user/nodes/{node_name}", wrapper(authWrapper(nodeWrapper(DeleteNode)))},

    Route{"POST",   "/api/v1/user/nodes/{node_name}/tags",            wrapper(authWrapper(nodeWrapper(AddNodeTag)))},
    Route{"DELETE", "/api/v1/user/nodes/{node_name}/tags/{tag_name}", wrapper(authWrapper(nodeWrapper(DeleteNodeTag)))},
    Route{"POST",   "/api/v1/user/nodes/{node_name}/nics/{nic_name}/tags",            wrapper(authWrapper(nodeWrapper(AddNicTag)))},
    Route{"DELETE", "/api/v1/user/nodes/{node_name}/nics/{nic_name}/tags/{tag_name}", wrapper(authWrapper(nodeWrapper(DeleteNicTag)))},

    // group
    Route{"GET",    "/api/v1/user/groups",              wrapper(authWrapper(ListGroup))},
    Route{"POST",   "/api/v1/user/groups",              wrapper(authWrapper(PostGroup))},
    Route{"GET",    "/api/v1/user/groups/{group_name}", wrapper(authWrapper(groupWrapper(GetGroup)))},
    Route{"PUT",    "/api/v1/user/groups/{group_name}", wrapper(authWrapper(groupWrapper(PutGroup)))},
    Route{"DELETE", "/api/v1/user/groups/{group_name}", wrapper(authWrapper(groupWrapper(DeleteGroup)))},

    Route{"GET",    "/api/v1/user/groups/{group_name}/nodes",             wrapper(authWrapper(groupWrapper(ListGroupNode)))},
    Route{"POST",   "/api/v1/user/groups/{group_name}/nodes",             wrapper(authWrapper(groupWrapper(AddGroupNode)))},
    Route{"DELETE", "/api/v1/user/groups/{group_name}/nodes/{node_name}", wrapper(authWrapper(groupWrapper(nodeWrapper(DeleteGroupNode))))},

    // deployment
    Route{"GET",    "/api/v1/user/groups/{group_name}/deployment",         wrapper(authWrapper(groupWrapper(GetDeployment)))},
    Route{"POST",   "/api/v1/user/groups/{group_name}/deployment",         wrapper(authWrapper(groupWrapper(PostDeployment)))},
    Route{"PUT",    "/api/v1/user/groups/{group_name}/deployment/prepare", wrapper(authWrapper(groupWrapper(Prepare)))},
    Route{"PUT",    "/api/v1/user/groups/{group_name}/deployment/execute", wrapper(authWrapper(groupWrapper(Deploy)))},
    Route{"GET",    "/api/v1/user/groups/{group_name}/deployment/process", wrapper(authWrapper(groupWrapper(GetProcess)))},
    Route{"DELETE", "/api/v1/user/groups/{group_name}/deployment",         wrapper(authWrapper(groupWrapper(DeleteDeployment)))},

    // repo
    Route{"GET",   "/api/v1/repos/{namespace}/{name}",            wrapper(authWrapper(globalRepoWrapper(GetGlobalRepo)))},
    Route{"GET",   "/api/v1/repos/{namespace}/{name}/tags",       wrapper(authWrapper(globalRepoWrapper(ListGlobalRepoTag)))},
    Route{"GET",   "/api/v1/repos/{namespace}/{name}/tags/{tag}", wrapper(authWrapper(globalRepoWrapper(GetGlobalRepoTag)))},
}

type InnerResponseWriter struct {
    statusCode int
    setted     bool
    http.ResponseWriter
}

func (i *InnerResponseWriter) WriteHeader(status int) {
    if !i.setted {
        i.statusCode = status
        i.setted = true
    }

    i.ResponseWriter.WriteHeader(status)
}

func (i *InnerResponseWriter) Write(b []byte) (int, error) {
    i.setted = true
    return i.ResponseWriter.Write(b)
}

func wrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        s := time.Now()
        wr := &InnerResponseWriter{
            statusCode:     200,
            setted:         false,
            ResponseWriter: w,
        }

        defer func() {
            if err := recover(); err != nil {
                debug.PrintStack()
                wr.WriteHeader(http.StatusInternalServerError)
                log.Criticalf("Panic: %v", err)
                fmt.Fprintf(w, fmt.Sprintln(err))
            }

            d := time.Now().Sub(s)
            log.Infof("Wrapper %s %s %d %s", r.Method, r.RequestURI, wr.statusCode, d.String())
        }()

        inner.ServeHTTP(wr, r)
    })
}

func NewRouter() *mux.Router {
    router := mux.NewRouter()
    for _, route := range routes {
        router.Methods(route.Method).Path(route.Pattern).HandlerFunc(route.HandlerFunc)
    }

    return router
}
