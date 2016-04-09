package migrate

import (
    "github.com/BurntSushi/migration"
    _ "github.com/go-sql-driver/mysql"
)

func Run(uri string) error {
    migration.DefaultGetVersion = GetVersion
    migration.DefaultSetVersion = SetVersion
    var migrations = []migration.Migrator{
        Migrate_1,
    }

    db, err := migration.Open("mysql", uri, migrations)
    if err != nil {
        return err
    }

    db.Close()
    return nil
}
