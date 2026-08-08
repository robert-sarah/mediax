package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// Reverse plays video backwards
type Reverse struct{}

func (r *Reverse) Name() string {
	return "reverse"
}

func (r *Reverse) Description() string {
	return "Play video backwards (reverse)"
}

func (r *Reverse) Usage() string {
	return "mediax reverse <input> <output> [--audio]"
}

func (r *Reverse) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", r.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Build filter for video reverse
	vf := "reverse"

	// Check if audio should also be reversed
	reverseAudio := flags["audio"] != "" || flags["a"] != ""

	var af string
	if reverseAudio {
		af = "areverse"
	}

	cmdArgs := []string{"-i", input}

	// Add video filter
	cmdArgs = append(cmdArgs, "-vf", vf)

	// Add audio filter if requested
	if reverseAudio {
		cmdArgs = append(cmdArgs, "-af", af)
	} else {
		// Copy audio without reversing
		cmdArgs = append(cmdArgs, "-c:a", "copy")
	}

	cmdArgs = append(cmdArgs, "-y", output)

	if reverseAudio {
		fmt.Printf("Reversing video and audio: %s -> %s\n", input, output)
	} else {
		fmt.Printf("Reversing video (keeping audio): %s -> %s\n", input, output)
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Reverse{})
}
