# FFmpeg HTTP Server

## HLS

```
ffmpeg -f lavfi -re -i testsrc=size=1280x720:rate=25 -vf drawtext="fontfile=monofonto.ttf: fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/4): text='%{localtime\:%X}'" -c:v h264 -f hls -hls_time 4 -hls_segment_type fmp4 -method PUT "http://localhost:8080/upload/main.m3u8"
```

## LHLS x264

```
ffmpeg -re -i ~/videos/bigbuckbunny.mp4 -loglevel info   -map 0 -map 0 -map 0 -c:a aac -c:v libx264 -tune zerolatency   -b:v:0 2000k -s:v:0 1280x720 -profile:v:0 high -b:v:1 1500k -s:v:1 640x340  -profile:v:1 main -b:v:2 500k -s:v:2 320x170  -profile:v:2 baseline -bf 1  -keyint_min 24 -g 24 -sc_threshold 0 -b_strategy 0 -ar:a:1 22050 -use_timeline 1 -use_template 1 -window_size 5  -adaptation_sets "id=0,streams=v id=1,streams=a" -hls_playlist 1 -seg_duration 3 -streaming 1  -strict experimental -lhls 1 -remove_at_exit 0 -master_m3u8_publish_rate 3  -f dash -method PUT -http_persistent 1  http://localhost:8080/upload/manifest.mpd
```