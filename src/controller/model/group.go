package model

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strings"
    "text/template"

    "gopkg.in/yaml.v2"
    "controller/model/yml"

    "github.com/hyper2stack/mooncommon/protocol"
)

type Group struct {
    Id              int64                          `json:"id"`
    Name            string                         `json:"name"`
    Description     string                         `json:"description"`
    OwnerId         int64                          `json:"owner_id"`
    Status          string                         `json:"status"`
    Deployment      *Deployment                    `json:"deployment"`
    Process         map[string][]*InstanceStatus   `json:"process"`
}

type InstanceStatus struct {
    NodeUuid string `json:"node_uuid"`
    Service  string `json:"service"`
    Name     string `json:"name"`
    Status   string `json:"status"`
}

type Deployment struct {
    Repo         string       `json:"repo"`
    Runtime      *yml.Runtime `json:"runtime"`
    InstanceList []*Instance  `json:"instance_list"`
    Yml          *yml.Yml     `json:"yml"`
}

type Service struct {
    Name     string   `json:"name"`
    Image    string   `json:"image"`
    Networks []string `json:"networks"`
    Depends  []string `json:"depends_on"`
}

type Instance struct {
    Host            *Node            `json:"host"`
    Service         *Service         `json:"service"`
    Name            string           `json:"name"`
    PrepareFile     *protocol.File   `json:"prepare_file"`
    PrepareCommand  *protocol.Script `json:"prepare_commands"`
    RunCommand      *protocol.Script `json:"run_command"`
    RestartCommand  *protocol.Script `json:"restart_command"`
    RmCommand       *protocol.Script `json:"rm_command"`
    Entrypoints     []*Entrypoint    `json:"entrypoints"`
    Config          *yml.Config      `json:"config"`
    Env             []string         `json:"env"`
}

type Entrypoint struct {
    Protocol      string `json:"protocol"`
    ListeningAddr string `json:"listening_addr"`
    ListeningPort string `json:"listening_port"`
    ContainerPort string `json:"container_port"`
}

const (
    StatusRaw            = "raw"
    StatusCreated        = "created"
    StatusPreparing      = "preparing"
    StatusPrepared       = "prepared"
    StatusPrepareTimeout = "prepare_timeout"
    StatusPrepareError   = "prepare_error"
    StatusDeploying      = "deploying"
    StatusDeployed       = "deployed"
    StatusDeployTimeout  = "deploy_timeout"
    StatusDeployError    = "deploy_error"
    StatusClearing       = "clearing"
)

var dataDir = "/var/run/moon"

////////////////////////////////////////////////////////////////////////////////

func ListGroup(cs []*Condition, o *Order, p *Paging) []*Group {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, name, description, owner_id, status, deployment, process
        FROM
            nodeGroup
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Group, 0)
    for rows.Next() {
        g := new(Group)
        var deployInfo string
        var processInfo string
        if err := rows.Scan(
            &g.Id, &g.Name, &g.Description, &g.OwnerId, &g.Status, &deployInfo, &processInfo,
        ); err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(deployInfo), &g.Deployment); err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(processInfo), &g.Process); err != nil {
            panic(err)
        }

        l = append(l, g)
    }

    return l
}

func GetGroupById(id int64) *Group {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListGroup(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetGroupByNameAndOwnerId(name string, ownerId int64) *Group {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("name", "=", name))
    conditions = append(conditions, NewCondition("owner_id", "=", ownerId))

    l := ListGroup(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (g *Group) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO nodeGroup(
            name, description, owner_id, status, deployment, process
        )
        VALUES(?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    deployInfo, err := json.Marshal(g.Deployment)
    if err != nil {
        panic(err)
    }

    processInfo, err := json.Marshal(g.Process)
    if err != nil {
        panic(err)
    }

    result, err := stmt.Exec(
        g.Name, g.Description, g.OwnerId, g.Status, string(deployInfo), string(processInfo),
    )
    if err != nil {
        panic(err)
    }

    g.Id, err = result.LastInsertId()
    if err != nil {
        panic(err)
    }
}

func (g *Group) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            nodeGroup
        SET
            name = ?,
            description = ?,
            status = ?,
            deployment = ?,
            process = ?
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    deployInfo, err := json.Marshal(g.Deployment)
    if err != nil {
        panic(err)
    }

    processInfo, err := json.Marshal(g.Process)
    if err != nil {
        panic(err)
    }

    if _, err := stmt.Exec(
        g.Name, g.Description, g.Status, string(deployInfo), string(processInfo), g.Id,
    ); err != nil {
        panic(err)
    }
}

func (g *Group) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            nodeGroup
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(g.Id); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (g *Group) Nodes() []*Node {
    rows, err := db.Query(`
        SELECT
            node.id, node.uuid, node.name, node.description, node.status, node.owner_id, node.tags, node.nics
        FROM
            node, nodeMembership
        WHERE
            nodeMembership.group_id =  ? and node.id = nodeMembership.node_id
        `, g.Id,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Node, 0)
    for rows.Next() {
        n := new(Node)
        var tagInfo string
        var nicInfo string
        if err := rows.Scan(
            &n.Id, &n.Uuid, &n.Name, &n.Description, &n.Status, &n.OwnerId, &tagInfo, &nicInfo,
        ); err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(tagInfo), &n.Tags); err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(nicInfo), &n.Nics); err != nil {
            panic(err)
        }

        l = append(l, n)
    }

    return l
}

func (g *Group) HasNode(node *Node) bool {
    for _, n := range g.Nodes() {
        if n.Uuid == node.Uuid {
            return true
        }
    }

    return false
}

func (g *Group) AddNode(node *Node) {
    stmt, err := db.Prepare(`
        INSERT INTO nodeMembership(
            group_id, node_id
        )
        VALUES(?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        g.Id, node.Id,
    )
    if err != nil {
        panic(err)
    }
}

func (g *Group) DeleteNode(node *Node) {
    stmt, err := db.Prepare(`
        DELETE FROM
            nodeMembership
        WHERE
            group_id = ? and node_id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(g.Id, node.Id); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (g *Group) InitDeployment(repo *Repo, tag *RepoTag, runtime *yml.Runtime) error {
    if g.Status != StatusRaw {
        return fmt.Errorf("invalid group status")
    }

    yml := new(yml.Yml)
    if err := yaml.Unmarshal([]byte(tag.Yml), yml); err != nil {
        return fmt.Errorf("parse yml error, %v", err)
    }

    // TBD, check physical infra available, un-necessary?

    d := new(Deployment)
    d.Yml = yml
    d.Repo = CompileRepoString(repo, tag)
    d.InstanceList = make([]*Instance, 0)
    if runtime == nil {
        d.Runtime = yml.Runtime
    } else {
        d.Runtime = runtime
    }

    // generate instances
    for k, v := range yml.Services {
        ds := new(Service)
        ds.Name = k
        ds.Image = v.Image
        ds.Networks = v.Networks
        ds.Depends = v.Depends

        num := 1
        if n := d.Runtime.GetPolicy(k).InstanceNum; n > 1 && !v.Singleton {
            num = n
        }

        instances := make([]*Instance, 0)
        for i := 0; i < num; i++ {
            ins := new(Instance)
            ins.Service = ds
            ins.Name = fmt.Sprintf("%s_%d", ds.Name, i)
            ins.Config = v.Config
            ins.Env = v.Environment
            ins.Entrypoints = make([]*Entrypoint, 0)
            for _, port := range v.Ports {
                ep := new(Entrypoint)
                ss := strings.Split(port, "/")
                if len(ss) < 2 {
                    ep.Protocol = "tcp"
                    ep.ContainerPort = port
                } else {
                    ep.Protocol = ss[1]
                    ep.ContainerPort = ss[0]
                }
                ins.Entrypoints = append(ins.Entrypoints, ep)
            }

            instances = append(instances, ins)
        }

        if err := g.mappingNodeAndPort(instances, ds, d.Runtime); err != nil {
            return fmt.Errorf("mapping node and port error, %v", err)
        }

        d.InstanceList = append(d.InstanceList, instances...)
    }

    // render
    if err := d.Render(); err != nil {
        return err
    }

    // init deploy commands
    d.InitDeployCmd()

    // update database
    g.Deployment = d
    g.Status = StatusCreated
    g.Update()

    return nil
}

func (g *Group) InitDeployProcess() {
    g.Process = make(map[string][]*InstanceStatus)
    for _, i := range g.Deployment.InstanceList {
        is := new(InstanceStatus)
        is.NodeUuid = i.Host.Uuid
        is.Status = StatusCreated
        is.Service = i.Service.Name
        is.Name = i.Name
        g.Process[i.Host.Uuid] = append(g.Process[i.Host.Uuid], is)
    }
}

func (g *Group) FindNodeByService(s *Service) []*Node {
    result := make([]*Node, 0)
    for _, node := range g.Nodes() {
        if node.HasTag(s.Name) && node.HasNicTags(s.Networks) {
            result = append(result, node)
        }
    }

    return result
}

func (g *Group) ParsePrepareProcessStatus() (status string, done bool) {
    done = true
    status = StatusPrepared

    for _, insts := range g.Process {
        for _, is := range insts {
            if is.Status == StatusPrepared {
                continue
            }
            if is.Status == StatusPrepareTimeout {
                if status == StatusPrepared {
                    status = StatusPrepareTimeout
                }
                continue
            }

            if is.Status == StatusPrepareError {
                status = StatusPrepareError
                continue
            }

            done = false
            return
        }
    }

    return
}

func (g *Group) GetProcessOfServiceMap() map[string][]*InstanceStatus {
    result := make(map[string][]*InstanceStatus)
    for _, isList := range g.Process {
        for _, is := range isList {
            result[is.Service] = append(result[is.Service], is)
        }
    }

    return result
}

func (g *Group) mappingNodeAndPort(instances []*Instance, service *Service, runtime *yml.Runtime) error {
    nodes := g.FindNodeByService(service)
    policy := runtime.GetPolicy(service.Name)

    switch policy.PortMapping {
    case yml.FixedPortMapping:
        return fixedMapping(instances, nodes)
    case yml.RandomPortMapping:
        return randomMapping(instances, nodes)
    case yml.CustomizedPortMapping:
        return customMapping(instances, nodes, policy.PortRange)
    default:
        return fmt.Errorf("invalid port mapping policy, %s", policy.PortMapping)
    }

    return nil
}

func fixedMapping(instances []*Instance, nodes []*Node) error {
    if len(nodes) < len(instances) {
        return fmt.Errorf("fixed port mapping, node num < instance num, %d < %d", len(nodes), len(instances))
    }

    for i := 0; i < len(instances); i++ {
        instances[i].Host = nodes[i]
        for _, ep := range instances[i].Entrypoints {
            ep.ListeningAddr = "0.0.0.0"
            ep.ListeningPort = ep.ContainerPort
        }
    }

    return nil
}

func randomMapping(instances []*Instance, nodes []*Node) error {
    // TBD
    return nil
}

func customMapping(instances []*Instance, nodes []*Node, portRange string) error {
    // TBD
    return nil
}

func (d *Deployment) FindInstanceByName(name string) *Instance {
    for _, i := range d.InstanceList {
        if i.Name == name {
            return i
        }
    }

    return nil
}

func (d *Deployment) TopSortServices() ([]string, error) {
    edges := make([]*Edge, 0)
    nodes := make([]string, 0)

    for _, i := range d.InstanceList {
        if containElement(nodes, i.Service.Name) {
            continue
        }

        nodes = append(nodes, i.Service.Name)
        for _, parent := range i.Service.Depends {
            edge := new(Edge)
            edge.From = parent
            edge.To = i.Service.Name
            edges = append(edges, edge)
        }
    }

    for _, edge := range edges {
        if !containElement(nodes, edge.From) {
            return nil, fmt.Errorf("invalid edge %s -> %s, missing node %s", edge.From, edge.To, edge.From)
        }
    }

    return TopologicSort(nodes, edges)
}

func (d *Deployment) Render() error {
    for _, ins := range d.InstanceList {
        // Config & Env will be rendered, so every instance should have a copy of these two fields
        if ins.Config != nil {
            c := *ins.Config
            var err error
            if c.Content, err = d.render(ins.Config.Content); err != nil {
                return err
            }
            ins.Config = &c
        }

        env := make([]string, 0)
        for _, e := range ins.Env {
            s, err := d.render(e)
            if err != nil {
                return err
            }
            env = append(env, s)
        }
        ins.Env = env
    }

    return nil
}

func (d *Deployment) render(s string) (string, error) {
    var b bytes.Buffer
    t, err := template.New("").Parse(s)
    if err != nil {
        return "", err
    }

    if err := t.Execute(&b, d); err != nil {
        return "", err
    }

    return b.String(), nil
}

func (d *Deployment) InitDeployCmd() {
    for _, ins := range d.InstanceList {
        ins.initPrepareFile()
        ins.initPrepareCmd()
        ins.initRunCmd(d.Runtime)
        ins.initRmCmd()
        ins.initRestartCmd()
    }
}

func (ins *Instance) initPrepareFile() {
    if ins.Config == nil {
        return
    }

    ins.PrepareFile = &protocol.File{
        Path: fmt.Sprintf("%s/%s", dataDir, ins.Name),
        Mode: ins.Config.Mode,
        Content: ins.Config.Content,
    }
}

func (ins *Instance) initPrepareCmd() {
    ins.PrepareCommand = new(protocol.Script)
    command := new(protocol.Command)
    command.Command = "bash"
    command.Args = []string{"-c", fmt.Sprintf("docker inspect %s > /dev/null 2>&1 || docker pull %s", ins.Service.Image, ins.Service.Image)}
    command.Restrict = true
    ins.PrepareCommand.Commands = append(ins.PrepareCommand.Commands, command)
}

func (ins *Instance) initRunCmd(runtime *yml.Runtime) {
    ins.RunCommand = new(protocol.Script)
    command := new(protocol.Command)
    command.Command = "docker"
    command.Restrict = true
    command.Args = []string{"run", "-d"}
    command.Args = append(command.Args, fmt.Sprintf("--name=%s", ins.Name))

    for _, envStr := range ins.Env {
        command.Args = append(command.Args, "-e")
        command.Args = append(command.Args, envStr)
    }

    if ins.Config != nil {
        command.Args = append(command.Args, "-v")
        command.Args = append(command.Args, fmt.Sprintf("%s/%s:%s", dataDir, ins.Name, ins.Config.Path))
    }

    for _, ep := range ins.Entrypoints {
        command.Args = append(command.Args, "-p")
        command.Args = append(
            command.Args,
            fmt.Sprintf("%s:%s:%s/%s", ep.ListeningAddr, ep.ListeningPort, ep.ContainerPort, ep.Protocol),
        )
    }

    command.Args = append(
        command.Args,
        fmt.Sprintf("--restart=%s", runtime.GetPolicy(ins.Service.Name).Restart),
    )

    command.Args = append(command.Args, ins.Service.Image)
    ins.RunCommand.Commands = append(ins.RunCommand.Commands, command)
}

func (ins *Instance) initRmCmd() {
    ins.RmCommand = new(protocol.Script)
    command := new(protocol.Command)
    command.Command = "docker"
    command.Restrict = false
    command.Args = []string{"rm", "-vf", ins.Name}
    ins.RmCommand.Commands = append(ins.RmCommand.Commands, command)
}

func (ins *Instance) initRestartCmd() {
    ins.RestartCommand = new(protocol.Script)
    command := new(protocol.Command)
    command.Command = "docker"
    command.Restrict = true
    command.Args = []string{"restart", ins.Name}
    ins.RestartCommand.Commands = append(ins.RmCommand.Commands, command)
}

////////////////////////////////////////////////////////////////////////////////

func (d *Deployment) Instances(service string) []*Instance {
    result := make([]*Instance, 0)
    for _, i := range d.InstanceList {
        if i.Service != nil && i.Service.Name == service {
            result = append(result, i)
        }
    }

    return result
}

func (d *Deployment) Singleton(service string) *Instance {
    result := d.Instances(service)
    if len(result) != 1 {
        return nil
    }

    return result[0]
}

func (i *Instance) AddressOf(network string) string {
    for _, nic := range i.Host.Nics {
        for _, tag := range nic.Tags {
            if tag == network {
                return nic.Ip4Addr
            }
        }
    }

    return ""
}

func (i *Instance) Address() string {
    return i.AddressOf(i.Service.Networks[0])
}

func (i *Instance) PortOf(port string) string {
    ss := strings.Split(port , "/")
    for _, ep := range i.Entrypoints {
        if ss[0] != ep.ContainerPort {
            continue
        }

        if len(ss) == 1 || ss[1] == ep.Protocol {
            return ep.ListeningPort
        }
    }

    return ""
}

func (i *Instance) Port() string {
    return i.Entrypoints[0].ListeningPort
}
