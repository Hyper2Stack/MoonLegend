package model

import (
)

type Repo struct {
    Id          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Owner       int64  `json:"owner"`
    YmlPath     string `json:"ymlPath"`
}

type RepoTag struct {
    Name   string `json:"name"`
}

////////////////////////////////////////////////////////////////////////////////

func ListRepo(cs []*Condition, o *Order, p *Paging) []*Repo {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, name, description, owner, ymlPath
        FROM
            repo
        ` + where + order + limit, vs...,
    )
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    l := make([]*Repo, 0)
    for rows.Next() {
        r := new(Repo)
        if err := rows.Scan(
            &r.Id, &r.Name, &r.Description, &r.Owner, &r.YmlPath,
        ); err != nil {
            panic(err)
        }

        l = append(l, r)
    }

    return l
}

func GetRepoById(id int64) *Repo {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("id", "=", id))

    l := ListRepo(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetRepoByNameAndOwner(name string, owner int64) *Repo {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("name", "=", name))
    conditions = append(conditions, NewCondition("owner", "=", owner))

    l := ListRepo(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func (r *Repo) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO repo(
            name, description, owner, ymlPath
        )
        VALUES(?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(
        r.Name, r.Description, r.Owner, r.YmlPath,
    )
    if err != nil {
        panic(err)
    }

    r.Id, err = result.LastInsertId()
    if err != nil {
        panic(err)
    }
}

func (r *Repo) Update() {
    stmt, err := db.Prepare(`
        UPDATE
            repo
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
        r.Name, r.Description, r.Id,
    ); err != nil {
        panic(err)
    }
}

func (r *Repo) Delete() {
    stmt, err := db.Prepare(`
        DELETE FROM
            repo
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.Id); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func (r *Repo) Tags() []*RepoTag {
    rows, err := db.Query(`
        SELECT
            name
        FROM
            repoTag
        WHERE
            repoId =  ?
        `, r.Id,
    )
    if err != nil {
        panic(err)
    }   
    defer rows.Close()
    
    l := make([]*RepoTag, 0)
    for rows.Next() { 
        t := new(RepoTag)
        if err := rows.Scan(
            &t.Name,
        ); err != nil {
            panic(err)
        }   
        
        l = append(l, t)
    }   
    
    return l
}

func (r *Repo) AddTag(t *RepoTag) {
    stmt, err := db.Prepare(`
        INSERT INTO repoTag(
            repoId, name
        )
        VALUES(?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        r.Id, t.Name,
    )
    if err != nil {
        panic(err)
    }
}

func (r *Repo) RemoveTag(name string) {
    stmt, err := db.Prepare(`
        DELETE FROM
            repoTag
        WHERE
            repoId = ? and name = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.Id, name); err != nil {
        panic(err)
    }
}
