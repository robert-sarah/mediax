package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Thumbnail extracts a single frame at specific time
type Thumbnail struct{}

func (t *Thumbnail) Name() string {
	return "thumbnail"
}

func (t *Thumbnail) Description() string {
	return "Extract thumbnail image from video at specific time"
}

func (t *Thumbnail) Usage() string {
	return "mediax thumbnail <input> <output.png> [--time 00:00:01] [--width 1920]"
}

func (t *Thumbnail) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", t.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse time
	captureTime := flags["time"]
	if captureTime == "" {
		captureTime = flags["at"]
	}
	if captureTime == "" {
		captureTime = "00:00:01" // Default to 1 second
	}

	// Parse width
	width := flags["width"]
	if width == "" {
		width = "1920"
	}

	// Determine output format from extension
	ext := filepath.Ext(output)
	if ext == "" {
		ext = ".png"
	}

	// Build ffmpeg command
	vf := fmt.Sprintf("select='eq(n,0)',scale=%s:-1", width)

	cmdArgs := []string{
		"-ss", captureTime,
		"-i", input,
		"-vf", vf,
		"-vframes", "1",
		"-update", "1",
		"-y",
		output,
	}

	fmt.Printf("Extracting thumbnail at %s: %s -> %s\n", captureTime, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Thumbnail{})
}
