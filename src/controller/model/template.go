package model

import (
    "strings"
)

type Runtime struct {
    Environment map[string]string
}

func (r *Runtime) Env(name string) string {
    return r.Environment[name]
}

func (d *Deployment) Service(name string) *Service {
    for k, v := range d.Services {
        if k == name {
            return v
        }
    }

    return nil
}

func (s *Service) Instances() []*Instance {
    return s.InstanceList
}

func (s *Service) Singoleton() *Instance {
    if len(s.InstanceList) != 1 {
        return nil
    }

    return s.InstanceList[0]
}

func (i *Instance) Node() *Node {
    return i.Host
}

func (i *Instance) Entrypoint(port string) *Entrypoint {
    ss := strings.Split(port , "/")
    for _, ep := range i.Entrypoints {
        if ss[0] != ep.ContainerPort {
            continue
        }

        if len(ss) == 1 || ss[1] == ep.Protocol {
            return ep
        }
    }

    return nil
}

func (ep *Entrypoint) HostPort() string  {
    return ep.ListeningPort
}

func (node *Node) Network(name string) *Nic {
    for _, network := range node.Nics {
        for _, tag := range network.Tags {
            if tag == name {
                return network
            }
        }
    }

    return nil
}

func (nic *Nic) Address() string {
    return nic.Ip4Addr
}
