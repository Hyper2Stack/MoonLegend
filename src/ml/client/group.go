package client

import (
    "encoding/json"
    "fmt"

    "controller/model"
)

func (c *Client) Groups() ([]*model.Group, error) {
    status, outbody, _, err := c.do("GET", "/api/v1/user/groups", nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    var groups []*model.Group
    if err := json.Unmarshal(outbody, &groups); err != nil {
        return nil, err
    }

    return groups, nil
}

func (c *Client) CreateGroup(name string) error {
    in := struct {
        Name        string `json:"name"`
    }{}

    in.Name = name

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", "/api/v1/user/groups", nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) Group(name string) (*model.Group, error) {
    status, outbody, _, err := c.do("GET", fmt.Sprintf("/api/v1/user/groups/%s", name), nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    group := new(model.Group)
    if err := json.Unmarshal(outbody, group); err != nil {
        return nil, err
    }

    return group, nil
}

func (c *Client) DeleteGroup(name string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/groups/%s", name), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) GroupNodes(group string) ([]*model.Node, error) {
    status, outbody, _, err := c.do("GET", fmt.Sprintf("/api/v1/user/groups/%s/nodes", group), nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    var nodes []*model.Node
    if err := json.Unmarshal(outbody, &nodes); err != nil {
        return nil, err
    }

    return nodes, nil
}

func (c *Client) CreateGroupNode(group, node string) error {
    in := struct {
        Name string `json:"name"`
    }{}

    in.Name = node

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", fmt.Sprintf("/api/v1/user/groups/%s/nodes", group), nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) GroupNode(group, node string) (*model.Node, error) {
    status, outbody, _, err := c.do("GET", fmt.Sprintf("/api/v1/user/groups/%s/nodes/%s", group, node), nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    result := new(model.Node)
    if err := json.Unmarshal(outbody, result); err != nil {
        return nil, err
    }

    return result, nil
}

func (c *Client) DeleteGroupNode(group, node string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/groups/%s/nodes/%s", group, node), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}
