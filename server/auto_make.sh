#!/bin/bash

go build
cp thomas ../deploy/web/
cp conf ../deploy/web/
cp static ../deploy/web/
cp html ../deploy/web

cd scheduler 
go build
cp scheduler ../../deploy/scheduler
cp conf ../../deploy/scheduler/

