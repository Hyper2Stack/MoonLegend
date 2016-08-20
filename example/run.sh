#!/bin/bash

[ $# -ne 1 ] && echo "usage: $0 <env/config_file>" && exit 1

export MOONLEGEND_SERVER=localhost:8080
export ROOT_DIR=`dirname $0`/..
export HOST_NAME=`hostname`
export INTERFACE_NAME=`ip route get 8.8.8.8 | head -1 | awk '{print $5}'`
export YML=${1}.yml

## pre download images, avoid waiting long time for preparing
echo "# pre pull image ..."
docker inspect daocloud.io/hyper2stack/web:1.2 > /dev/null 2>&1 || docker pull daocloud.io/hyper2stack/web:1.2
docker inspect daocloud.io/mysql:5.5 > /dev/null 2>&1 || docker pull daocloud.io/mysql:5.5

## init
echo ""
echo "# ping moonlegend server ..."
${ROOT_DIR}/bin/ml ping

echo ""
echo "# login ..."
${ROOT_DIR}/bin/ml login hyper2stack password

echo ""
echo "# create repo p001 ..."
${ROOT_DIR}/bin/ml create-repo p001

echo ""
echo "# list repos ..."
${ROOT_DIR}/bin/ml list-repo

echo ""
echo "# create tag v1 with ${YML} ..."
${ROOT_DIR}/bin/ml create-repo-tag p001 v1 `dirname $0`/${YML}

echo ""
echo "# list tags of repo p001 ..."
${ROOT_DIR}/bin/ml list-repo-tag p001

echo ""
echo "# create tag on node ${HOST_NAME} ..."
${ROOT_DIR}/bin/ml create-node-tag ${HOST_NAME} mysql
${ROOT_DIR}/bin/ml create-node-tag ${HOST_NAME} web

echo ""
echo "# show details of node ${HOST_NAME} ..."
${ROOT_DIR}/bin/ml show-node ${HOST_NAME}

echo ""
echo "# create group g001 ..."
${ROOT_DIR}/bin/ml create-group g001

echo ""
echo "# list group ..."
${ROOT_DIR}/bin/ml list-group

echo ""
echo "# add node ${HOST_NAME} to group g001 ..."
${ROOT_DIR}/bin/ml add-node-to-group g001 ${HOST_NAME}

echo ""
echo "# list node of group g001 ..."
${ROOT_DIR}/bin/ml list-group-node g001

echo ""
echo "# init deployment of group g001 ..."
${ROOT_DIR}/bin/ml deploy-init g001 hyper2stack/p001:v1

echo ""
echo "# prepare deployment of group g001 ..."
${ROOT_DIR}/bin/ml deploy-prepare g001

echo ""
echo "# show details of group g001 ..."
${ROOT_DIR}/bin/ml show-group g001

echo ""
echo "# waiting until prepare complete ..."
for x in `seq 1 10`
do
  sleep 3
  ${ROOT_DIR}/bin/ml show-group g001 | grep "^status.*prepared" > /dev/null
  if [ $? -eq 0 ]; then
    echo "prepare complete in $((3*x))s"
    break
  fi
  [ $x -eq 10 ] && echo "prepare can not complete in 30s"
done

echo ""
echo "# execute deployment of group g001 ..."
${ROOT_DIR}/bin/ml deploy-execute g001

echo ""
echo "# show details of group g001 ..."
${ROOT_DIR}/bin/ml show-group g001

echo ""
echo "# waiting until deploy complete ..."
done=0
for x in `seq 1 10`
do
  sleep 3
  ${ROOT_DIR}/bin/ml show-group g001 | grep "^status.*deployed" > /dev/null
  if [ $? -eq 0 ]; then
    echo "deploy complete in $((3*x))s"
    done=1
    break
  fi
done

if [ $done -eq 1 ]; then
  sleep 5

  echo ""
  echo "# check if web app is deployed ..."
  sudo docker ps | egrep "web_0|mysql_0"

  echo ""
  echo "# test app with command 'curl -i localhost:8000/hello' ..."
  curl -i localhost:8000/hello

  echo ""
  echo "# test app with command 'curl -i localhost:8000/env' ..."
  curl -i localhost:8000/env
else
  echo "deploy can not complete in 30s"
fi

## destroy
echo ""
echo "# clear deployment of group g001 ..."
${ROOT_DIR}/bin/ml deploy-clear g001

echo ""
echo "# waiting until clear complete ..."
for x in `seq 1 10`
do
  sleep 3
  ${ROOT_DIR}/bin/ml show-group g001 | grep "^status.*prepared" > /dev/null
  if [ $? -eq 0 ]; then
    echo "clear complete in $((3*x))s"
    break
  fi
  [ $x -eq 10 ] && echo "clear can not complete in 30s"
done

echo ""
echo "# delete deployment of group g001 ..."
${ROOT_DIR}/bin/ml deploy-delete g001

echo ""
echo "# remove node ${HOST_NAME} from group g001 ..."
${ROOT_DIR}/bin/ml remove-node-from-group g001 ${HOST_NAME}

echo ""
echo "# delete group g001 ..."
${ROOT_DIR}/bin/ml delete-group g001

echo ""
echo "# delete tag of node ${HOST_NAME} ..."
${ROOT_DIR}/bin/ml delete-node-tag ${HOST_NAME} mysql
${ROOT_DIR}/bin/ml delete-node-tag ${HOST_NAME} web

echo ""
echo "# delete tag v1 of repo p001 ..."
${ROOT_DIR}/bin/ml delete-repo-tag p001 v1

echo ""
echo "# delete repo p001 ..."
${ROOT_DIR}/bin/ml delete-repo p001
