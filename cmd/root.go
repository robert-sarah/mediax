// Package cmd provides the CLI commands for Mediax
// Created by Levi Enama

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"mediax/assets"
	"mediax/internal/tui"
	"mediax/internal/verbs"
)

var rootCmd = &cobra.Command{
	Use:   "mediax",
	Short: "Mediax - FFmpeg Simplified with 36 Powerful Verbs",
	Long: `Mediax is a powerful FFmpeg wrapper that transforms complex commands
into simple, intuitive verbs. Perfect for beginners and professionals alike.

Usage:
  mediax <verb> <input> [output] [flags]

Examples:
  mediax probe video.mp4                    # Analyze media file
  mediax convert input.mov output.gif       # Convert format
  mediax trim video.mp4 --start 00:10 --duration 30 output.mp4
  mediax compress large.mp4 small.mp4 --quality medium
  mediax resize video.mp4 720p              # Resize to 720p

Run 'mediax verbs' to see all available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Launch interactive TUI
			p := tui.NewMainMenu()
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				cmd.Help()
			}
			return
		}
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listVerbsCmd)
	rootCmd.AddCommand(docsCmd)

	// Execute verbs dynamically
	for name, verb := range verbs.GetAll() {
		// Capture loop variables
		v := verb
		n := name

		verbCmd := &cobra.Command{
			Use:                n + " [args]",
			Short:              v.Description(),
			Long:               fmt.Sprintf("%s\n\nUsage: %s", v.Description(), v.Usage()),
			DisableFlagParsing: true,
			RunE: func(cmd *cobra.Command, rawArgs []string) error {
				for _, a := range rawArgs {
					if a == "-h" || a == "--help" {
						return cmd.Help()
					}
				}

				positional, flags, err := parseVerbArgs(rawArgs)
				if err != nil {
					return err
				}

				return v.Execute(positional, flags)
			},
		}

		rootCmd.AddCommand(verbCmd)
	}
}

// parseVerbArgs splits raw CLI tokens into positional arguments and flags.
// Supported forms: --flag value | --flag=value | --flag (boolean, "true")
// and their single-dash equivalents (-x value, -x=value, -x).
func parseVerbArgs(rawArgs []string) ([]string, map[string]string, error) {
	positional := []string{}
	flags := make(map[string]string)

	for i := 0; i < len(rawArgs); i++ {
		tok := rawArgs[i]

		if !strings.HasPrefix(tok, "-") || tok == "-" {
			positional = append(positional, tok)
			continue
		}

		name := strings.TrimLeft(tok, "-")
		if name == "" {
			return nil, nil, fmt.Errorf("invalid flag: %s", tok)
		}

		if eq := strings.Index(name, "="); eq != -1 {
			flags[name[:eq]] = name[eq+1:]
			continue
		}

		// Look ahead: if the next token exists and isn't itself a flag,
		// treat it as this flag's value. Otherwise this is a boolean flag.
		if i+1 < len(rawArgs) && !strings.HasPrefix(rawArgs[i+1], "-") {
			flags[name] = rawArgs[i+1]
			i++
		} else {
			flags[name] = "true"
		}
	}

	return positional, flags, nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display the version of Mediax and its dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		assets.PrintLogo()
		fmt.Println()
		fmt.Println("Version:    1.0.0")
		fmt.Println("Author:     Created by Levi Enama")
		fmt.Println("License:    MIT")
		fmt.Println()
		fmt.Println("Dependencies:")
		fmt.Println("  - FFmpeg (https://ffmpeg.org)")
		fmt.Println("  - FFprobe (bundled with FFmpeg)")
	},
}

var listVerbsCmd = &cobra.Command{
	Use:   "verbs",
	Short: "List all available verbs",
	Long:  `Display a complete list of all 36 available verbs in Mediax.`,
	Run: func(cmd *cobra.Command, args []string) {
		assets.PrintLogo()
		fmt.Println()

		categories := []struct {
			name  string
			verbs []string
		}{
			{
				name: "Analysis & Information",
				verbs: []string{"probe"},
			},
			{
				name: "Conversion & Compression",
				verbs: []string{"convert", "compress", "gif"},
			},
			{
				name: "Audio Processing",
				verbs: []string{"extract-audio", "mute", "volume"},
			},
			{
				name: "Video Manipulation",
				verbs: []string{"extract-video", "trim", "crop", "resize", "rotate", "flip"},
			},
			{
				name: "Effects & Filters",
				verbs: []string{"speed", "reverse", "blur", "fade-in", "fade-out"},
			},
			{
				name: "Composition",
				verbs: []string{"concat", "watermark"},
			},
		}

		for _, cat := range categories {
			fmt.Printf("\n  [ %s ]\n", cat.name)
			fmt.Println("  " + strings.Repeat("-", 60))
			for _, name := range cat.verbs {
				if v, ok := verbs.Get(name); ok {
					fmt.Printf("    %-18s  %s\n", name, v.Description())
				}
			}
		}

		fmt.Println()
		fmt.Println("  Use 'mediax <verb> --help' for detailed usage information.")
		fmt.Println()
	},
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Show documentation",
	Long:  `Display helpful documentation and examples for Mediax.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
MEDIAX - QUICK START GUIDE
==========================

Mediax wraps FFmpeg with 36 simple verbs for everyday media tasks.

BASIC USAGE
-----------
  mediax <verb> <input> [output] [flags]

COMMON EXAMPLES
---------------

1. Analyze a video file:
   mediax probe video.mp4

2. Convert between formats:
   mediax convert input.mov output.mp4
   mediax convert video.mp4 audio.mp3

3. Compress a video:
   mediax compress large.mp4 small.mp4 --quality medium

4. Extract audio:
   mediax extract-audio video.mp4 soundtrack.mp3

5. Trim a video:
   mediax trim video.mp4 clip.mp4 --start 00:01:30 --duration 10

6. Resize video:
   mediax resize video.mp4 720p
   mediax resize video.mp4 1920:1080

7. Create GIF:
   mediax gif video.mp4 animation.gif --fps 15 --width 480

8. Add watermark:
   mediax watermark video.mp4 branded.mp4 --text "(c) 2026" --pos bottom-right

9. Change speed:
   mediax speed video.mp4 slowmo.mp4 --rate 0.5
   mediax speed video.mp4 fast.mp4 --rate 2.0

10. Join videos:
    mediax concat final.mp4 part1.mp4 part2.mp4 part3.mp4

FOR MORE INFORMATION
--------------------
  mediax verbs      - List all available commands
  mediax version    - Show version and credits
  mediax <verb> --help

Created by Levi Enama
`)
	},
}
