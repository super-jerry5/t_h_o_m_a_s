#!/bin/bash

go build
cp thomas ../deploy/web/
cd scheduler 
go build
cp scheduler ../../deploy/scheduler


