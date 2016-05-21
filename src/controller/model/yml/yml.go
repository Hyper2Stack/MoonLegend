package yml

import (
    "strings"
)

const (
    RestartAlways         = "always"
    RestartNo             = "no"
    FixedPortMapping      = "fix"
    RandomPortMapping     = "random"
    CustomizedPortMapping = "customized"
)

type Yml struct {
    Infra    *Infra              `yaml:"infrastructure"`
    Services map[string]*Service `yaml:"services"`
    Runtime  *Runtime            `yaml:"runtime"`
}

type Infra struct {
    Network []string `yaml:"network"`
}

type Service struct {
    Image       string   `yaml:"image"`
    Singleton   bool     `yaml:"singleton"`
    Ports       []string `yaml:"ports"`
    Config      *Config  `yaml:"config_file"`
    Networks    []string `yaml:"networks"`
    Volumes     []string `yaml:"volumes"`
    Depends     []string `yaml:"depends_on"`
    Environment []string `yaml:"environment"`
}

type Config struct {
    Path    string `yaml:"path" json:"path"`
    Mode    string `yaml:"mode" json:"mode"`
    Content string `yaml:"content" json:"content"`
}

type Runtime struct {
    Environment   []string                  `yaml:"env" json:"env"`
    GlobalPolicy  *GlobalPolicy             `yaml:"global_policy" json:"global_policy"`
    ServicePolicy map[string]*ServicePolicy `yaml:"services" json:"services"`
}

type GlobalPolicy struct {
    Restart     string `yaml:"restart" json:"restart"`
    PortMapping string `yaml:"port_mapping" json:"port_mapping"`
}

type ServicePolicy struct {
    Restart     string `yaml:"restart" json:"restart"`
    InstanceNum int    `yaml:"instance_num" json:"instance_num"`
    PortMapping string `yaml:"port_mapping" json:"port_mapping"`
    PortRange   string `yaml:"port_range" json:"port_range"`
}

func (r *Runtime) Env(name string) string {
    for _, s := range r.Environment {
        ss := strings.SplitN(s, "=", 2)
        if ss[0] == name && len(ss) > 1 {
            return ss[1]
        }
    }

    return ""
}

func (r *Runtime) GetPolicy(service string) *ServicePolicy {
    sp := r.ServicePolicy[service]
    if sp == nil {
        sp = new(ServicePolicy)
    }

    if sp.Restart == "" {
        if r.GlobalPolicy != nil && r.GlobalPolicy.Restart != "" {
            sp.Restart = r.GlobalPolicy.Restart
        } else {
            sp.Restart = RestartNo
        }
    }

    if sp.InstanceNum == 0 {
        sp.InstanceNum = 1
    }

    if sp.PortMapping == "" {
        if r.GlobalPolicy != nil && r.GlobalPolicy.PortMapping != "" {
            sp.PortMapping = r.GlobalPolicy.PortMapping
        } else {
            sp.PortMapping = FixedPortMapping
        }
    }

    return sp
}
