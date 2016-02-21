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
        setNodeGroupForeignKeyOwner,
        setNodeGroupForeignKeyRepo,
        nodeTable,
        setNodeForeignKey,
        nodeMembershipTable,
        setNodeMembershipForeignKeyGroup,
        setNodeMembershipForeignKeyNode,
        nodeTagTable,
        setNodeTagForeignKey,
        nicTable,
        setNicForeignKey,
        nicTagTable,
        setNicTagForeignKey,
        serviceTable,
        setServiceForeignKey,
        instanceTable,
        setInstanceForeignKeyNode,
        setInstanceForeignKeyService,
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
    displayName            VARCHAR(32),
    password               VARCHAR(64),
    userKey                VARCHAR(64),
    email                  VARCHAR(64),
    createTime             BIGINT,
    isAdmin                BOOLEAN DEFAULT false,
    isActive               BOOLEAN DEFAULT true,
    UNIQUE (name)
)AUTO_INCREMENT=10000;
`

// repo
var repoTable = `
CREATE TABLE IF NOT EXISTS repo (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(32),
    description MEDIUMBLOB,
    owner       INT,
    ymlPath     VARCHAR(64)
)AUTO_INCREMENT=10000;
`

var setRepoForeignKey = `
ALTER TABLE repo ADD CONSTRAINT fk__repo_user_owner FOREIGN KEY(owner) REFERENCES user(id) ON DELETE CASCADE;
`

// repoTag
var repoTagTable = `
CREATE TABLE IF NOT EXISTS repoTag (
    repoId      INT,
    name        VARCHAR(32),
    PRIMARY KEY (repoId, name)
);
`

var setRepoTagForeignKey = `
ALTER TABLE repoTag ADD CONSTRAINT fk__repoTag_repo_id FOREIGN KEY(repoId) REFERENCES repo(id) ON DELETE CASCADE;
`

// nodeGroup
var nodeGroupTable = `
CREATE TABLE IF NOT EXISTS nodeGroup (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(32),
    description MEDIUMBLOB,
    owner       INT,
    repoId      INT,
    status      INT
)AUTO_INCREMENT=10000;
`

var setNodeGroupForeignKeyOwner = `
ALTER TABLE nodeGroup ADD CONSTRAINT fk__nodeGroup_owner_id FOREIGN KEY(owner) REFERENCES user(id) ON DELETE CASCADE;
`

var setNodeGroupForeignKeyRepo = `
ALTER TABLE nodeGroup ADD CONSTRAINT fk__nodeGroup_repo_id FOREIGN KEY(repoId) REFERENCES repo(id) ON DELETE CASCADE;
`

// node
var nodeTable = `
CREATE TABLE IF NOT EXISTS node (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(32),
    description MEDIUMBLOB,
    owner       INT
)AUTO_INCREMENT=10000;
`

var setNodeForeignKey = `
ALTER TABLE node ADD CONSTRAINT fk__node_user_owner FOREIGN KEY(owner) REFERENCES user(id) ON DELETE CASCADE;
`

// nodeMembership
var nodeMembershipTable = `
CREATE TABLE IF NOT EXISTS nodeMembership (
    groupId INT,
    nodeId  INT,
    PRIMARY KEY (groupId, nodeId)
);
`

var setNodeMembershipForeignKeyGroup = `
ALTER TABLE nodeMembership ADD CONSTRAINT fk__nodeMembership_group_id FOREIGN KEY(groupId) REFERENCES nodeGroup(id) ON DELETE CASCADE;
`

var setNodeMembershipForeignKeyNode = `
ALTER TABLE nodeMembership ADD CONSTRAINT fk__nodeMembership_node_id FOREIGN KEY(nodeId) REFERENCES node(id) ON DELETE CASCADE;
`

// nodeTag
var nodeTagTable = `
CREATE TABLE IF NOT EXISTS nodeTag (
    nodeId      INT,
    name        VARCHAR(32),
    PRIMARY KEY (nodeId, name)
);
`

var setNodeTagForeignKey = `
ALTER TABLE nodeTag ADD CONSTRAINT fk__nodeTag_node_id FOREIGN KEY(nodeId) REFERENCES node(id) ON DELETE CASCADE;
`

// nic
var nicTable = `
CREATE TABLE IF NOT EXISTS nic (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(32),
    nodeId      INT,
    ip4Addr     VARCHAR(32)
)AUTO_INCREMENT=10000;
`

var setNicForeignKey = `
ALTER TABLE nic ADD CONSTRAINT fk__nic_node_id FOREIGN KEY(nodeId) REFERENCES node(id) ON DELETE CASCADE;
`

// nicTag
var nicTagTable = `
CREATE TABLE IF NOT EXISTS nicTag (
    nicId       INT,
    name        VARCHAR(32),
    PRIMARY KEY (nicId, name)
);
`

var setNicTagForeignKey = `
ALTER TABLE nicTag ADD CONSTRAINT fk__nicTag_nic_id FOREIGN KEY(nicId) REFERENCES nic(id) ON DELETE CASCADE;
`

// service
var serviceTable = `
CREATE TABLE IF NOT EXISTS service (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(32),
    groupId     INT,
    configJson  MEDIUMBLOB
)AUTO_INCREMENT=10000;
`

var setServiceForeignKey = `
ALTER TABLE service ADD CONSTRAINT fk__service_group_id FOREIGN KEY(groupId) REFERENCES nodeGroup(id) ON DELETE CASCADE;
`

// instance
var instanceTable = `
CREATE TABLE IF NOT EXISTS instance (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(32),
    nodeId      INT,
    serviceId   INT,
    status      INT,
    configJson  MEDIUMBLOB
)AUTO_INCREMENT=10000;
`

var setInstanceForeignKeyNode = `
ALTER TABLE instance ADD CONSTRAINT fk__instance_node_id FOREIGN KEY(nodeId) REFERENCES node(id) ON DELETE CASCADE;
`

var setInstanceForeignKeyService = `
ALTER TABLE instance ADD CONSTRAINT fk__instance_service_id FOREIGN KEY(serviceId) REFERENCES service(id) ON DELETE CASCADE;
`

// Init DB
// CREATE DATABASE moonlegend CHARACTER SET utf8 COLLATE utf8_general_ci;

