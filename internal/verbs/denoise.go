package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Denoise reduces noise in video using various filters
type Denoise struct{}

func (d *Denoise) Name() string {
	return "denoise"
}

func (d *Denoise) Description() string {
	return "Reduce noise in video"
}

func (d *Denoise) Usage() string {
	return "mediax denoise <input> <output> [--strength medium] [--type nlmeans|hqdn3d]"
}

func (d *Denoise) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", d.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse denoise type
	denoiseType := flags["type"]
	if denoiseType == "" {
		denoiseType = "hqdn3d"
	}

	// Parse strength
	strength := flags["strength"]
	if strength == "" {
		strength = "medium"
	}

	var vf string

	switch denoiseType {
	case "nlmeans":
		// Non-local means denoising (high quality, slower)
		sigma, _ := strconv.ParseFloat(strength, 64)
		if sigma <= 0 {
			sigma = 10.0
		}
		vf = fmt.Sprintf("nlmeans=s=%.1f", sigma)

	case "hqdn3d", "default":
		// High quality 3D denoiser (faster)
		var lumaSp, chromaSp, lumaTmp, chromaTmp int
		switch strength {
		case "low", "light":
			lumaSp, chromaSp, lumaTmp, chromaTmp = 2, 2, 3, 3
		case "medium", "normal":
			lumaSp, chromaSp, lumaTmp, chromaTmp = 4, 4, 6, 6
		case "high", "strong":
			lumaSp, chromaSp, lumaTmp, chromaTmp = 8, 8, 12, 12
		default:
			lumaSp, chromaSp, lumaTmp, chromaTmp = 4, 4, 6, 6
		}
		vf = fmt.Sprintf("hqdn3d=%d:%d:%d:%d", lumaSp, chromaSp, lumaTmp, chromaTmp)
	}

	cmdArgs := []string{
		"-i", input,
		"-vf", vf,
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Denoising video (%s, strength: %s): %s -> %s\n", denoiseType, strength, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Denoise{})
}
