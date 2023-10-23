#!/usr/bin/env bash

gst-launch-1.0 udpsrc port=5004 ! \
    application/x-rtp,encoding-name=H264 ! \
        rtpjitterbuffer ! \
            rtph264depay ! \
                avdec_h264 ! \
                    videoconvert ! \
                        ximagesink




