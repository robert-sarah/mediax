package assets

import (
	"fmt"
	"strings"
)

const LogoColor = `
    ███╗   ███╗███████╗██████╗ ██╗ █████╗ ██╗  ██╗
    ████╗ ████║██╔════╝██╔══██╗██║██╔══██╗╚██╗██╔╝
    ██╔████╔██║█████╗  ██║  ██║██║███████║ ╚███╔╝ 
    ██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║ ██╔██╗ 
    ██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║██╔╝ ██╗
    ╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝
`

const LogoPlain = `
    ███╗   ███╗███████╗██████╗ ██╗ █████╗ ██╗  ██╗
    ████╗ ████║██╔════╝██╔══██╗██║██╔══██╗╚██╗██╔╝
    ██╔████╔██║█████╗  ██║  ██║██║███████║ ╚███╔╝ 
    ██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║ ██╔██╗ 
    ██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║██╔╝ ██╗
    ╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝
`

const Tagline = "FFmpeg's Cooler Cousin • 36 Powerful Verbs • Interactive Interface"

const BoxTop = "┌─────────────────────────────────────────────────────────────────────────────┐"
const BoxBottom = "└─────────────────────────────────────────────────────────────────────────────┘"
const BoxSide = "│"

func PrintLogo() {
	fmt.Println(LogoPlain)
	fmt.Println()
	fmt.Printf("    %s\n", Tagline)
	fmt.Println()
}

func PrintBoxed(text string) {
	fmt.Println(BoxTop)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if len(line) > 73 {
			line = line[:73]
		}
		fmt.Printf("│ %-73s │\n", line)
	}
	fmt.Println(BoxBottom)
}

func PrintHelpHeader(title string) {
	fmt.Println()
	fmt.Printf("  ╔══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("  ║  %-66s  ║\n", title)
	fmt.Printf("  ╚══════════════════════════════════════════════════════════════════════╝\n")
	fmt.Println()
}
