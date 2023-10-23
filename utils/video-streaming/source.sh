#!/usr/bin/env bash

gst-launch-1.0 nvarguscamerasrc sensor-id=0 sensor-mode=-1 bufapi-version=true ! \
    'video/x-raw(memory:NVMM),framerate=30/1,width=640,height=360' ! \
        nvvideoconvert gpu-id=0 ! \
            video/x-raw,format=I420,width=640,height=360 ! \
                x264enc tune=zerolatency ! \
                    rtph264pay ! \
                        udpsink host=192.168.24.13 port=5004

