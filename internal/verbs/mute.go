package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// Mute removes audio track from video
type Mute struct{}

func (m *Mute) Name() string {
	return "mute"
}

func (m *Mute) Description() string {
	return "Remove audio track from video"
}

func (m *Mute) Usage() string {
	return "mediax mute <input> <output>"
}

func (m *Mute) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", m.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	cmdArgs := []string{
		"-i", input,
		"-an", // Remove audio
		"-c:v", "copy", // Copy video without re-encoding
		"-y",
		output,
	}

	fmt.Printf("Removing audio: %s -> %s\n", input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Mute{})
}
