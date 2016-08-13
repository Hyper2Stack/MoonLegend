package main

import (
    "fmt"

    "controller/model"
)

func printProfile(profile *model.User) {
    fmt.Printf("--- profile ---\n")
    fmt.Printf("name:  %s\n", profile.Name)
    fmt.Printf("email: %s\n", profile.Email)
    fmt.Printf("key:   %s\n", profile.Key)
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
