
THIS_PATH = $(shell pwd)
FFMPEG_BUILD_PATH = $(THIS_PATH)/ffmpeg_build/
CPUNUM = $(shell nproc)

#echo $(FFMPEG_BUILD_PATH)
all:ffmpeg
	
ffmpeg:
	cd ffmpeg-4.0.2 && \
	./configure --prefix=$(FFMPEG_BUILD_PATH) --enable-gpl --enable-nonfree --enable-libass \
	--enable-libfdk-aac --enable-libfreetype --enable-libmp3lame \
	--enable-libx264 --enable-libx265  --enable-cuda --enable-cuvid --enable-nvenc --enable-nonfree --enable-libnpp \
	--disable-shared  --enable-static --extra-cflags=-I/usr/local/cuda/include --extra-ldflags=-L/usr/local/cuda/lib64 && \
	make -j$(CPUNUM) && \
	make install

clean:
	-cd ffmpeg-4.0.2 && \
	make clean > /dev/null 2>&1

