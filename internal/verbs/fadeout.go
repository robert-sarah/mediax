package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// FadeOut applies fade-out effect at the end
type FadeOut struct{}

func (f *FadeOut) Name() string {
	return "fade-out"
}

func (f *FadeOut) Description() string {
	return "Apply fade-out effect at the end"
}

func (f *FadeOut) Usage() string {
	return "mediax fade-out <input> <output> [--duration 2s] [--start auto | --start 00:01:00]"
}

func (f *FadeOut) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", f.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse duration
	duration := flags["duration"]
	if duration == "" {
		duration = "2"
	}
	if len(duration) > 0 && duration[len(duration)-1] == 's' {
		duration = duration[:len(duration)-1]
	}

	// Determine real media duration via ffprobe so the fade starts
	// exactly `duration` seconds before the end of the file.
	totalDuration, err := getMediaDurationSeconds(input)
	if err != nil {
		return fmt.Errorf("failed to read media duration: %w", err)
	}

	fadeDuration, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		return fmt.Errorf("invalid duration value: %s", duration)
	}

	startAt := totalDuration - fadeDuration
	if startAt < 0 {
		startAt = 0
	}
	startStr := strconv.FormatFloat(startAt, 'f', 3, 64)

	// Build video fade filter that fades out at the end
	// fade=t=out:st=start_time:d=duration
	vf := fmt.Sprintf("fade=t=out:st=%s:d=%s", startStr, duration)

	// Audio fade filter
	af := fmt.Sprintf("afade=t=out:st=%s:d=%s", startStr, duration)

	cmdArgs := []string{
		"-i", input,
		"-vf", vf,
		"-af", af,
		"-y",
		output,
	}

	fmt.Printf("Applying fade-out (duration: %ss) to: %s -> %s\n", duration, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&FadeOut{})
}
