package model

import (
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Initialize(driver string, uri string) error {
    var err error 
    if db, err = sql.Open(driver, uri); err != nil {
        return err
    }

    return db.Ping()
}
