#!/usr/bin/with-contenv bash

mkdir -p /data/logs/scheduler
mkdir -p /data/logs/thomas
mkdir -p /etc/services.d/thomas
cp /home/thomas_web.run /etc/services.d/thomas/run
chmod +x  /etc/services.d/thomas/*
ln -s /usr/lib/nvidia-390/libnvcuvid.so.390.12 /usr/local/lib/libnvcuvid.so
ln -s /usr/lib/nvidia-390/libnvcuvid.so.390.12 /usr/local/lib/libnvcuvid.so.1
ln -s /usr/lib/nvidia-390/libnvidia-encode.so /usr/local/lib/libnvidia-encode.so
ln -s /usr/lib/nvidia-390/libnvidia-encode.so /usr/local/lib/libnvidia-encode.so.1
ln -s /usr/lib/nvidia-390/libnvidia-encode.so /usr/local/lib/libnvidia-encode.so.390.12
ln -s /usr/lib/nvidia-390/libEGL_nvidia.so.390.12 /usr/local/lib/libEGL_nvidia.so

mkdir -p /etc/services.d/scheduler
cp /home/thomas_scheduler.run /etc/services.d/scheduler/run
chmod +x /etc/services.d/scheduler/*
