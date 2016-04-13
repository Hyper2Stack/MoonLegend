package migrate

import (
    "github.com/BurntSushi/migration"
)

//////////////////////////////////////////////////////////////////////

func Migrate_1(tx migration.LimitedTx) error {
    scripts := []string{
        userTable,
        repoTable,
        setRepoForeignKey,
        repoTagTable,
        setRepoTagForeignKey,
        nodeGroupTable,
        setNodeGroupForeignKeyOwnerId,
        nodeTable,
        setNodeForeignKey,
        nodeMembershipTable,
        setNodeMembershipForeignKeyGroup,
        setNodeMembershipForeignKeyNode,
    }

    for _, cmd := range scripts {
        if _, err := tx.Exec(cmd); err != nil {
            return err
        }
    }

    return nil
}

// user
var userTable = `
CREATE TABLE IF NOT EXISTS user (
    id                     INT AUTO_INCREMENT PRIMARY KEY,
    name                   VARCHAR(32),
    display_name           VARCHAR(32),
    password               VARCHAR(64),
    agent_key              VARCHAR(64),
    email                  VARCHAR(64),
    create_ts              BIGINT,
    is_admin               BOOLEAN DEFAULT false,
    is_active              BOOLEAN DEFAULT true,
    UNIQUE (name)
)AUTO_INCREMENT=10000;
`

// repo
var repoTable = `
CREATE TABLE IF NOT EXISTS repo (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(32),
    is_public   BOOLEAN,
    description MEDIUMBLOB,
    readme      MEDIUMBLOB,
    owner_id    INT
)AUTO_INCREMENT=10000;
`

var setRepoForeignKey = `
ALTER TABLE repo ADD CONSTRAINT fk__repo_user_owner FOREIGN KEY(owner_id) REFERENCES user(id) ON DELETE CASCADE;
`

// repoTag
var repoTagTable = `
CREATE TABLE IF NOT EXISTS repoTag (
    repo_id     INT,
    name        VARCHAR(32),
    yml         MEDIUMBLOB,
    PRIMARY KEY (repo_id, name)
);
`

var setRepoTagForeignKey = `
ALTER TABLE repoTag ADD CONSTRAINT fk__repoTag_repo_id FOREIGN KEY(repo_id) REFERENCES repo(id) ON DELETE CASCADE;
`

// nodeGroup
var nodeGroupTable = `
CREATE TABLE IF NOT EXISTS nodeGroup (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(32),
    description     MEDIUMBLOB,
    owner_id        INT,
    status          VARCHAR(32),
    deployment      MEDIUMBLOB,
    process         MEDIUMBLOB
)AUTO_INCREMENT=10000;
`

var setNodeGroupForeignKeyOwnerId = `
ALTER TABLE nodeGroup ADD CONSTRAINT fk__nodeGroup_owner_id FOREIGN KEY(owner_id) REFERENCES user(id) ON DELETE CASCADE;
`

// node
var nodeTable = `
CREATE TABLE IF NOT EXISTS node (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    uuid        VARCHAR(128),
    name        VARCHAR(32),
    description MEDIUMBLOB,
    status      VARCHAR(32),
    owner_id    INT,
    tags        MEDIUMBLOB,
    nics        MEDIUMBLOB
)AUTO_INCREMENT=10000;
`

var setNodeForeignKey = `
ALTER TABLE node ADD CONSTRAINT fk__node_user_owner FOREIGN KEY(owner_id) REFERENCES user(id) ON DELETE CASCADE;
`

// nodeMembership
var nodeMembershipTable = `
CREATE TABLE IF NOT EXISTS nodeMembership (
    group_id INT,
    node_id  INT,
    PRIMARY KEY (group_id, node_id)
);
`

var setNodeMembershipForeignKeyGroup = `
ALTER TABLE nodeMembership ADD CONSTRAINT fk__nodeMembership_group_id FOREIGN KEY(group_id) REFERENCES nodeGroup(id) ON DELETE CASCADE;
`

var setNodeMembershipForeignKeyNode = `
ALTER TABLE nodeMembership ADD CONSTRAINT fk__nodeMembership_node_id FOREIGN KEY(node_id) REFERENCES node(id) ON DELETE CASCADE;
`

// Init DB
// CREATE DATABASE moonlegend CHARACTER SET utf8 COLLATE utf8_general_ci;

