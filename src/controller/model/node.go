package model

import (
)

type Node struct {
    Id          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Owner       int64  `json:"owner"`
}

type NodeTag struct {
    Name   string `json:"name"`
}

type Nic struct {
    Id      string `json:id`
    Name    string `json:"name"`
    Ip4Addr string `json:"ip"`
}

type NicTag struct {
    Name   string `json:"name"`
}

////////////////////////////////////////////////////////////////////////////////

func ListNode(cs []*Condition, o *Order, p *Paging) []*Node {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, name, description, owner
        FROM
            node
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Node, 0)
    for rows.Next() {
        n := new(Node)
        if err := rows.Scan(
            &n.Id, &n.Name, &n.Description, &n.Owner,
        ); err != nil {
            panic(err)
        }

        l = append(l, n)
    }

    return l
}

func GetNodeById(id int64) *Node {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListNode(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetNodeByNameAndOwner(name string, owner int64) *Node {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("name", "=", name))
    conditions = append(conditions, NewCondition("owner", "=", owner))

    l := ListNode(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetNodesByGroupId(groupId int64) []*Node {
    rows, err := db.Query(`
        SELECT
            id, name, description, owner
        FROM
            node, nodeMembership
        WHERE
            groupId =  ? and id = nodeId
        `, groupId,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Node, 0)
    for rows.Next() {
        n := new(Node)
        if err := rows.Scan(
            &n.Id, &n.Name, &n.Description, &n.Owner,
        ); err != nil {
            panic(err)
        }

        l = append(l, n)
    }

    return l
}

func (n *Node) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO node(
            name, description, owner
        )
        VALUES(?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(
        n.Name, n.Description, n.Owner,
    )
    if err != nil {
        panic(err)
    }

    n.Id, err = result.LastInsertId()
    if err != nil {
        panic(err)
    }
}

func (n *Node) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            node
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
        n.Name, n.Description,
    ); err != nil {
        panic(err)
    }
}

func (n *Node) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            node
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(n.Id); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (n *Node) Tags() []*NodeTag {
    rows, err := db.Query(`
        SELECT
            name
        FROM
            nodeTag
        WHERE
            nodeId =  ?
        `, n.Id,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*NodeTag, 0)
    for rows.Next() {
        t := new(NodeTag)
        if err := rows.Scan(
            &t.Name,
        ); err != nil {
            panic(err)
        }

        l = append(l, t)
    }

    return l
}

func (n *Node) AddTag(t *NodeTag) {
    stmt, err := db.Prepare(`
        INSERT INTO nodeTag(
            nodeId, name
        )
        VALUES(?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        n.Id, t.Name,
    )
    if err != nil {
        panic(err)
    }
}

func (n *Node) RemoveTag(name string) {
    stmt, err := db.Prepare(`
        DELETE FROM
            nodeTag
        WHERE
            nodeId = ? and name = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(n.Id, name); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func ListNic(cs []*Condition, o *Order, p *Paging) []*Nic {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            name, ip4Addr
        FROM
            nic
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Nic, 0)
    for rows.Next() {
        n := new(Nic)
        if err := rows.Scan(
            &n.Name, &n.Ip4Addr,
        ); err != nil {
            panic(err)
        }

        l = append(l, n)
    }

    return l
}

func (n *Node) Nics() []*Nic {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("nodeId", "=", n.Id))

    l := ListNic(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l
}

func (n *Node) GetNicByName(name string) *Nic {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("nodeId", "=", n.Id))
    conditions = append(conditions, NewCondition("name", "=", name))

    l := ListNic(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (n *Node) AddNic(nic *Nic) {
    stmt, err := db.Prepare(`
        INSERT INTO nic(
            nodeId, name, ip4addr
        )
        VALUES(?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        n.Id, nic.Name, nic.Ip4Addr,
    )
    if err != nil {
        panic(err)
    }
}

func (n *Node) RemoveNic(name string) {
    stmt, err := db.Prepare(`
        DELETE FROM
            nic
        WHERE
            nodeId = ? and name = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(n.Id, name); err != nil {
        panic(err)
    }
}
////////////////////////////////////////////////////////////////////////////////

func (n *Nic) NicTags() []*NicTag {
    rows, err := db.Query(`
        SELECT
            name
        FROM
            nicTag
        WHERE
            nicId =  ?
        `, n.Id,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*NicTag, 0)
    for rows.Next() {
        t := new(NicTag)
        if err := rows.Scan(
            &t.Name,
        ); err != nil {
            panic(err)
        }

        l = append(l, t)
    }

    return l
}

func (n *Nic) AddNicTag(t *NicTag) {
    stmt, err := db.Prepare(`
        INSERT INTO nicTag(
            nicId, name
        )
        VALUES(?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        n.Id, t.Name,
    )
    if err != nil {
        panic(err)
    }
}

func (n *Nic) RemoveNicTag(name string) {
    stmt, err := db.Prepare(`
        DELETE FROM
            nicTag
        WHERE
            nicId = ? and name = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(n.Id, name); err != nil {
        panic(err)
    }
}
