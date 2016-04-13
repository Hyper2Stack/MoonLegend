package model

import (
    "encoding/json"
)

type Group struct {
    Id              int64               `json:"id"`
    Name            string              `json:"name"`
    Description     string              `json:"description"`
    OwnerId         int64               `json:"owner_id"`
    Status          string              `json:"status"`
    Deployment      *Deployment         `json:"deployment"`
    Process         []*InstanceStatus   `json:"process"`
}

type InstanceStatus struct {
    Service  string `json:"service"`
    Instance string `json:"instance"`
    Status   string `json:"status"`
}

type Deployment struct {
    Repo       string               `json:"repo"`
    Services   map[string]*Service  `json:"services"`
}

type Service struct {
    Image      string      `json:"image"`
    Depends    []string    `json:"depends"`
    Instances  []*Instance `json:"instance"`
}

type Instance struct {
    Node        *Node         `json:"node"`
    Container   *Container    `json:"container"`
    Entrypoints []*Entrypoint `json:"entrypoints"`
}

type Entrypoint struct {
    Protocol      string `json:"protocol"`
    ListeningAddr string `json:"listening_addr"`
    ListeningPort string `json:"listening_port"`
    ContainerPort string `json:"container_port"`
}

type Container struct {
    Name         string `json:"name"`
    RunCommand   string `json:"run_command"`
    RmCommand    string `json:"rm_command"`
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
