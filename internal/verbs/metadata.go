package verbs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Metadata reads or modifies video metadata
type Metadata struct{}

func (m *Metadata) Name() string {
	return "metadata"
}

func (m *Metadata) Description() string {
	return "Read or modify video metadata"
}

func (m *Metadata) Usage() string {
	return "mediax metadata <input> [--get title | --set title=\"My Video\" --set artist=\"John\"]"
}

func (m *Metadata) Execute(args []string, flags map[string]string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s", m.Usage())
	}

	input := args[0]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Check if we're getting a specific metadata field
	if getField := flags["get"]; getField != "" {
		cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", fmt.Sprintf("format_tags=%s", getField), "-of", "default=noprint_wrappers=1:nokey=1", input)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get metadata: %v", err)
		}
		value := strings.TrimSpace(string(output))
		if value == "" {
			fmt.Printf("Metadata field '%s' is empty or not set\n", getField)
		} else {
			fmt.Printf("%s: %s\n", getField, value)
		}
		return nil
	}

	// Check if we're setting metadata
	setMetadata := flags["set"]
	if setMetadata != "" {
		// Parse multiple set flags if provided
		// Build ffmpeg command to copy and set metadata
		cmdArgs := []string{"-i", input, "-c", "copy"}

		// Parse set flags - user can pass multiple --set flags
		// For now, handle simple key=value format
		if strings.Contains(setMetadata, "=") {
			parts := strings.SplitN(setMetadata, "=", 2)
			if len(parts) == 2 {
				cmdArgs = append(cmdArgs, "-metadata", fmt.Sprintf("%s=%s", parts[0], parts[1]))
			}
		}

		// Add output
		if len(args) >= 2 {
			cmdArgs = append(cmdArgs, "-y", args[1])
		} else {
			// Overwrite input (not recommended, but supported)
			cmdArgs = append(cmdArgs, "-y", input+"_new"+filepath.Ext(input))
		}

		cmd := exec.Command("ffmpeg", cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Default: show all metadata
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", input)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to read metadata: %v", err)
	}

	// Pretty print JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err == nil {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println(string(output))
	}

	return nil
}

func init() {
	Register(&Metadata{})
}
