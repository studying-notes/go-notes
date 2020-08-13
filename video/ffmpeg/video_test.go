package ffmpeg

import "testing"

// go test -run TestExtractVideo
func TestExtractVideo(t *testing.T) {
	err = ExtractVideo("mp4/merge.mp4")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestExtractAudio
func TestExtractAudio(t *testing.T) {
	err = ExtractAudio("mp4/merge.mp4")
	if err != nil {
		t.Error(err)
	}
}

// 拼接多个 MP4 视频文件
// go test -run TestConcatVideos
func TestConcatVideos(t *testing.T) {
	err = ConcatVideos("mp4/merge.mp4", "mp4/1.mp4",
		"mp4/2.mp4", "mp4/3.mp4", "mp4/4.mp4", "mp4/5.mp4")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestCutVideo
func TestCutVideo(t *testing.T) {
	err = CutVideo("mp4/long.mp4",
		"mp4/long_cut.mp4", 2, 5)
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestTransposeVideo
func TestTransposeVideo(t *testing.T) {
	err = TransposeVideo("mp4/long.mp4", "mp4/long_90.mp4", 90)
	if err != nil {
		t.Error(err)
	}
}
