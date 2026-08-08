package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Watermark adds text or image overlay to video
type Watermark struct{}

func (w *Watermark) Name() string {
	return "watermark"
}

func (w *Watermark) Description() string {
	return "Add text or image watermark to video"
}

func (w *Watermark) Usage() string {
	return "mediax watermark <input> <output> [--text \"Copyright 2026\" | --image logo.png] [--pos bottom-right] [--opacity 0.5]"
}

func (w *Watermark) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", w.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	text := flags["text"]
	image := flags["image"]

	if text == "" && image == "" {
		return fmt.Errorf("must specify --text or --image for watermark")
	}

	// Parse position
	pos := flags["pos"]
	if pos == "" {
		pos = flags["position"]
	}
	if pos == "" {
		pos = "bottom-right"
	}

	opacity, _ := strconv.ParseFloat(flags["opacity"], 64)
	if opacity <= 0 || opacity > 1 {
		opacity = 0.7
	}

	var filterComplex string

	if text != "" {
		// Text watermark using drawtext
		fontFile := "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
		if _, err := os.Stat(fontFile); err != nil {
			fontFile = "/System/Library/Fonts/Helvetica.ttc"
			if _, err := os.Stat(fontFile); err != nil {
				fontFile = "c:/Windows/Fonts/arial.ttf"
			}
		}

		// Calculate position
		xPos, yPos := "w-text_w-10", "h-text_h-10"
		switch pos {
		case "top-left":
			xPos, yPos = "10", "10"
		case "top-center", "top":
			xPos, yPos = "(w-text_w)/2", "10"
		case "top-right":
			xPos, yPos = "w-text_w-10", "10"
		case "center-left", "left":
			xPos, yPos = "10", "(h-text_h)/2"
		case "center", "middle":
			xPos, yPos = "(w-text_w)/2", "(h-text_h)/2"
		case "center-right", "right":
			xPos, yPos = "w-text_w-10", "(h-text_h)/2"
		case "bottom-left":
			xPos, yPos = "10", "h-text_h-10"
		case "bottom-center", "bottom":
			xPos, yPos = "(w-text_w)/2", "h-text_h-10"
		case "bottom-right":
			xPos, yPos = "w-text_w-10", "h-text_h-10"
		}

		filterComplex = fmt.Sprintf("drawtext=fontfile=%s:text='%s':x=%s:y=%s:fontsize=24:fontcolor=white@%0.2f:box=1:boxcolor=black@0.5",
			fontFile, text, xPos, yPos, opacity)
	} else {
		// Image watermark using overlay
		if _, err := os.Stat(image); err != nil {
			return fmt.Errorf("watermark image not found: %s", image)
		}

		// Calculate overlay position
		overlayPos := "overlay=W-w-10:H-h-10"
		switch pos {
		case "top-left":
			overlayPos = "overlay=10:10"
		case "top-center", "top":
			overlayPos = "overlay=(W-w)/2:10"
		case "top-right":
			overlayPos = "overlay=W-w-10:10"
		case "center-left", "left":
			overlayPos = "overlay=10:(H-h)/2"
		case "center", "middle":
			overlayPos = "overlay=(W-w)/2:(H-h)/2"
		case "center-right", "right":
			overlayPos = "overlay=W-w-10:(H-h)/2"
		case "bottom-left":
			overlayPos = "overlay=10:H-h-10"
		case "bottom-center", "bottom":
			overlayPos = "overlay=(W-w)/2:H-h-10"
		case "bottom-right":
			overlayPos = "overlay=W-w-10:H-h-10"
		}

		// Format opacity for overlay
		opacityStr := fmt.Sprintf("%.2f", opacity)

		// Use colorchannelmixer for opacity
		filterComplex = fmt.Sprintf("[1:v]format=rgba,colorchannelmixer=aa=%s[wm];[0:v][wm]%s",
			opacityStr, overlayPos)
	}

	var cmdArgs []string

	if text != "" {
		cmdArgs = []string{"-i", input, "-vf", filterComplex, "-c:a", "copy", "-y", output}
	} else {
		cmdArgs = []string{"-i", input, "-i", image, "-filter_complex", filterComplex, "-c:a", "copy", "-y", output}
	}

	fmt.Printf("Adding watermark to: %s -> %s\n", input, output)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Watermark{})
}
