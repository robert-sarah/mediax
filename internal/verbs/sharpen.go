package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Sharpen applies sharpening filter to enhance edges
type Sharpen struct{}

func (s *Sharpen) Name() string {
	return "sharpen"
}

func (s *Sharpen) Description() string {
	return "Sharpen video to enhance edges and details"
}

func (s *Sharpen) Usage() string {
	return "mediax sharpen <input> <output> [--amount 1.5]"
}

func (s *Sharpen) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", s.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse sharpen amount
	amount := 1.5
	if a := flags["amount"]; a != "" {
		if val, err := strconv.ParseFloat(a, 64); err == nil {
			amount = val
		}
	}

	// Build unsharp filter: luma_msize_x:luma_msize_y:luma_amount
	// Default: 5:5:1.0, we scale the amount
	vf := fmt.Sprintf("unsharp=5:5:%.2f", amount)

	cmdArgs := []string{
		"-i", input,
		"-vf", vf,
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Sharpening video (amount: %.2f): %s -> %s\n", amount, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Sharpen{})
}
