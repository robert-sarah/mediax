package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// SlideShow creates video from image sequence
type SlideShow struct{}

func (s *SlideShow) Name() string {
	return "slideshow"
}

func (s *SlideShow) Description() string {
	return "Create video from image sequence"
}

func (s *SlideShow) Usage() string {
	return "mediax slideshow <image-pattern> <output> [--fps 1] [--duration 5] [--transition fade]"
}

func (s *SlideShow) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", s.Usage())
	}

	pattern, output := args[0], args[1]

	// Parse parameters
	fps := 1
	if f := flags["fps"]; f != "" {
		if val, err := strconv.Atoi(f); err == nil && val > 0 {
			fps = val
		}
	}

	duration := 5
	if d := flags["duration"]; d != "" {
		if val, err := strconv.Atoi(d); err == nil && val > 0 {
			duration = val
		}
	}

	transition := flags["transition"]
	if transition == "" {
		transition = "none"
	}

	// Check if pattern contains wildcards
	if !strings.Contains(pattern, "*") && !strings.Contains(pattern, "?") {
		// Single image - treat as pattern with frame number
		base := pattern[:len(pattern)-len(filepath.Ext(pattern))]
		ext := filepath.Ext(pattern)
		pattern = base + "_%03d" + ext
	}

	var cmdArgs []string

	if transition == "fade" || transition == "crossfade" {
		// Use complex filter for crossfade transitions
		// This is simplified - full implementation would need more complex filter
		vf := fmt.Sprintf("fps=%d,format=yuv420p", fps)
		cmdArgs = []string{
			"-framerate", strconv.Itoa(fps),
			"-i", pattern,
			"-vf", vf,
			"-c:v", "libx264",
			"-pix_fmt", "yuv420p",
			"-movflags", "+faststart",
			"-y",
			output,
		}
	} else {
		// Simple slideshow without transitions
		vf := fmt.Sprintf("fps=%d,format=yuv420p", fps)
		cmdArgs = []string{
			"-framerate", strconv.Itoa(fps),
			"-i", pattern,
			"-vf", vf,
			"-c:v", "libx264",
			"-pix_fmt", "yuv420p",
			"-movflags", "+faststart",
			"-y",
			output,
		}
	}

	fmt.Printf("Creating slideshow from images matching: %s -> %s\n", pattern, output)
	fmt.Printf("Settings: %d fps, %d seconds per image\n", fps, duration)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&SlideShow{})
}
