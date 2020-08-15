package ffmpeg

import "testing"

var err error

// 视频与背景音乐混合
// go test -run TestMixBackgroundMusic
func TestMixBackgroundMusic(t *testing.T) {
	err = MixBackgroundMusic("mp4/merge.mp4",
		"mp4/1.mp3", "mp4/merge_bgm.mp4")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestAudioFormatConv
func TestAudioFormatConv(t *testing.T) {
	err = AudioFormatConv("mp4/merge_audio.m4a",
		"mp4/merge_audio.mp3")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestVideoConvAudio
func TestVideoConvAudio(t *testing.T) {
	err = VideoConvAudio("mp4/merge.mp4",
		"mp4/merge.mp3")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestOverlayOriginAudio
func TestOverlayOriginAudio(t *testing.T) {
	err = OverlayOriginAudio("mp4/merge.mp4",
		"mp4/merge.mp3", "mp4/merge_overlay.mp4")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestReplaceOriginAudio
func TestReplaceOriginAudio(t *testing.T) {
	err = ReplaceOriginAudio("mp4/merge.mp4",
		"mp4/merge.mp3", "mp4/merge_overlay.mp4")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestMergeVideoAudio
func TestMergeVideoAudio(t *testing.T) {
	err = MergeVideoAudio("mp4/merge_noaudio.mp4",
		"mp4/mix_audio.mp3", "mp4/merge_bgm.mp4")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestExtractVideoAudio
func TestExtractVideoAudio(t *testing.T) {
	err = ExtractVideoAudio("mp4/merge.mp4")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestShowStream
func TestShowStream(t *testing.T) {
	err = ExtractVideoAudio("mp4/merge.mp4")
	if err != nil {
		t.Error(err)
	}
}
