#!/bin/bash

image="daocloud.io/nginx:1.9"

curl -H "Authorization:test" -XPUT 127.0.0.1:8080/api/v1/user/test -d '{"commands":[{"command":"docker", "args":["pull", "'$image'"], "restrict":true}]}'
curl -H "Authorization:test" -XPUT 127.0.0.1:8080/api/v1/user/test -d '{"commands":[{"command":"docker", "args":["run", "'$image'"], "restrict":true}]}'
