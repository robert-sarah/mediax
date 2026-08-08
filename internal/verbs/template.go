package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Template resizes videos for specific platforms
type Template struct{}

func (t *Template) Name() string {
	return "template"
}

func (t *Template) Description() string {
	return "Resize video for specific platforms (Instagram, YouTube, Twitter, TikTok, etc.)"
}

func (t *Template) Usage() string {
	return "mediax template <input> <output> --platform instagram|youtube|twitter|tiktok [--aspect square|portrait|landscape]"
}

func (t *Template) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", t.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	platform := flags["platform"]
	if platform == "" {
		return fmt.Errorf("must specify --platform (instagram, youtube, twitter, tiktok)")
	}

	var width, height int
	var bitrate string
	var fps int

	switch strings.ToLower(platform) {
	case "instagram", "ig":
		aspect := flags["aspect"]
		if aspect == "" {
			aspect = "portrait"
		}
		switch strings.ToLower(aspect) {
		case "square":
			width, height = 1080, 1080
		case "portrait":
			width, height = 1080, 1350
		case "landscape":
			width, height = 1080, 1920
		default:
			width, height = 1080, 1080
		}
		bitrate = "3500k"
		fps = 30

	case "youtube", "yt":
		width, height = 1920, 1080
		bitrate = "8000k"
		fps = 30

	case "twitter", "x":
		width, height = 1280, 720
		bitrate = "5000k"
		fps = 30

	case "tiktok":
		width, height = 1080, 1920
		bitrate = "4000k"
		fps = 30

	case "facebook", "fb":
		width, height = 1280, 720
		bitrate = "4000k"
		fps = 30

	case "linkedin":
		width, height = 1920, 1080
		bitrate = "5000k"
		fps = 30

	default:
		return fmt.Errorf("unsupported platform: %s (supported: instagram, youtube, twitter, tiktok, facebook, linkedin)", platform)
	}

	fmt.Printf("Resizing for %s: %dx%d @ %dfps\n", platform, width, height, fps)

	cmdArgs := []string{
		"-i", input,
		"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2", width, height, width, height),
		"-c:v", "libx264",
		"-preset", "medium",
		"-b:v", bitrate,
		"-maxrate", bitrate,
		"-bufsize", fmt.Sprintf("%s", strings.TrimSuffix(bitrate, "k")+"k"),
		"-r", fmt.Sprintf("%d", fps),
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		"-y",
		output,
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Template{})
}
