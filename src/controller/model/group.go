package model

import (
)

type Group struct {
    Id          int64       `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Owner       int64       `json:"owner"`
    Status      string      `json:"status"`
    Nodes       []*Node     `json:"nodes"`
    Deployment  *Deployment `json:"deployment"`
}

type Deployment struct {
    RepoName   string     `json:"name"`
    Services   []*Service `json:"services"`
}

type Service struct {
    Name       string      `json:"name"`
    Instances  []*Instance `json:"instance"`
}

type Instance struct {
    Node        *Node         `json:"node"`
    Container   *Container    `json:"container"`
    Entrypoints []*Entrypoint `json:"entrypoints"`
}

type Entrypoint struct {
    Protocol      string `json:"protocol"`
    ListeningAddr string `json:"listeningaddr"`
    ListeningPort string `json:"listeningport"`
    ContainerPort string `json:"containerport"`
}

type Container struct {
    Name           string `json:"name"`
    StartCommand   string `json:"startcommand"`
    DestroyCommand string `json:"destroycommand"`
}

const (
    StatusRaw          = "raw"
    StatusCreated      = "created"
    StatusPreparing    = "preparing"
    StatusPrepared     = "prepared"
    StatusDeploying    = "deploying"
    StatusDeployed     = "deployed"
)

////////////////////////////////////////////////////////////////////////////////

func ListGroup(cs []*Condition, o *Order, p *Paging) []*Group {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, name, description, owner, status, repoId
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
        var repoId int64
        g := new(Group)
        if err := rows.Scan(
            &g.Id, &g.Name, &g.Description, &g.Owner, &g.Status, &repoId,
        ); err != nil {
            panic(err)
        }
        
        if repoId != 0 {
            d := new(Deployment)
            d.RepoName = GetRepoById(repoId).Name
            d.Services = GetServicesByGroupId(g.Id)
            g.Deployment = d
        }

        g.Nodes = GetNodesByGroupId(g.Id)

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

func GetGroupByNameAndOwner(name string, owner int64) *Group {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("name", "=", name))
    conditions = append(conditions, NewCondition("owner", "=", owner))

    l := ListGroup(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (g *Group) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO nodeGroup(
            name, description, owner, status
        )
        VALUES(?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(
        g.Name, g.Description, g.Owner, g.Status,
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
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        g.Name, g.Description, g.Id,
    ); err != nil {
        panic(err)
    }
}

func (g *Group) Delete() {
    if (g.Nodes != nil) {
        for _, n := range g.Nodes {
            // delete from nodeTable
            n.Delete()
            // delete from nodeMembershipTable
            g.DeleteMembershipAll()
        }
    }
    
    if (g.Deployment != nil) {
        g.DeleteDeployment()
    }

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

func (g *Group) AddMembership(nodeId int64) {
    stmt, err := db.Prepare(`
        INSERT INTO nodeMembership(
            groupId, nodeId
        )
        VALUES(?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        g.Id, nodeId,
    )
    if err != nil {
        panic(err)
    }
}

func (g *Group) DeleteMembershipAll() {
    stmt, err := db.Prepare(`
        DELETE FROM
            nodeMembership
        WHERE
            groupId = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(g.Id); err != nil {
        panic(err)
    }
}

func (g *Group) DeleteMembership(nodeId int64) {
    stmt, err := db.Prepare(`
        DELETE FROM
            nodeMembership
        WHERE
            groupId = ? and nodeId = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(g.Id, nodeId); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (g *Group) DeleteDeployment() {
    // delete from instance table
    g.DeleteInstance()
    // delete from service table
    g.DeleteService()
    stmt, err := db.Prepare(`
        UPDATE
            nodeGroup
        SET
            repoId = 0,
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        g.Id,
    ); err != nil {
        panic(err)
    }
}

func (g *Group) DeleteService() {
    stmt, err := db.Prepare(`
        DELETE FROM
            service
        WHERE
            groupId = 
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(g.Id); err != nil {
        panic(err)
    }
}

func (g *Group) DeleteInstance() {
    stmt, err := db.Prepare(`
        DELETE FROM
            instance
        WHERE
            serviceId in (SELECT id from service where groupId = ?)
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

func (g *Group) Execute(repoId int64) {
    for _, s := range g.Deployment.Services {
        s.Execute(g.Id)
    }
    stmt, err := db.Prepare(`
        UPDATE
            nodeGroup
        SET
            repoId = ?,
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        repoId, g.Id,
    ); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (g *Group) ParseYml(path string) error {
    // TBD
    return nil
}

////////////////////////////////////////////////////////////////////////////////

func GetServicesByGroupId(groupId int64) []*Service {
    rows, err := db.Query(`
        SELECT
            id, name, configJson
        FROM
            service
        WHERE
            groupId = ?
        `, groupId,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Service, 0)
    for rows.Next() {
        var serviceId int64
        // TODO: how to convert configJson
        var configJson string
        s := new(Service)
        if err := rows.Scan(
            &serviceId, &s.Name, &configJson,
        ); err != nil {
            panic(err)
        }
        
        s.Instances = GetInstancesByServiceId(serviceId)

        l = append(l, s)
    }

    return l
}

func (s *Service) Execute(groupId int64) {
// TODO: how to convert to configJson
    stmt, err := db.Prepare(`
        INSERT INTO service(
            name, groupId, configJson
        )
        VALUES(?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(
        s.Name, groupId, "",
    )
    if err != nil {
        panic(err)
    }

    serviceId, err := result.LastInsertId()
    if err != nil {
        panic(err)
    }

    for _, i := range s.Instances {
        i.Execute(serviceId)
    }
}

////////////////////////////////////////////////////////////////////////////////

func GetInstancesByServiceId(serviceId int64) []*Instance {
    rows, err := db.Query(`
        SELECT
            name, nodeId, configJson
        FROM
            instance
        WHERE
            serviceId = ?
        `, serviceId,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Instance, 0)
    for rows.Next() {
        var nodeId int64
        // TODO: how to convert configJson
        var configJson string
        i := new(Instance)
        c := new(Container)
        if err := rows.Scan(
            &c.Name, &nodeId, &configJson,
        ); err != nil {
            panic(err)
        }
        
        i.Node = GetNodeById(nodeId)
        i.Container = c

        l = append(l, i)
    }

    return l
}

func (i *Instance) Execute(serviceId int64) {
// TODO: how to convert to configJson
    stmt, err := db.Prepare(`
        INSERT INTO instance(
            name, nodeId, status, serviceId, configJson
        )
        VALUES(?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        i.Container.Name, i.Node, StatusRaw, serviceId, "",
    )
    if err != nil {
        panic(err)
    }
}
