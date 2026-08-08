package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ReplaceAudio replaces the audio track of a video
type ReplaceAudio struct{}

func (r *ReplaceAudio) Name() string {
	return "replace-audio"
}

func (r *ReplaceAudio) Description() string {
	return "Replace audio track in video (full or partial)"
}

func (r *ReplaceAudio) Usage() string {
	return "mediax replace-audio <video> <audio> <output> [--start HH:MM:SS] [--duration HH:MM:SS] [--time HH:MM:SS-HH:MM:SS]"
}

func (r *ReplaceAudio) Execute(args []string, flags map[string]string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: %s", r.Usage())
	}

	video, audio, output := args[0], args[1], args[2]

	if _, err := os.Stat(video); err != nil {
		return fmt.Errorf("video file not found: %s", video)
	}
	if _, err := os.Stat(audio); err != nil {
		return fmt.Errorf("audio file not found: %s", audio)
	}

	// Parse time flags
	start := flags["start"]
	duration := flags["duration"]
	timeRange := flags["time"]

	// If time range is provided, parse it
	if timeRange != "" {
		parts := strings.Split(timeRange, "-")
		if len(parts) == 2 {
			start = parts[0]
			duration = calculateDuration(parts[0], parts[1])
		}
	}

	// If no time flags, replace entire audio
	if start == "" && duration == "" && timeRange == "" {
		return r.replaceFullAudio(video, audio, output)
	}

	// Partial replacement using filter_complex
	return r.replacePartialAudio(video, audio, output, start, duration)
}

func (r *ReplaceAudio) replaceFullAudio(video, audio, output string) error {
	cmdArgs := []string{
		"-i", video,
		"-i", audio,
		"-c:v", "copy",
		"-map", "0:v:0",
		"-map", "1:a:0",
		"-shortest",
		"-y",
		output,
	}

	fmt.Printf("Replacing full audio: %s -> %s\n", video, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (r *ReplaceAudio) replacePartialAudio(video, audio, output, start, duration string) error {
	startSec := parseTimeToSeconds(start)
	durationSec := parseTimeToSeconds(duration)
	endSec := startSec + durationSec
	endTime := formatSecondsToTime(endSec)

	filter := fmt.Sprintf(
		"[0:a]atrim=0:%s,asetpts=PTS-STARTPTS[before];[0:a]atrim=%s:%s,asetpts=PTS-STARTPTS[after];[1:a]asetpts=PTS-STARTPTS[new];[before][new][after]concat=n=3:v=0:a=1[aout]",
		start, start, endTime,
	)

	cmdArgs := []string{
		"-i", video,
		"-i", audio,
		"-filter_complex", filter,
		"-c:v", "copy",
		"-map", "0:v",
		"-map", "[aout]",
		"-y",
		output,
	}

	fmt.Printf("Replacing partial audio from %s (duration %s): %s -> %s\n", start, duration, video, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&ReplaceAudio{})
}
