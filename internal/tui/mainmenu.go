package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"mediax/assets"
	"mediax/internal/verbs"
)

var (
	// Professional color palette
	primaryColor   = lipgloss.Color("#7C3AED")   // Violet
	secondaryColor = lipgloss.Color("#06B6D4")  // Cyan
	accentColor    = lipgloss.Color("#F59E0B")  // Amber
	successColor   = lipgloss.Color("#10B981")  // Emerald
	warningColor   = lipgloss.Color("#F97316")  // Orange
	errorColor     = lipgloss.Color("#EF4444")  // Red
	mutedColor     = lipgloss.Color("#6B7280")  // Gray
	darkBg         = lipgloss.Color("#1F2937")  // Dark gray
	lightBg        = lipgloss.Color("#374151")  // Lighter gray

	// Styles
	titleStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginTop(1).
		MarginBottom(1).
		Width(80)

	subtitleStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Italic(true).
		MarginBottom(2)

	itemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("#E5E7EB"))

	selectedItemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(primaryColor).
		Background(lipgloss.Color("#374151")).
		Bold(true)

	descriptionStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginLeft(4).
		MarginTop(0).
		MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		BorderBackground(darkBg).
		Padding(1, 3).
		Margin(1, 0).
		Background(darkBg)

	footerStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		MarginTop(1)

	verbTitleStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	verbDescStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginLeft(2).
		MarginBottom(1)

	spinnerStyle = lipgloss.NewStyle().
		Foreground(successColor)
)

type categoryItem struct {
	title       string
	description string
	verbs       []string
	icon        string
}

func (i categoryItem) FilterValue() string { return i.title }
func (i categoryItem) Title() string        { return i.icon + " " + i.title }
func (i categoryItem) Description() string   { return i.description }

type verbItem struct {
	name        string
	description string
	usage       string
}

func (i verbItem) FilterValue() string { return i.name }
func (i verbItem) Title() string        { return i.name }
func (i verbItem) Description() string   { return i.description }

type model struct {
	state       state
	list        list.Model
	verbList    list.Model
	choice      string
	quitting    bool
	selectedCat categoryItem
	spinner     spinner.Model
	loading     bool
	lastAction  string
}

type state int

const (
	stateMenu state = iota
	stateVerbs
	stateExecuting
	stateDone
)

func NewMainMenu() *tea.Program {
	items := []list.Item{
		categoryItem{
			title:       "Analysis & Information",
			description: "Analyze media files, extract metadata, detect issues",
			verbs:       []string{"probe", "info", "metadata", "wtf"},
			icon:        "[+]",
		},
		categoryItem{
			title:       "Conversion & Compression",
			description: "Convert formats, compress, adjust quality",
			verbs:       []string{"convert", "compress", "gif"},
			icon:        "[>]",
		},
		categoryItem{
			title:       "Audio",
			description: "Extract, modify, enhance audio",
			verbs:       []string{"extract-audio", "mute", "volume", "normalize", "replace-audio"},
			icon:        "[A]",
		},
		categoryItem{
			title:       "Video",
			description: "Crop, resize, trim, assemble videos",
			verbs:       []string{"extract-video", "trim", "crop", "resize", "rotate", "flip", "concat", "split", "loop"},
			icon:        "[V]",
		},
		categoryItem{
			title:       "Effects & Filters",
			description: "Apply visual effects and transitions",
			verbs:       []string{"speed", "reverse", "blur", "sharpen", "fade-in", "fade-out", "stabilize", "denoise", "chroma"},
			icon:        "[*]",
		},
		categoryItem{
			title:       "Advanced",
			description: "Watermark, subtitles, templates, batch processing, sharing",
			verbs:       []string{"watermark", "subtitle", "thumbnail", "slide", "batch", "template", "share"},
			icon:        "[@]",
		},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = descriptionStyle.Foreground(lipgloss.Color("#9CA3AF"))
	delegate.Styles.NormalTitle = itemStyle
	delegate.Styles.NormalDesc = descriptionStyle

	l := list.New(items, delegate, 80, 15)
	l.Title = "  MEDIAX - Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.Styles.Title = titleStyle

	sp := spinner.New()
	sp.Spinner = spinner.Points

	m := model{
		state:   stateMenu,
		list:    l,
		spinner: sp,
	}
	return tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width - 4)
		if m.verbList.Width() > 0 {
			m.verbList.SetWidth(msg.Width - 4)
		}
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			if m.state == stateVerbs {
				m.state = stateMenu
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.state == stateMenu {
				i, ok := m.list.SelectedItem().(categoryItem)
				if ok {
					m.selectedCat = i
					m.state = stateVerbs
					return m, m.initVerbList(i.verbs)
				}
			} else if m.state == stateVerbs {
				i, ok := m.verbList.SelectedItem().(verbItem)
				if ok {
					m.choice = i.name
					m.state = stateExecuting
					m.loading = true
					m.lastAction = fmt.Sprintf("Executing: %s", i.name)
					return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
						return executeVerbMsg{i.name}
					})
				}
			}
		}

	case executeVerbMsg:
		m.loading = false
		m.state = stateDone
		m.lastAction = fmt.Sprintf("Command ready: %s\nUsage: mediax %s", msg.verb, msg.verb)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	if m.state == stateMenu {
		m.list, cmd = m.list.Update(msg)
	} else if m.state == stateVerbs {
		m.verbList, cmd = m.verbList.Update(msg)
	}
	return m, cmd
}

type executeVerbMsg struct {
	verb string
}

func (m model) initVerbList(verbNames []string) tea.Cmd {
	items := make([]list.Item, 0, len(verbNames))
	for _, name := range verbNames {
		if v, ok := verbs.Get(name); ok {
			items = append(items, verbItem{
				name:        name,
				description: v.Description(),
				usage:       v.Usage(),
			})
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = descriptionStyle.Foreground(lipgloss.Color("#9CA3AF"))
	delegate.Styles.NormalTitle = itemStyle
	delegate.Styles.NormalDesc = descriptionStyle

	m.verbList = list.New(items, delegate, 80, 12)
	m.verbList.Title = "  " + m.selectedCat.icon + " " + m.selectedCat.title
	m.verbList.SetShowStatusBar(false)
	m.verbList.SetFilteringEnabled(false)
	m.verbList.SetShowPagination(false)
	m.verbList.Styles.Title = verbTitleStyle

	return nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Header with logo
	s.WriteString("\n")
	s.WriteString(assets.LogoPlain)
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render("    " + assets.Tagline))
	s.WriteString("\n")

	switch m.state {
	case stateMenu:
		s.WriteString(m.list.View())
		s.WriteString("\n")
		s.WriteString(footerStyle.Render("↑/↓: Navigate  •  Enter: Select  •  q/esc: Quit"))

	case stateVerbs:
		s.WriteString(m.verbList.View())
		s.WriteString("\n")
		s.WriteString(footerStyle.Render("↑/↓: Navigate  •  Enter: Execute  •  esc: Back  •  q: Quit"))

	case stateExecuting:
		s.WriteString("\n")
		s.WriteString(boxStyle.Render("\n" + spinnerStyle.Render(m.spinner.View()) + " " + m.lastAction + "\n"))

	case stateDone:
		s.WriteString("\n")
		s.WriteString(boxStyle.Render("\n" + successStyle.Render("[OK] "+m.lastAction) + "\n\n" + mutedStyle.Render("Press q to quit or esc to return") + "\n"))
	}

	return s.String()
}

var mutedStyle = lipgloss.NewStyle().Foreground(mutedColor).Italic(true)
var successStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
