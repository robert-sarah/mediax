package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Blur applies blur effect to video
type Blur struct{}

func (b *Blur) Name() string {
	return "blur"
}

func (b *Blur) Description() string {
	return "Apply blur effect to video"
}

func (b *Blur) Usage() string {
	return "mediax blur <input> <output> [--strength 10] [--type gaussian|box]"
}

func (b *Blur) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", b.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse blur strength
	strength := 10
	if s := flags["strength"]; s != "" {
		if val, err := strconv.Atoi(s); err == nil && val > 0 {
			strength = val
		}
	}

	// Cap strength to reasonable values
	if strength > 100 {
		strength = 100
	}

	// Determine blur type
	blurType := flags["type"]
	if blurType == "" {
		blurType = "gaussian"
	}

	// Build filter
	var vf string
	switch blurType {
	case "box":
		// Box blur (faster, softer)
		vf = fmt.Sprintf("boxblur=%d:%d", strength, strength)
	case "gaussian":
		// Gaussian blur using gblur filter
		// Sigma controls blur amount
		sigma := float64(strength) / 10.0
		vf = fmt.Sprintf("gblur=sigma=%.1f", sigma)
	default:
		vf = fmt.Sprintf("boxblur=%d:%d", strength, strength)
	}

	cmdArgs := []string{
		"-i", input,
		"-vf", vf,
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Applying %s blur (strength: %d) to: %s -> %s\n", blurType, strength, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Blur{})
}
