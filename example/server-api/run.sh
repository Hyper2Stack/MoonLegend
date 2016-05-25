#!/bin/bash

### ping
curl -i localhost:8080/api/v1/ping

### user
session_key=`curl -s localhost:8080/api/v1/login  -X POST -d '{"username":"hyper2stack", "password":"password"}' | sed 's/"}$//g' | sed 's/^.*"//g'`
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user

### repo
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos -X POST -d '{"name":"p001"}'
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos/p001
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos/p001 -X PUT -d '{"is_public":true}'

### repo tag
cat > /tmp/1 << EOF
{"name":"v1", "yml":"infrastructure:
  network:
  - default

services:
  web:
    image: daocloud.io/haipeng/web:1.0
    ports:
    - 80/tcp
    environment:
    - MYSQL_ADDRESS={{with .Singleton \"mysql\"}}{{.Address}}{{end}}
    - MYSQL_PORT={{with .Singleton \"mysql\"}}{{.Port}}{{end}}
    - MYSQL_PASSWORD=root
    - MYSQL_PASSWORD={{.Runtime.Env \"MysqlRootPassword\"}}
    - MYSQL_DATABASE={{.Runtime.Env \"MysqlDatabaseName\"}}
    depends_on:
    - mysql
    networks:
    - default

  mysql:
    image: daocloud.io/mysql:5.5
    singleton: true
    ports:
    - 3306
    environment:
    - MYSQL_ROOT_PASSWORD={{.Runtime.Env \"MysqlRootPassword\"}}
    - MYSQL_DATABASE={{.Runtime.Env \"MysqlDatabaseName\"}}
    networks:
    - default

runtime:
  env:
  - MysqlRootPassword=password
  - MysqlDatabaseName=test
  global_policy:
    restart: always
    port_mapping: fixed
  service_policy:
    web:
      instance_num: 1"}
EOF

sed ':a;N;$!ba;s/\n/\\n/g' /tmp/1 > /tmp/2

curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos/p001/tags
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos/p001/tags -X POST --data @/tmp/2

### node
host_name=`hostname`
interface=`ip route get 8.8.8.8 | head -1 | awk '{print $5}'`

curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes/${host_name}
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes/${host_name}/tags -X POST -d '{"name":"mysql"}'
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes/${host_name}/tags -X POST -d '{"name":"web"}'
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes/${host_name}/nics/${interface}/tags -X POST -d '{"name":"default"}'

### group
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups -X POST -d '{"name":"first"}'
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first/nodes
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first/nodes -X POST -d '{"name":"'${host_name}'"}'

### deployment
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first/deployment
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first/deployment -X POST -d '{"repo":"hyper2stack/p001:v1"}'
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first/deployment

### destroy
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first/deployment -X DELETE
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first/nodes/${host_name} -X DELETE
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/groups/first -X DELETE
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes/${host_name}/nics/${interface}/tags/default -X DELETE
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes/${host_name}/tags/mysql -X DELETE
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/nodes/${host_name}/tags/web -X DELETE
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos/p001/tags/v1 -X DELETE
curl -i -H "Authorization: ${session_key}" localhost:8080/api/v1/user/repos/p001 -X DELETE
