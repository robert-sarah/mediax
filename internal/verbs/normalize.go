package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// Normalize applies EBU R128 audio normalization
type Normalize struct{}

func (n *Normalize) Name() string {
	return "normalize"
}

func (n *Normalize) Description() string {
	return "Normalize audio to EBU R128 standard"
}

func (n *Normalize) Usage() string {
	return "mediax normalize <input> <output> [--target -16] [--true-peak -1.5]"
}

func (n *Normalize) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", n.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse parameters with defaults
	target := flags["target"]
	if target == "" {
		target = "-16" // EBU R128 standard
	}

	truePeak := flags["true-peak"]
	if truePeak == "" {
		truePeak = "-1.5"
	}

	lra := flags["lra"]
	if lra == "" {
		lra = "11" // Loudness Range
	}

	// Build loudnorm filter
	// I=target:TP=true_peak:LRA=loudness_range
	af := fmt.Sprintf("loudnorm=I=%s:TP=%s:LRA=%s", target, truePeak, lra)

	cmdArgs := []string{
		"-i", input,
		"-af", af,
		"-c:v", "copy", // Keep video unchanged
		"-y",
		output,
	}

	fmt.Printf("Normalizing audio to EBU R128 (target: %s LUFS): %s -> %s\n", target, input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Normalize{})
}
