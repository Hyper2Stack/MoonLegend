#!/bin/bash

uuid=`sudo cat /etc/moon/key.json | python -m json.tool | grep "uuid" | sed "s/\"$//g" | sed "s/^.*\"//g"`

curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/node.info
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/agent.info
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/script.exec -d '{"commands":[{"command":"ls", "args":["/tmp/a/notfound"], "restrict":true}]}'
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/script.exec -d '{"commands":[{"command":"ls", "args":["/tmp/a/notfound"], "restrict":false}]}'
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/script.exec -d '{"commands":[{"command":"ls", "args":["/usr"], "restrict":true}]}'
curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/file.create -d '{"path":"/tmp/xyz", "mode":"644", "content":"123\n456"}'
sudo cat /tmp/xyz
echo ""

curl -XPUT 127.0.0.1:8080/api/v1/debug-agent/${uuid}/not-supported
