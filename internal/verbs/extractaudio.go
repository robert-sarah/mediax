package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ExtractAudio extracts audio track from video
type ExtractAudio struct{}

func (e *ExtractAudio) Name() string {
	return "extract-audio"
}

func (e *ExtractAudio) Description() string {
	return "Extract audio track from video (full or partial)"
}

func (e *ExtractAudio) Usage() string {
	return "mediax extract-audio <video> <output> [--start HH:MM:SS] [--duration HH:MM:SS] [--time HH:MM:SS-HH:MM:SS] [--format mp3|aac|wav|flac]"
}

func (e *ExtractAudio) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", e.Usage())
	}

	input := args[0]
	output := args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
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

	// Determine output format
	format := flags["format"]
	if format == "" {
		format = "mp3"
	}

	// Build ffmpeg command
	cmdArgs := []string{"-i", input, "-vn", "-y"}

	// Add time flags if specified
	if start != "" {
		cmdArgs = append(cmdArgs, "-ss", start)
	}
	if duration != "" {
		cmdArgs = append(cmdArgs, "-t", duration)
	}

	// Set codec based on format
	switch format {
	case "mp3":
		cmdArgs = append(cmdArgs, "-c:a", "libmp3lame", "-q:a", "2")
	case "aac":
		cmdArgs = append(cmdArgs, "-c:a", "aac", "-b:a", "192k")
	case "wav":
		cmdArgs = append(cmdArgs, "-c:a", "pcm_s16le")
	case "flac":
		cmdArgs = append(cmdArgs, "-c:a", "flac")
	default:
		cmdArgs = append(cmdArgs, "-c:a", "copy")
	}

	cmdArgs = append(cmdArgs, output)

	if start != "" || duration != "" {
		fmt.Printf("Extracting partial audio from %s to %s: %s -> %s\n", start, duration, input, output)
	} else {
		fmt.Printf("Extracting full audio: %s -> %s\n", input, output)
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func calculateDuration(start, end string) string {
	// Parse HH:MM:SS format and calculate duration in seconds
	startSec := parseTimeToSeconds(start)
	endSec := parseTimeToSeconds(end)
	
	if startSec >= 0 && endSec >= 0 && endSec > startSec {
		duration := endSec - startSec
		return formatSecondsToTime(duration)
	}
	
	// Fallback: return the duration as-is if it's already a number
	return end
}

func parseTimeToSeconds(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	if len(parts) == 3 {
		h, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		s, _ := strconv.Atoi(parts[2])
		return h*3600 + m*60 + s
	}
	if len(parts) == 2 {
		m, _ := strconv.Atoi(parts[0])
		s, _ := strconv.Atoi(parts[1])
		return m*60 + s
	}
	// Try parsing as plain seconds
	s, _ := strconv.Atoi(timeStr)
	return s
}

func formatSecondsToTime(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func init() {
	Register(&ExtractAudio{})
}
