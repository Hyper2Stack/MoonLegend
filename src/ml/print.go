package main

import (
    "fmt"

    "controller/model"
)

func printProfile(profile *model.User) {
    fmt.Printf("--- profile ---\n")
    fmt.Printf("name:\t%s\n", profile.Name)
    fmt.Printf("email:\t%s\n", profile.Email)
    fmt.Printf("key:\t%s\n", profile.Key)
}

func printRepos(repos []*model.Repo) {
    for _, repo := range repos {
        fmt.Println(repo.Name)
    }
}

func printRepoTags(tags []*model.RepoTag) {
    for _, tag := range tags {
        fmt.Println(tag.Name)
    }
}

func printRepoTag(tag *model.RepoTag) {
    fmt.Println(tag.Yml)
}

func printNodes(nodes []*model.Node) {
    for _, node := range nodes {
        fmt.Println(node.Name)
    }
}

func printNode(node *model.Node) {
    fmt.Printf("--- node details ---\n")
    fmt.Printf("name:\t%s\n", node.Name)
    fmt.Printf("uuid:\t%s\n", node.Uuid)
    fmt.Printf("status:\t%s\n", node.Status)
    fmt.Printf("tags:\t%v\n", node.Tags)
    fmt.Println("networks:")
    for _, nic := range node.Nics {
        fmt.Printf("  %s:\n", nic.Name)
        fmt.Printf("    address:\t%s\n", nic.Ip4Addr)
        fmt.Printf("    tags:\t%v\n", nic.Tags)
    }
}

type Node struct {
    Tags        []string `json:"tags"`
    Nics        []*Nic   `json:"nics"`
}

type Nic struct {
    Name    string   `json:"name"`
    Ip4Addr string   `json:"ip4addr"`
    Tags    []string `json:"tags"`
}

func printNodeTags(node *model.Node) {
    for _, tag := range node.Tags {
        fmt.Println(tag)
    }
}

func findNic(node *model.Node, name string) *model.Nic {
    for _, nic := range node.Nics {
        if nic.Name == name {
            return nic
        }
    }

    return nil
}

func printNicTags(nic *model.Nic) {
    for _, tag := range nic.Tags {
        fmt.Println(tag)
    }
}
