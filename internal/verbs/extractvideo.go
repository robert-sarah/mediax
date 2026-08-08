package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// ExtractVideo extracts video track without audio
type ExtractVideo struct{}

func (e *ExtractVideo) Name() string {
	return "extract-video"
}

func (e *ExtractVideo) Description() string {
	return "Extract video track without audio"
}

func (e *ExtractVideo) Usage() string {
	return "mediax extract-video <input> <output>"
}

func (e *ExtractVideo) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", e.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	cmdArgs := []string{
		"-i", input,
		"-an",       // No audio
		"-c:v", "copy", // Copy video
		"-y",
		output,
	}

	fmt.Printf("Extracting video: %s -> %s\n", input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&ExtractVideo{})
}
