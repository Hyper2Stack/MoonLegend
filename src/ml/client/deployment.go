package client

import (
    "encoding/json"
    "fmt"
)

func (c *Client) CreateDeployment(group, repo string) error {
    in := struct {
        Repo string `json:"repo"`
    }{}

    in.Repo = repo

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", fmt.Sprintf("/api/v1/user/groups/%s/deployment", group), nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) DeleteDeployment(group string) error {
    status, outbody, _, err := c.do("DELETE", fmt.Sprintf("/api/v1/user/groups/%s/deployment", group), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) PrepareDeployment(group string) error {
    status, outbody, _, err := c.do("PUT", fmt.Sprintf("/api/v1/user/groups/%s/deployment/prepare", group), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) ExecuteDeployment(group string) error {
    status, outbody, _, err := c.do("PUT", fmt.Sprintf("/api/v1/user/groups/%s/deployment/execute", group), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) ClearDeployment(group string) error {
    status, outbody, _, err := c.do("PUT", fmt.Sprintf("/api/v1/user/groups/%s/deployment/clear", group), nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}
