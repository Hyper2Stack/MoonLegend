#!/bin/bash

cid=`sudo docker ps | grep nginx:1.9 | grep -v grep | awk '{print $1}'`
sudo docker stop $cid

ps -ef | grep moonlegend | grep -v grep | awk '{print $2}' | xargs kill -9
sudo /usr/sbin/moon -s quit

