#!/bin/bash

uuid=`cat /etc/moon/key.cfg | python -m json.tool | grep "uuid" | sed "s/\"$//g" | sed "s/^.*\"//g"`

curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/ping
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/get-node-info
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/get-agent-info
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/exec-shell-script -d '{"commands":[{"command":"rm", "args":["-rf", "/tmp/abc"], "restrict":true}]}'
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/exec-shell-script -d '{"commands":[{"command":"touch", "args":["/tmp/abc"], "restrict":true}]}'
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/not-supported
