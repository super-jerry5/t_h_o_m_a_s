#!/bin/bash

set -e
timetag=`date +%Y%m%d%H%M%S`


mkdir -p deploy 
mkdir -p deploy/web
mkdir -p deploy/scheduler
mkdir -p deploy/ffmpeg
echo "make deploy path succ"

cd ffmpeg 
./auto_make.sh
cd  ../
echo "make ffmpeg succ"

cd server
./auto_make.sh
cd ../

echo "make server(thomas/scheduler) succ"

cp -r docker/* ./ \
   && docker build --no-cache -t 172.18.199.246:5000/thomas:${timetag} -f Dockerfile.thomas .
docker push 172.18.199.246:5000/thomas:${timetag}
echo "New Docker Images:" "172.18.199.246:5000/thomas:"${timetag}
echo "New Docker Images:" "hub-docker.vrviu.com/thomas:"${timetag}


