#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))
IP=$(ip route get 8.8.8.8 | head -1 | awk '{print $7}')
MYSQL_PORT=13306

echo "init mysql..."
sudo docker inspect moonlegend_db > /dev/null 2>&1
if [ $? -ne 0 ]; then
  sudo docker run -d --name moonlegend_db -p ${MYSQL_PORT}:3306 -e MYSQL_ROOT_PASSWORD=root daocloud.io/mysql:5.5 > /dev/null
  sleep 15
  sudo docker exec moonlegend_db mysql -uroot -proot -e "create database if not exists moonlegend default character set utf8 collate utf8_general_ci;"
fi

echo "init moonlegend server..."
sed -e "s/127\.0\.0\.1:3306/${IP}:${MYSQL_PORT}/g" ${ROOT_DIR}/config/moonlegend.json > /tmp/moonlegend.json
${ROOT_DIR}/bin/moonlegend -conf /tmp/moonlegend.json &
sleep 5

echo "done."
echo ""
echo "====== Env Info ======"
echo "moonlegend config: /tmp/moonlegend.json"
echo "moonlegend log:    /tmp/moonlegend.log"
echo "moonlegend mysql:  root:root@tcp(${IP}:${MYSQL_PORT})/moonlegend"
