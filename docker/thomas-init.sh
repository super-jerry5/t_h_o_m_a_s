#!/usr/bin/with-contenv bash
mkdir -p /app/lss/instance
cp /app/lss/bin/VRVIUAA /app/lss/instance
cp /app/lss/bin/VRVIUBB /app/lss/instance
cp /app/lss/bin/VRVIUCC /app/lss/instance
cp /app/lss/bin/VRVIUDD /app/lss/instance

mkdir -p /app/lss/instance0 /app/stream_dispatch/logs
mkdir -p /etc/services.d/lss0
cp /home/lssforfz.run /etc/services.d/lss0/run
cp -r /app/lss/instance/* /app/lss/instance0
chmod +x  /etc/services.d/lss0/*
ln -s /usr/lib/nvidia-390/libnvcuvid.so.390.12 /usr/local/lib/libnvcuvid.so
ln -s /usr/lib/nvidia-390/libnvcuvid.so.390.12 /usr/local/lib/libnvcuvid.so.1
ln -s /usr/lib/nvidia-390/libnvidia-encode.so /usr/local/lib/libnvidia-encode.so
ln -s /usr/lib/nvidia-390/libnvidia-encode.so /usr/local/lib/libnvidia-encode.so.1
ln -s /usr/lib/nvidia-390/libnvidia-encode.so /usr/local/lib/libnvidia-encode.so.390.12
ln -s /usr/lib/nvidia-390/libEGL_nvidia.so.390.12 /usr/local/lib/libEGL_nvidia.so

mkdir -p /etc/services.d/stream_dispatch
cp /home/stream_dispatch.run /etc/services.d/stream_dispatch/run
chmod +x /etc/services.d/stream_dispatch/*
