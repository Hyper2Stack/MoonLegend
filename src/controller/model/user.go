package model

import (
    "time"
)

type User struct {
    Id               int64  `json:"id"`
    Name             string `json:"name"`
    DisplayName      string `json:"display_name"`
    Key              string `json:"key"`
    Email            string `json:"email"`
    CreateTime       int64  `json:"create_ts"`
}

////////////////////////////////////////////////////////////////////////////////

func ListUser(cs []*Condition, o *Order, p *Paging) []*User {
    if cs == nil {
        cs = make([]*Condition, 0)
    }

    cs = append(cs, NewCondition("is_active", "=", true))

    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, name, display_name, agent_key, email, create_ts
        FROM
            user
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*User, 0)
    for rows.Next() {
        u := new(User)
        if err := rows.Scan(
            &u.Id, &u.Name, &u.DisplayName, &u.Key, &u.Email, &u.CreateTime,
        ); err != nil {
            panic(err)
        }

        l = append(l, u)
    }

    return l
}

func GetUserById(id int64) *User {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListUser(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetUserByName(name string) *User {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("name", "=", name))

    l := ListUser(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetUserByAgentkey(key string) *User {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("agent_key", "=", key))

    l := ListUser(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func IsAdmin(id int64) bool {
    rows, err := db.Query(`
        SELECT
            is_admin
        FROM
            user
        WHERE
            id = ?
        `, id,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    var b bool
    for rows.Next() {
        if err := rows.Scan(&b); err != nil {
            panic(err)
        }
    }

    return b
}

func (u *User) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO user(
            name, display_name, agent_key, email, create_ts, is_active
        )
        VALUES(?, ?, ?, ?, ?, true)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    u.CreateTime = time.Now().UTC().Unix()

    result, err := stmt.Exec(
        u.Name, u.DisplayName, u.Key, u.Email, u.CreateTime,
    )
    if err != nil {
        panic(err)
    }

    u.Id, err = result.LastInsertId()
    if err != nil {
        panic(err)
    }
}

func (u *User) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            user
        SET
            display_name = ?,
            email = ?,
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        u.DisplayName,
        u.Email,
        u.Id,
    ); err != nil {
        panic(err)
    }
}

func (u *User) Delete() {
    stmt, err := db.Prepare(`
        UPDATE
            user
        SET
            is_active = false,
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(u.Id); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (u *User) VerifyPassword(password string) bool {
    rows, err := db.Query(`
        SELECT
            password
        FROM
            user
        WHERE
            name = ?
        `, u.Name,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    var hp string
    for rows.Next() {
        if err := rows.Scan(&hp); err != nil {
            panic(err)
        }
    }

    return hashPassword(password) == hp
}

func (u *User) ResetPassword(password string) {
    stmt, err := db.Prepare(`
        UPDATE
            user
        SET
            password = ?
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        hashPassword(password),
        u.Id,
    ); err != nil {
        panic(err)
    }
}

func (u *User) ResetKey() {
    stmt, err := db.Prepare(`
        UPDATE
            user
        SET
            agent_key = ?
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        generateKey(),
        u.Id,
    ); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (u *User) Repos() []*Repo {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("owner_id", "=", u.Id))

    return ListRepo(conditions, nil, nil)
}

func (u *User) Nodes() []*Node {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("owner_id", "=", u.Id))

    return ListNode(conditions, nil, nil)
}

func (u *User) Groups() []*Group {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("owner_id", "=", u.Id))

    return ListGroup(conditions, nil, nil)
}

func (u *User) HasRepo(repo *Repo) bool {
    return repo.OwnerId == u.Id
}

func (u *User) CanDeploy(repo *Repo) bool {
    return repo.OwnerId == u.Id || repo.IsPublic
}
