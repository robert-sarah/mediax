package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Crop crops video to specified rectangle
type Crop struct{}

func (c *Crop) Name() string {
	return "crop"
}

func (c *Crop) Description() string {
	return "Crop video to specified rectangle"
}

func (c *Crop) Usage() string {
	return "mediax crop <input> <output> --x 100 --y 100 --width 800 --height 600"
}

func (c *Crop) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", c.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse crop parameters
	x, _ := strconv.Atoi(flags["x"])
	if x < 0 {
		x = 0
	}

	y, _ := strconv.Atoi(flags["y"])
	if y < 0 {
		y = 0
	}

	width, _ := strconv.Atoi(flags["width"])
	if width <= 0 {
		return fmt.Errorf("width must be positive")
	}

	height, _ := strconv.Atoi(flags["height"])
	if height <= 0 {
		return fmt.Errorf("height must be positive")
	}

	// Build crop filter: crop=w:h:x:y
	cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d", width, height, x, y)

	cmdArgs := []string{
		"-i", input,
		"-vf", cropFilter,
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Cropping %s: %s\n", input, cropFilter)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Crop{})
}
