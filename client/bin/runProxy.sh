#!/usr/bin/env bash

nohup ./adslProxyClient-linux -host http://proxy.mumulife.com:9010 -id 1 -port 55155 -changeInterval 40 > client.log 2>&1 &

while :
do
    if [ $(ps -ef | grep "./gost" | grep -v "grep" | wc -l) -gt 0 ];then
        kill $(ps -ef|grep "./gost" |grep -v "grep"|awk '{print $2}')
        echo "killed"
        sleep 1
    echo "restart"
    else
    echo "start"
    fi
    nohup ./gost -L=:55155 > gost.log 2>&1 &
    sleep 300
done

