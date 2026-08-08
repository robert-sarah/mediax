package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Split divides video into multiple segments
type Split struct{}

func (s *Split) Name() string {
	return "split"
}

func (s *Split) Description() string {
	return "Split video into multiple segments"
}

func (s *Split) Usage() string {
	return "mediax split <input> [--duration 60 | --parts 4] [--output-prefix part]"
}

func (s *Split) Execute(args []string, flags map[string]string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s", s.Usage())
	}

	input := args[0]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse duration or parts
	duration := flags["duration"]
	parts := flags["parts"]

	if duration == "" && parts == "" {
		duration = "60" // default 60 seconds per segment
	}

	// Output prefix
	prefix := flags["output-prefix"]
	if prefix == "" {
		base := input[0 : len(input)-len(filepath.Ext(input))]
		prefix = base + "_part"
	}

	var cmdArgs []string

	if duration != "" {
		// Split by duration
		cmdArgs = []string{
			"-i", input,
			"-c", "copy",
			"-f", "segment",
			"-segment_time", duration,
			"-reset_timestamps", "1",
			fmt.Sprintf("%s%%03d%s", prefix, filepath.Ext(input)),
		}
	} else if parts != "" {
		// Split into N roughly equal parts: compute the real duration of
		// the file via ffprobe and derive a per-segment time for ffmpeg's
		// segment muxer (which has no native "split into N parts" option).
		p, _ := strconv.Atoi(parts)
		if p < 2 {
			p = 2
		}

		totalDuration, err := getMediaDurationSeconds(input)
		if err != nil {
			return fmt.Errorf("failed to read media duration: %w", err)
		}

		segmentTime := totalDuration / float64(p)
		if segmentTime <= 0 {
			return fmt.Errorf("could not compute a valid segment duration")
		}
		segmentTimeStr := strconv.FormatFloat(segmentTime, 'f', 3, 64)

		cmdArgs = []string{
			"-i", input,
			"-c", "copy",
			"-f", "segment",
			"-segment_time", segmentTimeStr,
			"-reset_timestamps", "1",
			fmt.Sprintf("%s%%03d%s", prefix, filepath.Ext(input)),
		}
	}

	fmt.Printf("Splitting %s into segments...\n", input)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Split{})
}
