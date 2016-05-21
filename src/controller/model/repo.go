package model

import (
    "fmt"
    "strings"
)

type Repo struct {
    Id          int64  `json:"id"`
    OwnerId     int64  `json:"owner_id"`
    Name        string `json:"name"`
    IsPublic    bool   `json:"is_public"`
    Description string `json:"description"`
    Readme      string `json:"readme"`
}

type RepoTag struct {
    Name   string `json:"name"`
    Yml    string `json:"yml"`
}

////////////////////////////////////////////////////////////////////////////////

func ListRepo(cs []*Condition, o *Order, p *Paging) []*Repo {
    where, vs := GenerateWhereSql(cs)
    order := GenerateOrderSql(o)
    limit := GenerateLimitSql(p)

    rows, err := db.Query(`
        SELECT
            id, owner_id, name, is_public, description, readme
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
            &r.Id, &r.OwnerId, &r.Name, &r.IsPublic, &r.Description, &r.Readme,
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

func GetRepoByNameAndOwnerId(name string, ownerId int64) *Repo {
    conditions := make([]*Condition, 0)
    conditions = append(conditions, NewCondition("name", "=", name))
    conditions = append(conditions, NewCondition("owner_id", "=", ownerId))

    l := ListRepo(conditions, nil, nil)
    if len(l) == 0 {
        return nil
    }

    return l[0]
}

func GetRepoByNamespaceAndName(namespace, name string) *Repo {
    u := GetUserByName(namespace)
    if u == nil {
        return nil
    }

    return GetRepoByNameAndOwnerId(name, u.Id)
}

func (r *Repo) Save() {
    stmt, err := db.Prepare(`
        INSERT INTO repo(
            owner_id, name, is_public, description, readme
        )
        VALUES(?, ?, ?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(
        r.OwnerId, r.Name, r.IsPublic, r.Description, r.Readme,
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
            is_public = ?,
            description = ?,
            readme = ?
        WHERE
            id = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(
        r.Name, r.IsPublic, r.Description, r.Readme, r.Id,
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
            name, yml
        FROM
            repoTag
        WHERE
            repo_id =  ?
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
            &t.Name, &t.Yml,
        ); err != nil {
            panic(err)
        }

        l = append(l, t)
    }

    return l
}

func (r *Repo) GetTag(name string) *RepoTag {
    for _, t := range r.Tags() {
        if t.Name == name {
            return t
        }
    }

    return nil
}

func (r *Repo) AddTag(t *RepoTag) {
    stmt, err := db.Prepare(`
        INSERT INTO repoTag(
            repo_id, name, yml
        )
        VALUES(?, ?, ?)
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        r.Id, t.Name, t.Yml,
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
            repo_id = ? and name = ?
    `)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(r.Id, name); err != nil {
        panic(err)
    }
}

////////////////////////////////////////////////////////////////////////////////

func ParseRepoString(s string) (*Repo, *RepoTag) {
    ss := strings.SplitN(s, "/", 2)
    if len(ss) != 2 {
        return nil, nil
    }

    namespace := ss[0]
    ss = strings.SplitN(ss[1], ":", 2)
    repo := GetRepoByNamespaceAndName(namespace, ss[0])
    if repo == nil {
        return nil, nil
    }

    if len(ss) < 2 {
        return repo, nil
    }

    return repo, repo.GetTag(ss[1])
}

func CompileRepoString(r *Repo, rt *RepoTag) string {
    if r == nil {
        return ""
    }

    s := fmt.Sprintf("%s/%s", GetUserById(r.OwnerId).Name, r.Name)
    if rt == nil {
        return s
    }

    return fmt.Sprintf("%s:%s", s, rt.Name)
}
