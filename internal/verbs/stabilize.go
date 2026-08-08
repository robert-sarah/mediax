package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// Stabilize applies video stabilization to reduce shakiness
type Stabilize struct{}

func (s *Stabilize) Name() string {
	return "stabilize"
}

func (s *Stabilize) Description() string {
	return "Stabilize shaky video footage"
}

func (s *Stabilize) Usage() string {
	return "mediax stabilize <input> <output> [--smooth 10] [--zoom 0]"
}

func (s *Stabilize) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", s.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Parse parameters
	smooth := flags["smooth"]
	if smooth == "" {
		smooth = "10"
	}

	zoom := flags["zoom"]
	if zoom == "" {
		zoom = "0"
	}

	// Two-pass stabilization
	// Pass 1: Analyze motion
	fmt.Println("Pass 1: Analyzing motion...")
	
	analyzeCmd := exec.Command("ffmpeg",
		"-i", input,
		"-vf", "vidstabdetect=shakiness=10:accuracy=15:result=/tmp/transforms.trf",
		"-f", "null", "-",
	)
	analyzeCmd.Stdout = os.Stdout
	analyzeCmd.Stderr = os.Stderr
	
	if err := analyzeCmd.Run(); err != nil {
		return fmt.Errorf("motion analysis failed: %v", err)
	}

	// Pass 2: Apply stabilization
	fmt.Println("Pass 2: Applying stabilization...")
	
	vf := fmt.Sprintf("vidstabtransform=input=/tmp/transforms.trf:smoothing=%s:zoom=%s:optzoom=2", smooth, zoom)
	
	stabilizeCmd := exec.Command("ffmpeg",
		"-i", input,
		"-vf", vf,
		"-c:a", "copy",
		"-y",
		output,
	)
	stabilizeCmd.Stdout = os.Stdout
	stabilizeCmd.Stderr = os.Stderr
	
	// Clean up transform file
	defer os.Remove("/tmp/transforms.trf")
	
	return stabilizeCmd.Run()
}

func init() {
	Register(&Stabilize{})
}
