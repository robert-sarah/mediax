package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Rotate rotates video by specified degrees
type Rotate struct{}

func (r *Rotate) Name() string {
	return "rotate"
}

func (r *Rotate) Description() string {
	return "Rotate video by 90, 180, 270 degrees or custom angle"
}

func (r *Rotate) Usage() string {
	return "mediax rotate <input> <output> [--degrees 90] [--flip]"
}

func (r *Rotate) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", r.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Get rotation degrees
	degrees := 90
	if d := flags["degrees"]; d != "" {
		if val, err := strconv.Atoi(d); err == nil {
			degrees = val
		}
	}

	// Build transpose filter for 90 degree rotations
	var vf string
	switch degrees {
	case 90:
		vf = "transpose=1" // 90 degrees clockwise
	case 180:
		vf = "transpose=2,transpose=2" // 180 degrees
	case 270, -90:
		vf = "transpose=2" // 90 degrees counter-clockwise
	default:
		// Custom angle using rotate filter
		vf = fmt.Sprintf("rotate=%d*PI/180", degrees)
	}

	cmdArgs := []string{
		"-i", input,
		"-vf", vf,
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Rotating %s by %d degrees -> %s\n", input, degrees, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Rotate{})
}
