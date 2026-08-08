package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Resize resizes video to specific dimensions
type Resize struct{}

func (r *Resize) Name() string {
	return "resize"
}

func (r *Resize) Description() string {
	return "Resize video to specific dimensions or preset"
}

func (r *Resize) Usage() string {
	return "mediax resize <input> <output> [1080p | 720p | 480p | 360p | 50% | 1920:1080]"
}

func (r *Resize) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", r.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Determine target size
	size := flags["size"]
	if size == "" && len(args) >= 3 {
		size = args[2]
	}
	if size == "" {
		size = "1080p" // default
	}

	// Parse size to ffmpeg scale filter
	var scaleFilter string
	switch strings.ToLower(size) {
	case "4k", "2160p", "uhd":
		scaleFilter = "3840:2160"
	case "1440p", "2k":
		scaleFilter = "2560:1440"
	case "1080p", "fhd":
		scaleFilter = "1920:1080"
	case "720p", "hd":
		scaleFilter = "1280:720"
	case "480p", "sd":
		scaleFilter = "854:480"
	case "360p":
		scaleFilter = "640:360"
	default:
		// Check if it's a percentage
		if strings.HasSuffix(size, "%") {
			pct := strings.TrimSuffix(size, "%")
			if val, err := strconv.ParseFloat(pct, 64); err == nil {
				scale := val / 100
				scaleFilter = fmt.Sprintf("iw*%.2f:ih*%.2f", scale, scale)
			}
		} else if strings.Contains(size, ":") {
			// Direct scale format: 1920:1080
			scaleFilter = size
		} else if strings.Contains(size, "x") {
			// Format: 1920x1080
			scaleFilter = strings.Replace(size, "x", ":", 1)
		}
	}

	if scaleFilter == "" {
		scaleFilter = "1920:1080"
	}

	cmdArgs := []string{
		"-i", input,
		"-vf", fmt.Sprintf("scale=%s", scaleFilter),
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Resizing %s to %s (%s)\n", input, size, scaleFilter)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Resize{})
}
