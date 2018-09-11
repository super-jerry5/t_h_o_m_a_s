#!/bin/bash

go build
cp server ../deploy/web/thomas
cp -rf conf ../deploy/web/
cp -rf static ../deploy/web/
cp -rf html ../deploy/web

cd scheduler 
go build
cp scheduler ../../deploy/scheduler/
cp -rf conf ../../deploy/scheduler/

