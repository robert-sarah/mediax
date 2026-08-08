package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Subtitle burns subtitles into video
type Subtitle struct{}

func (s *Subtitle) Name() string {
	return "subtitle"
}

func (s *Subtitle) Description() string {
	return "Burn subtitles into video"
}

func (s *Subtitle) Usage() string {
	return "mediax subtitle <input> <subtitle-file> <output> [--style fontcolor=white:fontsize=24]"
}

func (s *Subtitle) Execute(args []string, flags map[string]string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: %s", s.Usage())
	}

	input, subtitleFile, output := args[0], args[1], args[2]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}
	if _, err := os.Stat(subtitleFile); err != nil {
		return fmt.Errorf("subtitle file not found: %s", subtitleFile)
	}

	// Parse style options
	style := flags["style"]
	if style == "" {
		// Default style
		style = "FontName=Arial:FontSize=24:PrimaryColour=&H00FFFFFF:OutlineColour=&H00000000:Outline=2:Shadow=0"
	}

	// Build subtitle filter. ffmpeg's filtergraph syntax uses ':' to
	// separate filter options, so any ':' inside the force_style value
	// (e.g. "FontName=Arial:FontSize=24") must be escaped, or ffmpeg
	// misinterprets each key=value pair as a separate filter option.
	escapedStyle := strings.ReplaceAll(style, ":", "\\:")
	vf := fmt.Sprintf("subtitles=%s:force_style='%s'", subtitleFile, escapedStyle)

	cmdArgs := []string{
		"-i", input,
		"-vf", vf,
		"-c:a", "copy",
		"-y",
		output,
	}

	fmt.Printf("Burning subtitles into: %s -> %s\n", input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Subtitle{})
}
