package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Chroma applies chroma key (green screen) effect
type Chroma struct{}

func (c *Chroma) Name() string {
	return "chroma"
}

func (c *Chroma) Description() string {
	return "Apply chroma key (green screen) effect"
}

func (c *Chroma) Usage() string {
	return "mediax chroma <input> <background> <output> [--color 00FF00] [--similarity 0.3] [--blend 0.1]"
}

func (c *Chroma) Execute(args []string, flags map[string]string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: %s", c.Usage())
	}

	input, background, output := args[0], args[1], args[2]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}
	if _, err := os.Stat(background); err != nil {
		return fmt.Errorf("background file not found: %s", background)
	}

	// Parse parameters
	color := flags["color"]
	if color == "" {
		color = "00FF00" // Default green
	}
	// Remove # if present
	color = strings.TrimPrefix(color, "#")
	// Convert to 0x format for ffmpeg
	color = "0x" + color

	similarity := 0.3
	if s := flags["similarity"]; s != "" {
		if val, err := strconv.ParseFloat(s, 64); err == nil {
			similarity = val
		}
	}

	blend := 0.1
	if b := flags["blend"]; b != "" {
		if val, err := strconv.ParseFloat(b, 64); err == nil {
			blend = val
		}
	}

	// Build chromakey filter
	// chromakey=color:similarity:blend
	chromakey := fmt.Sprintf("chromakey=%s:%.2f:%.2f", color, similarity, blend)

	// Complex filter: overlay chroma-keyed video over background
	// [0:v]chromakey[fg];[1:v][fg]overlay=0:0
	filterComplex := fmt.Sprintf("[0:v]%s[fg];[1:v][fg]overlay=0:0:shortest=1", chromakey)

	cmdArgs := []string{
		"-i", input,      // Foreground (with green screen)
		"-i", background, // Background
		"-filter_complex", filterComplex,
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Applying chroma key (%s): %s + %s -> %s\n", color, input, background, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Chroma{})
}
