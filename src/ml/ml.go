package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path"
//    "strings"

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
    server = app.Flag("server", "Server address.").Default("localhost:8000").OverrideDefaultFromEnvar("MOONLEGEND_SERVER").String()

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

    default:
        kingpin.Fatalf("command not specified, try --help")
    }
}
