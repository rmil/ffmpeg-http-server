# FFmpeg HTTP Server

## Endpoints

`/watch/` stream data, where you can feed your player with `manifest.mpd` or `master.m3u8`
`/publish/` FFmpeg HTTP output point
`/players/hls` HLS video player
`/players/dash` DASH video player

## HLS test src (single-quality)

```
ffmpeg -f lavfi -re -i testsrc=size=1280x720:rate=25 -vf drawtext="fontfile=monofonto.ttf: fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/4): text='%{localtime\:%X}'" -c:v h264 -f hls -hls_time 4 -hls_segment_type fmp4 -method PUT "http://localhost:8080/publish/main.m3u8"
```

## HLS test src (multi-quality)

```
ffmpeg -re -f lavfi  -i testsrc=size=1920x1080:rate=25 -vf drawtext="fontfile=monofonto.ttf: fontsize=96: box=1: boxcolor=black@0.75: boxborderw=5: fontcolor=white: x=(w-text_w)/2: y=((h-text_h)/2)+((h-text_h)/4): text='%{localtime\:%X}'" \
-map 0 -map 0 -map 0 -c:a aac -c:v libx264 \
-b:v:0 800k -s:v:0 1280x720 -profile:v:0 main -pix_fmt yuv420p \
-b:v:1 500k -s:v:1 640x340  -profile:v:1 main -pix_fmt yuv420p \
-b:v:2 300k -s:v:2 320x170  -profile:v:2 baseline \
-bf 1 \
-keyint_min 120 -g 120 -sc_threshold 0 -b_strategy 0 -ar:a:1 22050 -use_timeline 1 -use_template 1 \
-window_size 5 -adaptation_sets "id=0,streams=v id=1,streams=a" -hls_playlist 1 -seg_duration 4 -streaming 1 -remove_at_exit 1 -method PUT -f dash http://localhost:8080/publish/manifest.mpd
```

## HLS (multi-quality)

```
ffmpeg -re -stream_loop -1  -i ~/videos/bigbuckbunny.mp4 \
-map 0 -map 0 -map 0 -c:a aac -c:v libx264 \
-b:v:0 800k -s:v:0 1280x720 -profile:v:0 main \
-b:v:1 500k -s:v:1 640x340  -profile:v:1 main \
-b:v:2 300k -s:v:2 320x170  -profile:v:2 baseline \
-bf 1 \
-keyint_min 120 -g 120 -sc_threshold 0 -b_strategy 0 -ar:a:1 22050 -use_timeline 1 -use_template 1 \
-window_size 5 -adaptation_sets "id=0,streams=v id=1,streams=a" -hls_playlist 1 -seg_duration 4 -streaming 1 -remove_at_exit 1 -method PUT -f dash http://localhost:8080/publish/manifest.mpd
```

## LHLS x264 (multi-quality)

```
ffmpeg -re -i ~/videos/bigbuckbunny.mp4 -loglevel info   -map 0 -map 0 -map 0 -c:a aac -c:v libx264 -tune zerolatency   -b:v:0 2000k -s:v:0 1280x720 -profile:v:0 high -b:v:1 1500k -s:v:1 640x340  -profile:v:1 main -b:v:2 500k -s:v:2 320x170  -profile:v:2 baseline -bf 1  -keyint_min 24 -g 24 -sc_threshold 0 -b_strategy 0 -ar:a:1 22050 -use_timeline 1 -use_template 1 -window_size 5  -adaptation_sets "id=0,streams=v id=1,streams=a" -hls_playlist 1 -seg_duration 3 -streaming 1  -strict experimental -lhls 1 -remove_at_exit 0 -master_m3u8_publish_rate 3  -f dash -method PUT -http_persistent 1  http://localhost:8080/publish/manifest.mpd
```