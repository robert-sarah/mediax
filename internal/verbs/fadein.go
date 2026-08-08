package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// FadeIn applies fade-in effect to video and/or audio
type FadeIn struct{}

func (f *FadeIn) Name() string {
	return "fade-in"
}

func (f *FadeIn) Description() string {
	return "Apply fade-in effect at the beginning"
}

func (f *FadeIn) Usage() string {
	return "mediax fade-in <input> <output> [--duration 2s] [--video-only | --audio-only]"
}

func (f *FadeIn) Execute(args []string, flags map[string]string) error {
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
	// Remove 's' suffix if present
	if len(duration) > 0 && duration[len(duration)-1] == 's' {
		duration = duration[:len(duration)-1]
	}

	// Check for video-only or audio-only
	videoOnly := flags["video-only"] != "" || flags["v"] != ""
	audioOnly := flags["audio-only"] != "" || flags["a"] != ""

	// Build video fade filter
	vf := fmt.Sprintf("fade=t=in:st=0:d=%s", duration)

	// Build audio fade filter
	af := fmt.Sprintf("afade=t=in:st=0:d=%s", duration)

	cmdArgs := []string{"-i", input}

	if videoOnly {
		// Video fade only
		cmdArgs = append(cmdArgs, "-vf", vf, "-c:a", "copy")
	} else if audioOnly {
		// Audio fade only
		cmdArgs = append(cmdArgs, "-af", af, "-c:v", "copy")
	} else {
		// Both video and audio fade
		cmdArgs = append(cmdArgs, "-vf", vf, "-af", af)
	}

	cmdArgs = append(cmdArgs, "-y", output)

	fmt.Printf("Applying fade-in (duration: %ss) to: %s -> %s\n", duration, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&FadeIn{})
}
