package yml

import (
    "strings"
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
    Image     string   `yaml:"image"`
    Singleton bool     `yaml:"singleton"`
    Ports     []string `yaml:"ports"`
    Config    *Config  `yaml:"config_file"`
    Networks  []string `yaml:"networks"`
    Volumes   []string `yaml:"volumes"`
    Depends   []string `yaml:"depends_on"`
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
    Restart     bool   `yaml:"restart" json:"restart"`
    PortMapping string `yaml:"port_mapping" json:"port_mapping"`
}

type ServicePolicy struct {
    Instance    int    `yaml:"instance" json:"instance"`
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
