package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path"

    "ml/client"
    "gopkg.in/alecthomas/kingpin.v1"
)

var (
    AuthPath = fmt.Sprintf("%s/.ml/auth.json", os.Getenv("HOME"))
)

func readToken(host string) (string, error) {
    content, err := ioutil.ReadFile(AuthPath)
    if err != nil {
        if os.IsNotExist(err) {
            return "", nil
        }
        return "", err
    }

    var info map[string]string
    if err := json.Unmarshal(content, &info); err != nil {
        return "", err
    }

    if token, ok := info[host]; ok {
        return token, nil
    }

    return "", nil
}

func writeToken(host, token string) error {
    info := make(map[string]string)
    content, err := ioutil.ReadFile(AuthPath)
    if err == nil {
        if err := json.Unmarshal(content, &info); err != nil {
            return err
        }
    }

    info[host] = token
    body, _ := json.Marshal(info)
    if err := os.MkdirAll(path.Dir(AuthPath), 0700); err != nil {
        return err
    }
    if err := ioutil.WriteFile(AuthPath, body, 0644); err != nil {
        return err
    }

    return nil
}

func deleteToken(host string) error {
    content, err := ioutil.ReadFile(AuthPath)
    if err != nil {
        return nil
    }

    var info map[string]string
    if err := json.Unmarshal(content, &info); err != nil {
        return err
    }

    delete(info, host)

    body, _ := json.Marshal(info)
    if err := ioutil.WriteFile(AuthPath, body, 0644); err != nil {
        return err
    }

    return nil
}

func mustLogin(server string) string {
    token, err := readToken(server)
    if err != nil {
        log.Fatalln(err)
    }

    if token == "" {
        log.Fatalln("not login")
    }

    return token
}

var (
    app    = kingpin.New("ml", "A command-line client of moonlegend.")
    server = app.Flag("server", "Server address.").Default("localhost:8080").OverrideDefaultFromEnvar("MOONLEGEND_SERVER").String()

    ping   = app.Command("ping", "Ping moonlegend server.")

    signup        = app.Command("signup", "Register a new user.")
    signupUser    = signup.Arg("username", "Signup username.").Required().String()
    signupPasword = signup.Arg("password", "Signup password.").Required().String()

    login         = app.Command("login", "Login.")
    loginUser     = login.Arg("username", "Login username.").Required().String()
    loginPassword = login.Arg("password", "Login password.").Required().String()

    logout  = app.Command("logout", "Logout.")
    profile = app.Command("profile", "Print user profile.")
    resetK  = app.Command("reset-key", "Reset key.")

    resetP        = app.Command("reset-password", "Reset password.")
    resetPassword = resetP.Arg("new-password", "New password.").Required().String()

    listRepo = app.Command("list-repo", "List repo.")

    createRepo         = app.Command("create-repo", "Create repo.")
    createRepoIsPublic = createRepo.Flag("public", "Public repo.").Bool()
    createRepoName     = createRepo.Arg("name", "Repo name.").Required().String()

    deleteRepo     = app.Command("delete-repo", "Delete repo.")
    deleteRepoName = deleteRepo.Arg("name", "Repo name.").Required().String()

    listRepoTag     = app.Command("list-repo-tag", "List repo tag.")
    listRepoTagRepo = listRepoTag.Arg("repo", "Repo name.").Required().String()

    showRepoTag     = app.Command("show-repo-tag", "Print yml of specifed repo tag.")
    showRepoTagRepo = showRepoTag.Arg("repo", "Repo name.").Required().String()
    showRepoTagTag  = showRepoTag.Arg("tag", "Tag name.").Required().String()

    createRepoTag     = app.Command("create-repo-tag", "Create repo tag.")
    createRepoTagRepo = createRepoTag.Arg("repo", "Repo name.").Required().String()
    createRepoTagTag  = createRepoTag.Arg("tag", "Tag name.").Required().String()
    createRepoTagYml  = createRepoTag.Arg("path", "Yml file path.").Required().String()

    deleteRepoTag     = app.Command("delete-repo-tag", "Delete repo tag.")
    deleteRepoTagRepo = deleteRepoTag.Arg("repo", "Repo name.").Required().String()
    deleteRepoTagTag  = deleteRepoTag.Arg("tag", "Tag name.").Required().String()

    listNode = app.Command("list-node", "List node.")

    showNode     = app.Command("show-node", "Print details of specified node.")
    showNodeName = showNode.Arg("name", "Node name.").Required().String()

    deleteNode     = app.Command("delete-node", "Delete node.")
    deleteNodeName = deleteNode.Arg("name", "Node name.").Required().String()

    listNodeTag     = app.Command("list-node-tag", "List tags of specified node.")
    listNodeTagNode = listNodeTag.Arg("node", "Node name.").Required().String()

    createNodeTag     = app.Command("create-node-tag", "Create tag on specified node.")
    createNodeTagNode = createNodeTag.Arg("node", "Node name.").Required().String()
    createNodeTagTag  = createNodeTag.Arg("tag", "Tag name.").Required().String()

    deleteNodeTag     = app.Command("delete-node-tag", "Delete tag on specified node.")
    deleteNodeTagNode = deleteNodeTag.Arg("node", "Node name.").Required().String()
    deleteNodeTagTag  = deleteNodeTag.Arg("tag", "Tag name.").Required().String()

    listNicTag     = app.Command("list-nic-tag", "List tags of specified nic.")
    listNicTagNode = listNicTag.Arg("node", "Node name.").Required().String()
    listNicTagNic  = listNicTag.Arg("nic", "Nic name.").Required().String()

    createNicTag     = app.Command("create-nic-tag", "Create tag on specified node nic.")
    createNicTagNode = createNicTag.Arg("node", "Node name.").Required().String()
    createNicTagNic  = createNicTag.Arg("nic", "Nic name.").Required().String()
    createNicTagTag  = createNicTag.Arg("tag", "Tag name.").Required().String()

    deleteNicTag     = app.Command("delete-nic-tag", "Delete tag on specified node nic.")
    deleteNicTagNode = deleteNicTag.Arg("node", "Node name.").Required().String()
    deleteNicTagNic  = deleteNicTag.Arg("nic", "Nic name.").Required().String()
    deleteNicTagTag  = deleteNicTag.Arg("tag", "Tag name.").Required().String()

    listGroup = app.Command("list-group", "List group.")

    createGroup     = app.Command("create-group", "Create group.")
    createGroupName = createGroup.Arg("name", "Group name.").Required().String()

    showGroup     = app.Command("show-group", "Print group details.")
    showGroupName = showGroup.Arg("name", "Group name.").Required().String()

    deleteGroup     = app.Command("delete-group", "Delete group.")
    deleteGroupName = deleteGroup.Arg("name", "Group name.").Required().String()

    listGroupNode      = app.Command("list-group-node", "List group node.")
    listGroupNodeGroup = listGroupNode.Arg("group", "Group name.").Required().String()

    createGroupNode      = app.Command("add-node-to-group", "Add specified node to group.")
    createGroupNodeGroup = createGroupNode.Arg("group", "Group name.").Required().String()
    createGroupNodeNode  = createGroupNode.Arg("node", "Node name.").Required().String()

    deleteGroupNode      = app.Command("remove-node-from-group", "Remove specified node from group.")
    deleteGroupNodeGroup = deleteGroupNode.Arg("group", "Group name.").Required().String()
    deleteGroupNodeNode  = deleteGroupNode.Arg("node", "Node name.").Required().String()

    deployInit      = app.Command("deploy-init", "Init deployloyment for specified group.")
    deployInitGroup = deployInit.Arg("group", "Group name.").Required().String()
    deployInitRepo  = deployInit.Arg("repo", "Absolute repo name with format: namespace/repo:tag.").Required().String()

    deployPrepare      = app.Command("deploy-prepare", "Prepare deployloyment for specified group.")
    deployPrepareGroup = deployPrepare.Arg("group", "Group name.").Required().String()

    deployExecute      = app.Command("deploy-execute", "Execute deployloyment for specified group.")
    deployExecuteGroup = deployExecute.Arg("group", "Group name.").Required().String()

    deployClear      = app.Command("deploy-clear", "Clear deployloyment for specified group.")
    deployClearGroup = deployClear.Arg("group", "Group name.").Required().String()

    deployDelete      = app.Command("deploy-delete", "Delete deployloyment for specified group.")
    deployDeleteGroup = deployDelete.Arg("group", "Group name.").Required().String()
)

func main() {
    log.SetFlags(0)

    switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case ping.FullCommand():
        if err := client.New(*server, "").Ping(); err != nil {
            log.Fatalln(err)
        }
        fmt.Println("pong")

    case signup.FullCommand():
        if err := client.New(*server, "").Signup(*signupUser, *signupPasword); err != nil {
            log.Fatalln(err)
        }

    case login.FullCommand():
        token, err := client.New(*server, "").Login(*loginUser, *loginPassword)
        if err != nil {
            log.Fatalln(err)
        }
        if err := writeToken(*server, token); err != nil {
            log.Fatalln(err)
        }
        fmt.Printf("token stored in %s\n", AuthPath)

    case logout.FullCommand():
        if err := deleteToken(*server); err != nil {
            log.Fatalln(err)
        }

    case profile.FullCommand():
        pf, err := client.New(*server, mustLogin(*server)).Profile()
        if err != nil {
            log.Fatalln(err)
        }
        printProfile(pf)

    case resetP.FullCommand():
        if err := client.New(*server, mustLogin(*server)).ResetPassword(*resetPassword); err != nil {
            log.Fatalln(err)
        }

    case resetK.FullCommand():
        if err := client.New(*server, mustLogin(*server)).ResetKey(); err != nil {
            log.Fatalln(err)
        }

    case listRepo.FullCommand():
        repos, err := client.New(*server, mustLogin(*server)).Repos()
        if err != nil {
            log.Fatalln(err)
        }
        printRepos(repos)

    case createRepo.FullCommand():
        if err := client.New(*server, mustLogin(*server)).CreateRepo(*createRepoName, *createRepoIsPublic); err != nil {
            log.Fatalln(err)
        }

    case deleteRepo.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteRepo(*deleteRepoName); err != nil {
            log.Fatalln(err)
        }

    case listRepoTag.FullCommand():
        tags, err := client.New(*server, mustLogin(*server)).RepoTags(*listRepoTagRepo)
        if err != nil {
            log.Fatalln(err)
        }
        printRepoTags(tags)

    case showRepoTag.FullCommand():
        tag, err := client.New(*server, mustLogin(*server)).RepoTag(*showRepoTagRepo, *showRepoTagTag)
        if err != nil {
            log.Fatalln(err)
        }
        printRepoTag(tag)

    case createRepoTag.FullCommand():
        // To be fixed, here will add a new extra new line in yml
        yml, err := ioutil.ReadFile(*createRepoTagYml)
        if err != nil {
            log.Fatalln(err)
        }
        if err := client.New(*server, mustLogin(*server)).CreateRepoTag(*createRepoTagRepo, *createRepoTagTag, string(yml)); err != nil {
            log.Fatalln(err)
        }

    case deleteRepoTag.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteRepoTag(*deleteRepoTagRepo, *deleteRepoTagTag); err != nil {
            log.Fatalln(err)
        }

    case listNode.FullCommand():
        nodes, err := client.New(*server, mustLogin(*server)).Nodes()
        if err != nil {
            log.Fatalln(err)
        }
        printNodes(nodes)

    case showNode.FullCommand():
        node, err := client.New(*server, mustLogin(*server)).Node(*showNodeName)
        if err != nil {
            log.Fatalln(err)
        }
        printNode(node)

    case deleteNode.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteNode(*deleteNodeName); err != nil {
            log.Fatalln(err)
        }

    case listNodeTag.FullCommand():
        node, err := client.New(*server, mustLogin(*server)).Node(*listNodeTagNode)
        if err != nil {
            log.Fatalln(err)
        }
        printNodeTags(node)

    case createNodeTag.FullCommand():
        if err := client.New(*server, mustLogin(*server)).CreateNodeTag(*createNodeTagNode, *createNodeTagTag); err != nil {
            log.Fatalln(err)
        }

    case deleteNodeTag.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteNodeTag(*deleteNodeTagNode, *deleteNodeTagTag); err != nil {
            log.Fatalln(err)
        }

    case listNicTag.FullCommand():
        node, err := client.New(*server, mustLogin(*server)).Node(*listNicTagNode)
        if err != nil {
            log.Fatalln(err)
        }
        nic := findNic(node, *listNicTagNic)
        if nic == nil {
            log.Fatalf("nic %s not found\n", *listNicTagNic)
        }
        printNicTags(nic)

    case createNicTag.FullCommand():
        if err := client.New(*server, mustLogin(*server)).CreateNicTag(*createNicTagNode, *createNicTagNic, *createNicTagTag); err != nil {
            log.Fatalln(err)
        }

    case deleteNicTag.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteNicTag(*deleteNicTagNode, *deleteNicTagNic, *deleteNicTagTag); err != nil {
            log.Fatalln(err)
        }

    case listGroup.FullCommand():
        groups, err := client.New(*server, mustLogin(*server)).Groups()
        if err != nil {
            log.Fatalln(err)
        }
        printGroups(groups)

    case showGroup.FullCommand():
        group, err := client.New(*server, mustLogin(*server)).Group(*showGroupName)
        if err != nil {
            log.Fatalln(err)
        }
        printGroup(group)

    case createGroup.FullCommand():
        if err := client.New(*server, mustLogin(*server)).CreateGroup(*createGroupName); err != nil {
            log.Fatalln(err)
        }

    case deleteGroup.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteGroup(*deleteGroupName); err != nil {
            log.Fatalln(err)
        }

    case listGroupNode.FullCommand():
        nodes, err := client.New(*server, mustLogin(*server)).GroupNodes(*listGroupNodeGroup)
        if err != nil {
            log.Fatalln(err)
        }
        printGroupNodes(nodes)

    case createGroupNode.FullCommand():
        if err := client.New(*server, mustLogin(*server)).CreateGroupNode(*createGroupNodeGroup, *createGroupNodeNode); err != nil {
            log.Fatalln(err)
        }

    case deleteGroupNode.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteGroupNode(*deleteGroupNodeGroup, *deleteGroupNodeNode); err != nil {
            log.Fatalln(err)
        }

    case deployInit.FullCommand():
        if err := client.New(*server, mustLogin(*server)).CreateDeployment(*deployInitGroup, *deployInitRepo); err != nil {
            log.Fatalln(err)
        }

    case deployPrepare.FullCommand():
        if err := client.New(*server, mustLogin(*server)).PrepareDeployment(*deployPrepareGroup); err != nil {
            log.Fatalln(err)
        }

    case deployExecute.FullCommand():
        if err := client.New(*server, mustLogin(*server)).ExecuteDeployment(*deployExecuteGroup); err != nil {
            log.Fatalln(err)
        }

    case deployClear.FullCommand():
        if err := client.New(*server, mustLogin(*server)).ClearDeployment(*deployClearGroup); err != nil {
            log.Fatalln(err)
        }

    case deployDelete.FullCommand():
        if err := client.New(*server, mustLogin(*server)).DeleteDeployment(*deployDeleteGroup); err != nil {
            log.Fatalln(err)
        }

    default:
        kingpin.Fatalf("command not specified, try --help")
    }
}
