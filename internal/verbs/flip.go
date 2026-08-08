package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// Flip flips video horizontally or vertically
type Flip struct{}

func (f *Flip) Name() string {
	return "flip"
}

func (f *Flip) Description() string {
	return "Flip video horizontally or vertically (mirror)"
}

func (f *Flip) Usage() string {
	return "mediax flip <input> <output> [--horizontal | --vertical]"
}

func (f *Flip) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", f.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Determine flip type
	horizontal := flags["horizontal"] != "" || flags["h"] != ""
	vertical := flags["vertical"] != "" || flags["v"] != ""

	var vf string
	if horizontal {
		vf = "hflip"
	} else if vertical {
		vf = "vflip"
	} else {
		// Default to horizontal flip
		vf = "hflip"
	}

	cmdArgs := []string{
		"-i", input,
		"-vf", vf,
		"-c:a", "copy",
		"-y",
		output,
	}

	flipType := "horizontal"
	if vertical {
		flipType = "vertical"
	}

	fmt.Printf("Flipping %s: %s -> %s\n", flipType, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Flip{})
}
