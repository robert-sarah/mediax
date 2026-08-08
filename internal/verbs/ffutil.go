package verbs

import (
	"os/exec"
	"strconv"
	"strings"
)

// getMediaDurationSeconds returns the duration (in seconds) of a media file
// by querying ffprobe directly. It requires ffprobe to be installed and
// available on PATH, same as ffmpeg.
func getMediaDurationSeconds(path string) (float64, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(out))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}
