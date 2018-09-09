#!/bin/bash

#step1 install compile tools
apt install automake fdk-aac git lame libass libtool libvorbis libvpx \
opus sdl shtool texi2html theora wget x264 x265 xvid nasm

cd nv-codec-headers/
make  && make install
cd ../

#step2 clean the exist lib/include
#rm -rf ffmpeg_build/lib/*
#rm -rf ffmpeg_build/include/*
#rm -rf ffmpeg_build/ffmpeg.so.tar.gz

#step3 make
chmod +x -R ffmpeg-4.0.2/*
make
#
##step4 pack
cd ffmpeg_build
#
tar -cvzf ffmpeg.so.tar.gz lib/*
cp ffmpeg.so.tar.gz ../../deploy/ffmpeg/
cp bin/ffmpeg ../../deploy/ffmpeg/
cp bin/ffprobe ../../deploy/ffmpeg/
cd ../
