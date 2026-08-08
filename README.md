# Mediax

FFmpeg's Cooler Cousin • 36 Powerful Verbs • Interactive Interface

![Mediax Logo](https://img.shields.io/badge/Mediax-1.0.0-purple)
![Go Version](https://img.shields.io/badge/Go-1.21+-blue)
![License](https://img.shields.io/badge/License-MIT-green)

Mediax is a powerful CLI tool that wraps FFmpeg with 36 simple, intuitive verbs for easy media processing. Perfect for beginners and professionals alike.

## Features

- **36 Powerful Verbs**: Convert, compress, trim, crop, rotate, and more
- **Killer Features**: WTF analysis, share links, platform templates, batch processing
- **Interactive TUI**: Beautiful terminal user interface with category navigation
- **Simple Syntax**: Transform complex FFmpeg commands into intuitive verbs
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Professional Design**: Clean, modern interface with smooth animations

## Installation

### Prerequisites

- Go 1.21 or higher
- FFmpeg installed and available in your PATH

### Build from Source

```bash
git clone https://github.com/robert-sarah/mediax.git
cd mediax
go mod tidy
go build -o mediax
```

### Install FFmpeg

**Windows:**
```bash
winget install ffmpeg
```

**macOS:**
```bash
brew install ffmpeg
```

**Linux:**
```bash
sudo apt install ffmpeg  # Ubuntu/Debian
sudo dnf install ffmpeg  # Fedora
```

## Usage

### Interactive Mode

Simply run `mediax` without arguments to launch the interactive TUI:

```bash
mediax
```

Navigate through categories using arrow keys, select a verb with Enter, and follow the prompts.

### Command Line Mode

Use verbs directly from the command line:

```bash
mediax <verb> <input> [output] [flags]
```

## Available Verbs

### Analysis & Information
- `probe` - Analyze media file information
- `info` - Display detailed media information
- `metadata` - Extract and display metadata
- `wtf` - Detect video issues (bitrate, resolution, codec problems)

### Conversion & Compression
- `convert` - Convert between media formats
- `compress` - Compress video/audio files
- `gif` - Create animated GIFs from video

### Audio
- `extract-audio` - Extract audio track from video (full or partial)
- `mute` - Remove audio from video
- `volume` - Adjust audio volume
- `normalize` - Normalize audio levels
- `replace-audio` - Replace audio track in video (full or partial)

### Video
- `extract-video` - Extract video track without audio
- `trim` - Trim video to specific time range
- `crop` - Crop video to specific dimensions
- `resize` - Resize video to specific resolution
- `rotate` - Rotate video by degrees
- `flip` - Flip video horizontally or vertically
- `concat` - Concatenate multiple videos
- `split` - Split video into segments
- `loop` - Loop video for specified duration

### Effects & Filters
- `speed` - Change playback speed
- `reverse` - Reverse video/audio playback
- `blur` - Apply blur effect
- `sharpen` - Sharpen video
- `fade-in` - Add fade-in effect
- `fade-out` - Add fade-out effect
- `stabilize` - Stabilize shaky video
- `denoise` - Remove noise from video
- `chroma` - Chroma key (green screen) effect

### Advanced
- `watermark` - Add watermark to video
- `subtitle` - Add subtitles to video
- `thumbnail` - Extract thumbnail from video
- `slide` - Create slideshow from images
- `batch` - Batch process multiple files
- `template` - Resize for platforms (Instagram, YouTube, TikTok, etc.)
- `share` - Generate share links for processed files

## Examples

### Analyze a video file
```bash
mediax probe video.mp4
```

### Convert between formats
```bash
mediax convert input.mov output.mp4
mediax convert video.mp4 audio.mp3
```

### Compress a video
```bash
mediax compress large.mp4 small.mp4 --quality medium
```

### Extract audio
```bash
mediax extract-audio video.mp4 soundtrack.mp3
```

### Trim a video
```bash
mediax trim video.mp4 clip.mp4 --start 00:01:30 --duration 10
```

### Resize video
```bash
mediax resize video.mp4 720p
mediax resize video.mp4 1920:1080
```

### Create GIF
```bash
mediax gif video.mp4 animation.gif --fps 15 --width 480
```

### Add watermark
```bash
mediax watermark video.mp4 branded.mp4 --text "(c) 2026" --pos bottom-right
```

### Change speed
```bash
mediax speed video.mp4 slowmo.mp4 --rate 0.5
mediax speed video.mp4 fast.mp4 --rate 2.0
```

### Join videos
```bash
mediax concat final.mp4 part1.mp4 part2.mp4 part3.mp4
```

### Detect video issues (WTF)
```bash
mediax wtf video.mp4
```

### Resize for Instagram
```bash
mediax template video.mp4 insta_video.mp4 --platform instagram --aspect portrait
```

### Resize for YouTube
```bash
mediax template video.mp4 yt_video.mp4 --platform youtube
```

### Generate share link
```bash
mediax share video.mp4 --platform google-drive
```

### Batch process folder
```bash
mediax batch convert "*.mp4" "output/{name}_converted.mp4" --args "--quality high"
```

## Official Verb Specifications

### replace-audio
Replace audio track in video (full or partial).

**Syntax:**
```bash
mediax replace-audio <video> <new_audio> <output> [flags]
```

**Flags:**
- `--start HH:MM:SS` - Start time for replacement
- `--duration HH:MM:SS` - Duration of replacement
- `--time HH:MM:SS-HH:MM:SS` - Time range (alternative to start+duration)

**Examples:**
```bash
# Replace entire audio
mediax replace-audio video.mp4 new.mp3 output.mp4

# Replace from 1:30 to end
mediax replace-audio video.mp4 new.mp3 out.mp4 --start 00:01:30

# Replace from 1:30 to 1:40 (10 seconds)
mediax replace-audio video.mp4 new.mp3 out.mp4 --start 00:01:30 --duration 10

# Same with time range
mediax replace-audio video.mp4 new.mp3 out.mp4 --time 00:01:30-00:01:40
```

### extract-audio
Extract audio track from video (full or partial).

**Syntax:**
```bash
mediax extract-audio <video> <output> [flags]
```

**Flags:**
- `--start HH:MM:SS` - Start time for extraction
- `--duration HH:MM:SS` - Duration to extract
- `--time HH:MM:SS-HH:MM:SS` - Time range to extract
- `--format mp3|aac|wav|flac` - Output audio format

**Examples:**
```bash
# Extract entire audio
mediax extract-audio video.mp4 soundtrack.mp3

# Extract from 1:30 to end
mediax extract-audio video.mp4 clip.mp3 --start 00:01:30

# Extract 10 seconds from 1:30
mediax extract-audio video.mp4 clip.mp3 --start 00:01:30 --duration 10

# Same with time range
mediax extract-audio video.mp4 clip.mp3 --time 00:01:30-00:01:40
```

## Commands

- `mediax` - Launch interactive TUI
- `mediax verbs` - List all available verbs
- `mediax version` - Display version information
- `mediax docs` - Show documentation
- `mediax <verb> --help` - Show help for specific verb

## Development

### Project Structure

```
mediax/
├── cmd/           # CLI command definitions
├── internal/
│   ├── tui/       # Terminal user interface
│   └── verbs/     # Verb implementations
├── assets/        # Logo and branding
└── main.go        # Entry point
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Created by Levi Enama

## Acknowledgments

- FFmpeg team for the amazing media processing framework
- Charmbracelet for the beautiful TUI libraries (Bubbletea, Lipgloss, Bubbles)
- Cobra for the CLI framework

## Support

For issues, questions, or suggestions, please open an issue on GitHub.
