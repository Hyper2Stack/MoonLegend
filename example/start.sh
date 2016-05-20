#!/bin/bash

# Change your settings before separator
mysql_username="root"
mysql_password="root"
key="abc"

# ==================================================================
ROOT_DIR=$(readlink -f $(dirname $0)/..)

mysql -u${mysql_username} -p${mysql_password} -e "create database if not exists moonlegend default character set utf8 collate utf8_general_ci;"
sed -e "s/root:root/${mysql_username}:${mysql_password}/g" ${ROOT_DIR}/config/moonlegend.json > /tmp/moonlegend.json

${ROOT_DIR}/bin/moonlegend -conf /tmp/moonlegend.json &

sleep 5

res=`mysql -u${mysql_username} -p${mysql_password} -e "use moonlegend;select * from user where name='test'"| grep test`
if [ -z "$res" ]; then
  mysql -u${mysql_username} -p${mysql_password} -e "use moonlegend;insert into user (name, display_name, password, agent_key, email, create_ts) values ('test', 'test user', 'test', '"${key}"', 'test@moonlegend.com', `date +'%s'`);"
fi

sudo /usr/sbin/moon-config -key ${key}
sudo /usr/sbin/moon
