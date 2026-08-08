package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Gif creates GIF animation from video
type Gif struct{}

func (g *Gif) Name() string {
	return "gif"
}

func (g *Gif) Description() string {
	return "Create GIF animation from video"
}

func (g *Gif) Usage() string {
	return "mediax gif <input> <output.gif> [--fps 10] [--width 480] [--start 0] [--duration 5]"
}

func (g *Gif) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", g.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse parameters with defaults
	fps, _ := strconv.Atoi(flags["fps"])
	if fps <= 0 {
		fps = 15
	}

	width, _ := strconv.Atoi(flags["width"])
	if width <= 0 {
		width = 480
	}

	start := flags["start"]
	duration := flags["duration"]

	// Build filter chain for high-quality GIF
	// fps -> scale -> split -> palettegen -> paletteuse
	filter := fmt.Sprintf("fps=%d,scale=%d:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse=dither=bayer",
		fps, width)

	cmdArgs := []string{"-i", input}

	// Add start time if specified
	if start != "" && start != "0" {
		cmdArgs = append(cmdArgs, "-ss", start)
	}

	// Add duration if specified
	if duration != "" {
		cmdArgs = append(cmdArgs, "-t", duration)
	}

	// Optimization: use single thread for GIF encoding
	cmdArgs = append(cmdArgs,
		"-vf", filter,
		"-loop", "0",
		"-y",
		output,
	)

	fmt.Printf("Creating GIF from %s -> %s (%dfps, %dpx wide)\n", input, output, fps, width)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Gif{})
}
