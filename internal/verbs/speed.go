package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Speed changes playback speed (slow motion or fast forward)
type Speed struct{}

func (s *Speed) Name() string {
	return "speed"
}

func (s *Speed) Description() string {
	return "Change playback speed (0.5x slow-mo, 2x fast-forward)"
}

func (s *Speed) Usage() string {
	return "mediax speed <input> <output> [--rate 2.0] [--audio-pitch]"
}

func (s *Speed) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", s.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Get speed rate
	rate := 1.0
	if r := flags["rate"]; r != "" {
		if val, err := strconv.ParseFloat(r, 64); err == nil {
			rate = val
		}
	}

	if rate <= 0 {
		return fmt.Errorf("speed rate must be positive")
	}

	// Calculate PTS (presentation timestamp) factor
	// Speed > 1 = faster, Speed < 1 = slower
	ptsFactor := 1 / rate

	// Build video filter
	vf := fmt.Sprintf("setpts=%.4f*PTS", ptsFactor)

	// Build audio filter (if keeping audio)
	var af string
	if rate == 0.5 {
		af = "atempo=0.5"
	} else if rate == 2.0 {
		af = "atempo=2.0"
	} else {
		// For other rates, use atempo chain or asetrate
		af = fmt.Sprintf("atempo=%.2f", rate)
	}

	cmdArgs := []string{"-i", input}

	// Add video filter
	cmdArgs = append(cmdArgs, "-vf", vf)

	// Handle audio based on pitch flag
	if flags["audio-pitch"] != "" {
		// Change pitch to match speed (chipmunk effect for fast, slow-mo voice for slow)
		cmdArgs = append(cmdArgs, "-af", af)
	} else {
		// Keep original pitch using asetrate (more complex)
		cmdArgs = append(cmdArgs, "-af", fmt.Sprintf("asetrate=48000*%f,aresample=48000", rate))
	}

	cmdArgs = append(cmdArgs, "-y", output)

	var speedDesc string
	if rate > 1 {
		speedDesc = fmt.Sprintf("%gx fast-forward", rate)
	} else if rate < 1 {
		speedDesc = fmt.Sprintf("%gx slow-motion", rate)
	} else {
		speedDesc = "normal speed"
	}

	fmt.Printf("Changing speed: %s -> %s (%s)\n", input, output, speedDesc)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Speed{})
}
