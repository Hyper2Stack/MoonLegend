package model

import (
    "encoding/json"
)

type Node struct {
    Id          int64    `json:"id"`
    Uuid        string   `json:"uuid"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Status      string   `json:"status"`
    OwnerId     int64    `json:"owner_id"`
    Tags        []string `json:"tags"`
    Nics        []*Nic   `json:"nics"`
}

type Nic struct {
    Name    string   `json:"name"`
    Ip4Addr string   `json:"ip"`
    Tags    []string `json:"tags"`
}

////////////////////////////////////////////////////////////////////////////////

func ListNode(cs []*Condition, o *Order, p *Paging) []*Node {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, uuid, name, description, status, owner_id, tags, nics
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

func GetNodeById(id int64) *Node {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListNode(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetNodeByNameAndOwnerId(name string, ownerId int64) *Node {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("name", "=", name))
    conditions = append(conditions, NewCondition("owner_id", "=", ownerId))

    l := ListNode(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (n *Node) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO node(
            uuid, name, description, status, owner_id, tags, nics
        )
        VALUES(?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    tagInfo, err := json.Marshal(n.Tags)
    if err != nil {
        panic(err)
    }

    nicInfo, err := json.Marshal(n.Nics)
    if err != nil {
        panic(err)
    }

    result, err := stmt.Exec(
        n.Uuid, n.Name, n.Description, n.Status, n.OwnerId, string(tagInfo), string(nicInfo),
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
            status = ?,
            tags = ?,
            nics = ?
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    tagInfo, err := json.Marshal(n.Tags)
    if err != nil {
        panic(err)
    }

    nicInfo, err := json.Marshal(n.Nics)
    if err != nil {
        panic(err)
    }

    if _, err := stmt.Exec(
        n.Name, n.Description, n.Status, string(tagInfo), string(nicInfo), n.Id,
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

func (n *Node) AddTag(t string) {
    n.Tags = append(n.Tags, t)
    n.Update()
}

func (n *Node) RemoveTag(t string) {
    removeElement(&n.Tags, t)
    n.Update()
}

func (n *Node) AddNicTag(nicName, t string) {
    for _, nic := range n.Nics {
        if nic.Name == nicName {
            nic.Tags = append(nic.Tags, t)
            break
        }
    }

    n.Update()
}

func (n *Node) RemoveNicTag(nicName, t string) {
    for _, nic := range n.Nics {
        if nic.Name == nicName {
            removeElement(&nic.Tags, t)
            break
        }
    }

    n.Update()
}
