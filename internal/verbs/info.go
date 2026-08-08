package verbs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)
// Info displays comprehensive media file information
type Info struct{}

func (i *Info) Name() string {
	return "info"
}

func (i *Info) Description() string {
	return "Display comprehensive media file information"
}

func (i *Info) Usage() string {
	return "mediax info <input> [--format json | --format text]"
}

func (i *Info) Execute(args []string, flags map[string]string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s", i.Usage())
	}

	input := args[0]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	format := flags["format"]
	if format == "" {
		format = "text"
	}

	// Use ffprobe to get information
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		input,
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to analyze file: %v", err)
	}

	if format == "json" {
		// Pretty print JSON
		var result map[string]interface{}
		if err := json.Unmarshal(output, &result); err == nil {
			prettyJSON, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(prettyJSON))
		}
	} else {
		// Text format - parse and display nicely
		var result map[string]interface{}
		if err := json.Unmarshal(output, &result); err != nil {
			fmt.Println(string(output))
			return nil
		}

		fmt.Printf("\nFile: %s\n", input)
		fmt.Println(strings.Repeat("=", 50))

		if format_info, ok := result["format"].(map[string]interface{}); ok {
			if filename, ok := format_info["filename"].(string); ok {
				fmt.Printf("Filename: %s\n", filepath.Base(filename))
			}
			if duration, ok := format_info["duration"].(string); ok {
				fmt.Printf("Duration: %ss\n", duration)
			}
			if size, ok := format_info["size"].(string); ok {
				fmt.Printf("Size: %s bytes\n", size)
			}
			if bitrate, ok := format_info["bit_rate"].(string); ok {
				fmt.Printf("Bitrate: %s\n", bitrate)
			}
			if format_name, ok := format_info["format_name"].(string); ok {
				fmt.Printf("Format: %s\n", format_name)
			}
		}

		fmt.Println()
		fmt.Println("Streams:")
		fmt.Println(strings.Repeat("-", 50))

		if streams, ok := result["streams"].([]interface{}); ok {
			for i, s := range streams {
				if stream, ok := s.(map[string]interface{}); ok {
					codec_type, _ := stream["codec_type"].(string)
					codec_name, _ := stream["codec_name"].(string)
					fmt.Printf("  [%d] %s - %s\n", i, codec_type, codec_name)

					if codec_type == "video" {
						if width, ok := stream["width"].(float64); ok {
							if height, ok := stream["height"].(float64); ok {
								fmt.Printf("       Resolution: %dx%d\n", int(width), int(height))
							}
						}
						if fps, ok := stream["r_frame_rate"].(string); ok {
							fmt.Printf("       FPS: %s\n", fps)
						}
					} else if codec_type == "audio" {
						if sample_rate, ok := stream["sample_rate"].(string); ok {
							fmt.Printf("       Sample Rate: %s Hz\n", sample_rate)
						}
						if channels, ok := stream["channels"].(float64); ok {
							fmt.Printf("       Channels: %d\n", int(channels))
						}
					}
				}
			}
		}
		fmt.Println()
	}

	return nil
}

func init() {
	Register(&Info{})
}
