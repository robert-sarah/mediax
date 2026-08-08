package verbs

import (
	"fmt"
	"os"
	"os/exec"
)

// Probe analyzes media files using ffprobe
type Probe struct{}

func (p *Probe) Name() string {
	return "probe"
}

func (p *Probe) Description() string {
	return "Analyze media file with ffprobe"
}

func (p *Probe) Usage() string {
	return "mediax probe <input> [--format json]"
}

func (p *Probe) Execute(args []string, flags map[string]string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: mediax probe <input>")
	}

	input := args[0]
	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("file not found: %s", input)
	}

	format := flags["format"]
	if format == "" {
		format = "default"
	}

	var cmdArgs []string
	switch format {
	case "json":
		cmdArgs = []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", input}
	case "csv":
		cmdArgs = []string{"-v", "quiet", "-print_format", "csv", "-show_format", input}
	default:
		cmdArgs = []string{"-v", "error", "-show_format", "-show_streams", input}
	}

	cmd := exec.Command("ffprobe", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Probe{})
}
