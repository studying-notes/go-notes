package ffmpeg

import "testing"

// go test -run TestMix2Audios
func TestMix2Audios(t *testing.T) {
	err = Mix2Audios("mp4/1.mp3",
		"mp4/merge_audio.mp3", "mp4/mix_audio.mp3")
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestCutAudio
func TestCutAudio(t *testing.T) {
	err = CutAudio("mp4/long.mp3",
		"mp4/long_cut.mp3", 20, 15)
	if err != nil {
		t.Error(err)
	}
}

// go test -run TestAdjustVolumeMultiple
func TestAdjustVolumeMultiple(t *testing.T) {
	err = AdjustVolumeMultiple("mp4/long.mp3", "mp4/long_3vol.mp3", 0.3)
	if err != nil {
		t.Error(err)
	}
}
