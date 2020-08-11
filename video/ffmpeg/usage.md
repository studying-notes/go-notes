# ffmpeg 常用命令

- [ffmpeg 常用命令](#ffmpeg-常用命令)
  - [视频流和音频流](#视频流和音频流)
    - [显示流信息](#显示流信息)
    - [分离视频流](#分离视频流)
    - [分离音频流](#分离音频流)
  - [格式转换](#格式转换)
  - [音视频合成](#音视频合成)
  - [视频合并](#视频合并)
    - [将多个 MP4 文件合并为 1 个](#将多个-mp4-文件合并为-1-个)
  - [视频转码](#视频转码)
  - [视频剪切](#视频剪切)
  - [视频录制](#视频录制)

## 视频流和音频流

通常视频有两个要素：声音和画面。但其实严格意义上说视频中含有视频流和音频流，如果一个视频只有视频流，那么就只有画面没有声音，反之亦然。

### 显示流信息

```shell
ffmpeg -i 1.mp4
```

```
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from '1.mp4':

    # 视频流
    Stream #0:0(und): Video: h264 (Main) (avc1 / 0x31637661), yuv420p, 960x516 [SAR 1:1 DAR 80:43], 238 kb/s, 30 fps, 30 tbr, 15360 tbn, 60 tbc (default)

    # 音频流
    Stream #0:1(eng): Audio: aac (LC) (mp4a / 0x6134706D), 44100 Hz, stereo, fltp, 192 kb/s (default)
```

### 分离视频流

```shell
ffmpeg -i 1.mp4 -vcodec copy -an 1_video.mp4
```

```shell
ffmpeg -i 1_video.mp4
# 仅剩下视频流
```

### 分离音频流

```shell
ffmpeg -i 1.mp4 -acodec copy -vn 1_audio.m4a
```

视频中的音频是 ACC 格式，无法直接分离出 MP3 格式音频。必须进行转码：

```shell
ffmpeg -i 1_audio.m4a 1_audio.mp3
```

也可以直接将视频转换为 mp3 格式：

```shell
ffmpeg -i 1.mp4 1.mp3
```

## 格式转换

```shell
ffmpeg -i 1.mp4 1.avi
```

对于编码格式，一种方法是通过目标文件的扩展名来控制，另一种方法是通过 ``-c:v`` 参数来控制。

```shell
ffmpeg -i 1.mp4 -c:v libx265 1.avi
```

## 音视频合成

1. 去掉源文件里的音频

```shell
ffmpeg -i 1.mp4 -vcodec copy -an 1_an.mp4
```

- `-vcodec copy` 对源视频不解码，直接拷贝到目标文件
- `-an` 将源文件里的音频丢弃

2. 将这个视频文件与一个音频文件合成

```shell
ffmpeg -i 1_an.mp4 -ss 30 -t 52 -i 1.mp3 -vcodec copy 1_merge.mp4
```

`-ss 30 -t 52` 截取 1.mp3 文件的第 30 秒往后的 52 秒与视频合成。

## 视频合并

### 将多个 MP4 文件合并为 1 个

```shell
ffmpeg -i 1.mp4 -vcodec copy -acodec copy -vbsf h264_mp4toannexb 1.ts
ffmpeg -i 2.mp4 -vcodec copy -acodec copy -vbsf h264_mp4toannexb 2.ts

ffmpeg -i "concat:1.ts|2.ts" -acodec copy -vcodec copy -absf aac_adtstoasc 1merge2.mp4
```

## 视频转码

```sh
# 转码为码流原始文件
ffmpeg –i 1.mp4 –vcodec h264 –s 352*278 –an –f m4v 1.264

# 转码为码流原始文件
ffmpeg –i 1.mp4 –vcodec h264 –bf 0 –g 25 –s 352*278 –an –f m4v 1.264

# 转码为封装文件
ffmpeg –i 1.avi -vcodec mpeg4 –vtag xvid –qsame 1_xvid.avi
```

- `-bf` B 帧数目控制
- `-g` 关键帧间隔控制
- `-s` 分辨率控制

## 视频剪切

```
ffmpeg –i test.avi –r 1 –f image2 image-%3d.jpeg        //提取图片
ffmpeg -ss 0:1:30 -t 0:0:20 -i input.avi -vcodec copy -acodec copy output.avi    //剪切视频
```

- -r 提取图像的频率
- -ss 开始时间
- -t 持续时间

## 视频录制

```
ffmpeg –i rtsp://192.168.3.205:5555/test –vcodec copy out.avi
```
