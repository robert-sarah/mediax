package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Concat joins multiple videos together
type Concat struct{}

func (c *Concat) Name() string {
	return "concat"
}

func (c *Concat) Description() string {
	return "Join multiple videos into one"
}

func (c *Concat) Usage() string {
	return "mediax concat <output> <input1> <input2> [input3...]"
}

func (c *Concat) Execute(args []string, flags map[string]string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: %s", c.Usage())
	}

	output := args[0]
	inputs := args[1:]

	// Verify all input files exist
	for _, input := range inputs {
		if _, err := os.Stat(input); err != nil {
			return fmt.Errorf("input file not found: %s", input)
		}
	}

	// Create temporary file list for concat demuxer
	tmpFile, err := os.CreateTemp("", "mediax-concat-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write file list
	for _, input := range inputs {
		// Get absolute path
		absPath, err := filepath.Abs(input)
		if err != nil {
			absPath = input
		}
		// Escape single quotes in path
		escaped := strings.ReplaceAll(absPath, "'", "'\\''")
		fmt.Fprintf(tmpFile, "file '%s'\n", escaped)
	}
	tmpFile.Close()

	// Build ffmpeg command using concat demuxer
	cmdArgs := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", tmpFile.Name(),
		"-c", "copy",
		"-y",
		output,
	}

	fmt.Printf("Concatenating %d videos into: %s\n", len(inputs), output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Concat{})
}
