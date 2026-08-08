package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Trim extracts a portion of media (cut)
type Trim struct{}

func (t *Trim) Name() string {
	return "trim"
}

func (t *Trim) Description() string {
	return "Extract a segment from media (cut)"
}

func (t *Trim) Usage() string {
	return "mediax trim <input> <output> --start 00:01:30 [--duration 10s | --end 00:02:00]"
}

func (t *Trim) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: mediax trim <input> <output>")
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	start := flags["start"]
	if start == "" {
		start = flags["ss"]
	}

	end := flags["end"]
	duration := flags["duration"]
	if duration == "" {
		duration = flags["t"]
	}

	if start == "" && end == "" && duration == "" {
		return fmt.Errorf("must specify --start, --end, or --duration")
	}

	// Build ffmpeg command
	cmdArgs := []string{"-i", input}

	if start != "" {
		cmdArgs = append(cmdArgs, "-ss", start)
	}

	if duration != "" {
		cmdArgs = append(cmdArgs, "-t", duration)
	} else if end != "" {
		// Calculate duration from start and end
		cmdArgs = append(cmdArgs, "-to", end)
	}

	// Use fast seeking if starting from beginning
	if start != "" && !strings.HasPrefix(start, "00:00") {
		// Reorder for fast seeking
		cmdArgs = []string{"-ss", start, "-i", input}
		if duration != "" {
			cmdArgs = append(cmdArgs, "-t", duration)
		} else if end != "" {
			cmdArgs = append(cmdArgs, "-to", end)
		}
	}

	// Copy codecs for fast trimming if no re-encoding needed
	cmdArgs = append(cmdArgs, "-c", "copy")
	cmdArgs = append(cmdArgs, "-y", output)

	fmt.Printf("Trimming: %s -> %s (from %s", input, output, start)
	if duration != "" {
		fmt.Printf(", duration %s)\n", duration)
	} else if end != "" {
		fmt.Printf(" to %s)\n", end)
	} else {
		fmt.Println(")")
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Trim{})
}
