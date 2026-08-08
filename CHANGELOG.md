# Changelog

All notable changes to Mediax will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed (bug-fix pass, real testing against ffmpeg)
- **CLI flags were entirely broken**: `cmd/root.go` never registered any flags
  on the dynamically-generated verb subcommands, so every `--flag value` call
  (trim, resize, gif, watermark, speed, crop, etc.) was rejected by Cobra
  before reaching the verb's code. Rewrote flag handling with a custom parser
  (`--flag value`, `--flag=value`, boolean `--flag`) and `DisableFlagParsing`.
- `internal/verbs/fadeout.go`: invalid Go escape sequence (`\,`) prevented
  compilation. Replaced the ad-hoc `expr=max(0,t-x)` trick with a real
  ffprobe-based duration lookup so the fade starts at the correct offset.
- `internal/verbs/replaceaudio.go` / `extractaudio.go`: duplicate
  `calculateDuration`/`parseTimeToSeconds`/`formatSecondsToTime` declarations
  across both files (compile error). Removed the duplicates.
- `internal/verbs/wtf.go`: anonymous struct type mismatch prevented
  compilation; introduced a shared named `wtfStream` type.
- `internal/verbs/flip.go`: `Printf` format string had 4 verbs but 3
  arguments, producing `%!s(MISSING)` in the CLI output.
- `internal/verbs/gif.go`: malformed ffmpeg filtergraph
  (`palettegen=[s1]paletteuse=...`) made every GIF export fail; fixed to
  `palettegen[p];[s1][p]paletteuse=...`.
- `internal/verbs/split.go`: `--parts N` passed ffmpeg's segment muxer a
  nonexistent `-segment_count` option. Now computes the real file duration
  via ffprobe and derives a correct `-segment_time`.
- `internal/verbs/subtitle.go`: unescaped colons in `force_style` broke
  ffmpeg's filter option parser (`Error applying option 'FontSize'`).
  Colons are now escaped before being embedded in the filtergraph.
- `internal/verbs/thumbnail.go`: added `-update 1` to silence/avoid an
  ffmpeg image2 warning when writing a single-frame output.
- `internal/tui/mainmenu.go`: a local variable named `spinner` shadowed the
  `spinner` package, breaking `spinner.Points` (compile error).
- Removed unused imports across `batch.go`, `extractaudio.go`, `wtf.go`,
  `cmd/root.go`, `fadeout.go` that also prevented compilation.
- New shared helper `internal/verbs/ffutil.go` (`getMediaDurationSeconds`)
  used by `fade-out` and `split` for accurate ffprobe-based durations.

All 36 verbs were audited: every one calls real `ffmpeg`/`ffprobe` via
`exec.Command` — no stubs, mocks, or simulated output were found or
introduced. The above are genuine logic/compile bugs found by building the
project and running each verb against a real test video.

### Added
- Interactive Terminal User Interface (TUI) with category navigation
- Professional color scheme and styling
- Verb selection submenu with execution feedback
- Spinner animations for loading states
- 36 powerful media processing verbs
- Cross-platform support (Windows, macOS, Linux)
- Comprehensive documentation


### Changed
- Improved TUI design with professional aesthetics
- Removed emojis for cleaner interface
- Converted all text to English
- Enhanced color palette with modern design

### Fixed
- Import error in mainmenu.go (changed internal/commands to internal/verbs)
- TUI state management for proper navigation

## [1.0.0] - 2026-08-04

### Added
- Initial release of Mediax
- 36 FFmpeg wrapper verbs
- CLI interface with Cobra
- Basic TUI implementation
- Logo and branding assets
