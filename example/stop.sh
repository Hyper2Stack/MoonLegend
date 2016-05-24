#!/bin/bash

sudo /usr/sbin/moon -s quit > /dev/null 2>&1
ps -ef | grep moonlegend | grep -v grep | awk '{print $2}' | xargs kill -9
sudo docker rm -vf moonlegend_db > /dev/null 2>&1
