package client

import (
    "encoding/json"
    "fmt"

    "controller/model"
)

func (c *Client) Signup(user, password string) error {
    in := struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }{}

    in.Username = user
    in.Password = password

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", "/api/v1/signup", nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) Login(user, password string) (string, error) {
    in := struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }{}

    in.Username = user
    in.Password = password

    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("POST", "/api/v1/login", nil, inbody)
    if err != nil {
        return "", err
    }
    if status/100 != 2 {
        return "", fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    out := struct {
        Key string `json:"key"`
    }{}
    if err := json.Unmarshal(outbody, &out); err != nil {
        return "", err
    }

    return out.Key, nil
}

func (c *Client) Profile() (*model.User, error) {
    status, outbody, _, err := c.do("GET", "/api/v1/user", nil, nil)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    profile := new(model.User)
    if err := json.Unmarshal(outbody, profile); err != nil {
        return nil, err
    }

    return profile, nil
}

func (c *Client) ResetPassword(password string) error {
    in := struct {
        Password string `json:"password"`
    }{}

    in.Password = password
    inbody, _ := json.Marshal(in)
    status, outbody, _, err := c.do("PUT", "/api/v1/user/reset-password", nil, inbody)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}

func (c *Client) ResetKey() error {
    status, outbody, _, err := c.do("PUT", "/api/v1/user/reset-key", nil, nil)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("status code %d, %s", status, string(outbody))
    }

    return nil
}
