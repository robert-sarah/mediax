package verbs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// WTF analyzes video and detects potential issues
type WTF struct{}

func (w *WTF) Name() string {
	return "wtf"
}

func (w *WTF) Description() string {
	return "Analyze video and detect issues (bitrate, resolution, codec, etc.)"
}

func (w *WTF) Usage() string {
	return "mediax wtf <input>"
}

type Issue struct {
	Type     string
	Severity string
	Message  string
}

func (w *WTF) Execute(args []string, flags map[string]string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s", w.Usage())
	}

	input := args[0]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("file not found: %s", input)
	}

	fmt.Printf("Analyzing: %s\n\n", input)

	// Get video info using ffprobe
	cmdArgs := []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", input}
	cmd := exec.Command("ffprobe", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to analyze video: %v", err)
	}

	type wtfStream struct {
		CodecType string `json:"codec_type"`
		CodecName string `json:"codec_name"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		BitRate   string `json:"bit_rate"`
		Framerate string `json:"r_frame_rate"`
	}

	var probeData struct {
		Format struct {
			BitRate    string `json:"bit_rate"`
			Duration   string `json:"duration"`
			FormatName string `json:"format_name"`
		} `json:"format"`
		Streams []wtfStream `json:"streams"`
	}

	if err := json.Unmarshal(output, &probeData); err != nil {
		return fmt.Errorf("failed to parse video info: %v", err)
	}

	issues := []Issue{}

	// Analyze video stream
	var videoStream *wtfStream

	for i := range probeData.Streams {
		if probeData.Streams[i].CodecType == "video" {
			videoStream = &probeData.Streams[i]
			break
		}
	}

	if videoStream != nil {
		// Check resolution
		if videoStream.Width < 720 {
			issues = append(issues, Issue{
				Type:     "Resolution",
				Severity: "WARNING",
				Message:  fmt.Sprintf("Low resolution: %dx%d (recommended: 720p minimum)", videoStream.Width, videoStream.Height),
			})
		}

		// Check for non-standard resolutions
		standardResolutions := []struct{ w, h int }{
			{1920, 1080}, {1280, 720}, {3840, 2160}, {2560, 1440},
		}
		isStandard := false
		for _, res := range standardResolutions {
			if videoStream.Width == res.w && videoStream.Height == res.h {
				isStandard = true
				break
			}
		}
		if !isStandard && videoStream.Width > 0 {
			issues = append(issues, Issue{
				Type:     "Resolution",
				Severity: "INFO",
				Message:  fmt.Sprintf("Non-standard resolution: %dx%d", videoStream.Width, videoStream.Height),
			})
		}

		// Check codec
		if videoStream.CodecName != "h264" && videoStream.CodecName != "hevc" && videoStream.CodecName != "av1" {
			issues = append(issues, Issue{
				Type:     "Codec",
				Severity: "WARNING",
				Message:  fmt.Sprintf("Non-standard codec: %s (recommended: h264/hevc/av1)", videoStream.CodecName),
			})
		}

		// Check bitrate
		if videoStream.BitRate != "" {
			bitrate := parseInt(videoStream.BitRate)
			if bitrate > 0 && bitrate < 1000000 {
				issues = append(issues, Issue{
					Type:     "Bitrate",
					Severity: "WARNING",
					Message:  fmt.Sprintf("Low bitrate: %d kbps (recommended: 1000+ kbps for HD)", bitrate/1000),
				})
			}
		}
	}

	// Check overall format bitrate
	if probeData.Format.BitRate != "" {
		bitrate := parseInt(probeData.Format.BitRate)
		if bitrate > 0 && bitrate < 500000 {
			issues = append(issues, Issue{
				Type:     "Bitrate",
				Severity: "WARNING",
				Message:  fmt.Sprintf("Very low overall bitrate: %d kbps", bitrate/1000),
			})
		}
	}

	// Check duration
	if probeData.Format.Duration != "" {
		duration := parseFloat(probeData.Format.Duration)
		if duration > 0 && duration < 1 {
			issues = append(issues, Issue{
				Type:     "Duration",
				Severity: "INFO",
				Message:  fmt.Sprintf("Very short video: %.2f seconds", duration),
			})
		}
	}

	// Display results
	fmt.Println("=== VIDEO ANALYSIS ===")
	fmt.Printf("Codec: %s\n", func() string {
		if videoStream != nil {
			return videoStream.CodecName
		}
		return "N/A"
	}())
	fmt.Printf("Resolution: %dx%d\n", func() int {
		if videoStream != nil {
			return videoStream.Width
		}
		return 0
	}(), func() int {
		if videoStream != nil {
			return videoStream.Height
		}
		return 0
	}())
	fmt.Printf("Duration: %s\n", probeData.Format.Duration)
	fmt.Printf("Format: %s\n", probeData.Format.FormatName)
	fmt.Println()

	if len(issues) == 0 {
		fmt.Println("✓ No issues detected! Your video looks good.")
		return nil
	}

	fmt.Println("=== DETECTED ISSUES ===")
	for _, issue := range issues {
		severityColor := ""
		switch issue.Severity {
		case "WARNING":
			severityColor = "[WARNING]"
		case "INFO":
			severityColor = "[INFO]"
		default:
			severityColor = "[" + issue.Severity + "]"
		}
		fmt.Printf("%s %s: %s\n", severityColor, issue.Type, issue.Message)
	}

	return nil
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func parseFloat(s string) float64 {
	var result float64
	fmt.Sscanf(s, "%f", &result)
	return result
}

func init() {
	Register(&WTF{})
}
