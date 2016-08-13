package client

import (
    "encoding/json"
    "fmt"

    "controller/model"
)

func (c *Client) Repos() ([]*model.Repo, error) {
    status, outbody, _, err := c.do("GET", "/api/v1/user/repos", nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    var repos []*model.Repo
    if err := json.Unmarshal(outbody, &repos); err != nil {
        return nil, err
    }

    return repos, nil
}

func (c *Client) CreateRepo(name string, isPublic bool) error {
    in := struct {
        Name        string `json:"name"`
        IsPublic    bool   `json:"is_public"`
        Description string `json:"description"`
        Readme      string `json:"readme"`
    }{}

    in.Name = name
    in.IsPublic = isPublic

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", "/api/v1/user/repos", nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil    
}

func (c *Client) Repo(name string) (*model.Repo, error) {
    status, outbody, _, err := c.do("GET", fmt.Sprintf("/api/v1/user/repos/%s", name), nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    repo := new(model.Repo)
    if err := json.Unmarshal(outbody, repo); err != nil {
        return nil, err
    }

    return repo, nil
}

func (c *Client) DeleteRepo(name string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/repos/%s", name), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) RepoTags(repo string) ([]*model.RepoTag, error) {
    status, outbody, _, err := c.do("GET", fmt.Sprintf("/api/v1/user/repos/%s/tags", repo), nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    var tags []*model.RepoTag
    if err := json.Unmarshal(outbody, &tags); err != nil {
        return nil, err
    }

    return tags, nil
}

func (c *Client) CreateRepoTag(repo, tag, yml string) error {
    in := struct {
        Name string `json:"name"`
        Yml  string `json:"yml"`
    }{}

    in.Name = tag
    in.Yml = yml

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", fmt.Sprintf("/api/v1/user/repos/%s/tags", repo), nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) RepoTag(repo, tag string) (*model.RepoTag, error) {
    status, outbody, _, err := c.do("GET", fmt.Sprintf("/api/v1/user/repos/%s/tags/%s", repo, tag), nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    repoTag := new(model.RepoTag)
    if err := json.Unmarshal(outbody, repoTag); err != nil {
        return nil, err
    }

    return repoTag, nil
}

func (c *Client) DeleteRepoTag(repo, tag string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/repos/%s/tags/%s", repo, tag), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}
