#!/bin/bash

export MOONLEGEND_SERVER=localhost:8080
export ROOT_DIR=`dirname $0`/..

echo "starting moonlegend server ..."
${ROOT_DIR}/ctrl start

sleep 5

echo ""
echo "register user 'hyper2stack' with password 'password' ..."
${ROOT_DIR}/bin/ml signup hyper2stack password
${ROOT_DIR}/bin/ml login hyper2stack password

echo ""
echo "setup moon agent ..."
agent_key=`${ROOT_DIR}/bin/ml profile | grep "^key" | sed "s/^key:\s*//g"`
sudo /usr/sbin/moon-config -key ${agent_key}
sudo service moon restart
