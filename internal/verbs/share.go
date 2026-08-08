package verbs

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Share generates a share link for a file
type Share struct{}

func (s *Share) Name() string {
	return "share"
}

func (s *Share) Description() string {
	return "Generate a share link for the processed file"
}

func (s *Share) Usage() string {
	return "mediax share <input> [--platform google-drive|dropbox|onedrive]"
}

func (s *Share) Execute(args []string, flags map[string]string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s", s.Usage())
	}

	input := args[0]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("file not found: %s", input)
	}

	platform := flags["platform"]
	if platform == "" {
		platform = "google-drive"
	}

	absPath, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	var shareLink string

	switch platform {
	case "google-drive", "gdrive":
		shareLink = fmt.Sprintf("https://drive.google.com/drive/u/0/my-drive?path=%s", url.QueryEscape(filepath.Dir(absPath)))
	case "dropbox":
		shareLink = "https://www.dropbox.com/home"
	case "onedrive":
		shareLink = "https://onedrive.live.com/"
	default:
		return fmt.Errorf("unsupported platform: %s (supported: google-drive, dropbox, onedrive)", platform)
	}

	fmt.Printf("File: %s\n", absPath)
	fmt.Printf("Size: %.2f MB\n", getFileSizeMB(absPath))
	fmt.Printf("Platform: %s\n", platform)
	fmt.Printf("\nShare Link: %s\n", shareLink)
	fmt.Println("\nOpening in browser...")

	// Open in browser
	openBrowser(shareLink)

	return nil
}

func getFileSizeMB(path string) float64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return float64(info.Size()) / (1024 * 1024)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Run()
}

func init() {
	Register(&Share{})
}
