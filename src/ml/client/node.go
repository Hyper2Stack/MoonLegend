package client

import (
    "encoding/json"
    "fmt"

    "controller/model"
)

func (c *Client) Nodes() ([]*model.Node, error) {
    status, outbody, _, err := c.do("GET", "/api/v1/user/nodes", nil, nil)
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

func (c *Client) Node(name string) (*model.Node, error) {
    status, outbody, _, err := c.do("GET", fmt.Sprintf("/api/v1/user/nodes/%s", name), nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    node := new(model.Node)
    if err := json.Unmarshal(outbody, node); err != nil {
        return nil, err
    }

    return node, nil
}

func (c *Client) DeleteNode(name string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/nodes/%s", name), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) CreateNodeTag(node, tag string) error {
    in := struct {
        Name string `json:"name"`
    }{}

    in.Name = tag

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", fmt.Sprintf("/api/v1/user/nodes/%s/tags", node), nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) DeleteNodeTag(node, tag string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/nodes/%s/tags/%s", node, tag), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) CreateNicTag(node, nic, tag string) error {
    in := struct {
        Name string `json:"name"`
    }{}

    in.Name = tag

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", fmt.Sprintf("/api/v1/user/nodes/%s/nics/%s/tags", node, nic), nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) DeleteNicTag(node, nic, tag string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/nodes/%s/nics/%s/tags/%s", node, nic, tag), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}
