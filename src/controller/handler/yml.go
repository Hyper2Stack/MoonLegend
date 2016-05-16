package handler

import (
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
    Networks  []string `yaml:"network"`
    Volumes   []string `yaml:"volumes"`
    Depends   []string `yaml:"depends_on"`
}

type Config struct {
    Path    string `yaml:"path"`
    Mode    string `yaml:"mode"`
    Content string `yaml:"content"`
}

type Runtime struct {
    Env      []string                  `yaml:"env"`
    Policy   *GlobalPolicy             `yaml:"global_policy"`
    Services map[string]*ServicePolicy `yaml:"services"`
}

type GlobalPolicy struct {
    Restart     bool   `yaml:"restart"`
    PortMapping string `yaml:"port_mapping"`
}

type ServicePolicy struct {
    Instance    int    `yaml:"instance"`
    PortMapping string `yaml:"port_mapping"`
    PortRange   string `yaml:"port_range"`
}
