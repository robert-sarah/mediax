package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Convert converts media files between formats
type Convert struct{}

func (c *Convert) Name() string {
	return "convert"
}

func (c *Convert) Description() string {
	return "Convert media to different format"
}

func (c *Convert) Usage() string {
	return "mediax convert <input> <output> [--codec h264] [--quality high]"
}

func (c *Convert) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: mediax convert <input> <output>")
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Determine output format from extension
	outputExt := strings.ToLower(filepath.Ext(output))

	// Build ffmpeg command
	cmdArgs := []string{"-i", input, "-y"}

	// Apply codec if specified
	if codec := flags["codec"]; codec != "" {
		cmdArgs = append(cmdArgs, "-c:v", codec)
	}

	// Apply quality preset
	if quality := flags["quality"]; quality != "" {
		switch quality {
		case "high", "best":
			cmdArgs = append(cmdArgs, "-preset", "slow", "-crf", "18")
		case "medium", "normal":
			cmdArgs = append(cmdArgs, "-preset", "medium", "-crf", "23")
		case "low", "fast":
			cmdArgs = append(cmdArgs, "-preset", "fast", "-crf", "28")
		}
	}

	// Format-specific options
	switch outputExt {
	case ".gif":
		cmdArgs = append(cmdArgs, "-vf", "fps=30,scale=480:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=256[p];[s1][p]paletteuse=dither=bayer")
	case ".mp3", ".aac", ".ogg", ".flac":
		// Audio-only output
		cmdArgs = append(cmdArgs, "-vn")
	}

	cmdArgs = append(cmdArgs, output)

	fmt.Printf("Converting: %s -> %s\n", input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Convert{})
}
