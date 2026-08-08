package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Volume adjusts audio volume
type Volume struct{}

func (v *Volume) Name() string {
	return "volume"
}

func (v *Volume) Description() string {
	return "Adjust audio volume (increase, decrease, normalize)"
}

func (v *Volume) Usage() string {
	return "mediax volume <input> <output> [--level 2.0 | --increase 5dB | --normalize]"
}

func (v *Volume) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", v.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	var audioFilter string

	// Check for normalize
	if _, ok := flags["normalize"]; ok {
		// Use loudnorm filter for EBU R128 normalization
		audioFilter = "loudnorm=I=-16:TP=-1.5:LRA=11"
	} else if level := flags["level"]; level != "" {
		// Multiplier (2.0 = 200%, 0.5 = 50%)
		audioFilter = fmt.Sprintf("volume=%s", level)
	} else if increase := flags["increase"]; increase != "" {
		// Decibel change
		if !strings.HasSuffix(increase, "dB") && !strings.HasSuffix(increase, "db") {
			increase += "dB"
		}
		audioFilter = fmt.Sprintf("volume=%s", increase)
	} else {
		// Default: boost volume by 50%
		audioFilter = "volume=1.5"
	}

	cmdArgs := []string{
		"-i", input,
		"-af", audioFilter,
		"-c:v", "copy", // Keep video unchanged
		"-y",
		output,
	}

	fmt.Printf("Adjusting volume: %s -> %s (%s)\n", input, output, audioFilter)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Volume{})
}
