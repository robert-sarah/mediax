package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// Compress reduces file size while maintaining quality
type Compress struct{}

func (c *Compress) Name() string {
	return "compress"
}

func (c *Compress) Description() string {
	return "Compress video to reduce file size"
}

func (c *Compress) Usage() string {
	return "mediax compress <input> <output> [--quality high|medium|low] [--target 10mb]"
}

func (c *Compress) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", c.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	quality := flags["quality"]
	if quality == "" {
		quality = "medium"
	}

	cmdArgs := []string{"-i", input, "-y"}

	switch quality {
	case "high", "best":
		// Good quality, reasonable size
		cmdArgs = append(cmdArgs, "-c:v", "libx264", "-preset", "slow", "-crf", "20", "-c:a", "aac", "-b:a", "192k")
	case "medium", "normal":
		// Balanced
		cmdArgs = append(cmdArgs, "-c:v", "libx264", "-preset", "medium", "-crf", "23", "-c:a", "aac", "-b:a", "128k")
	case "low", "fast":
		// Smallest size
		cmdArgs = append(cmdArgs, "-c:v", "libx264", "-preset", "veryfast", "-crf", "28", "-c:a", "aac", "-b:a", "96k")
	}

	cmdArgs = append(cmdArgs, output)

	fmt.Printf("Compressing %s -> %s (quality: %s)\n", input, output, quality)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Compress{})
}
