#!/bin/bash

# mediax installer script
# Usage: curl -sSL https://raw.githubusercontent.com/robert-sarah/mediax/main/install.sh | sh
# or: wget -qO- https://raw.githubusercontent.com/robert-sarah/mediax/main/install.sh | sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="robert-sarah/mediax"
BINARY_NAME="mediax"
INSTALL_DIR="/usr/local/bin"

print_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

print_success() {
    echo -e "${GREEN}[OK] $1${NC}"
}

print_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

print_step() {
    echo -e "${YELLOW}[STEP] $1${NC}"
}

# Detect OS and architecture
detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        darwin) echo "darwin" ;;
        linux) echo "linux" ;;
        mingw*|msys*|cygwin*) echo "windows" ;;
        *) print_error "Unsupported OS: $os"; exit 1 ;;
    esac
}

detect_arch() {
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) print_error "Unsupported architecture: $arch"; exit 1 ;;
    esac
}

# Auto-install dependencies
install_ffmpeg() {
    print_step "FFmpeg not found. Installing automatically..."
    
    OS=$(detect_os)
    
    case "$OS" in
        darwin)
            if command -v brew &> /dev/null; then
                brew install ffmpeg
            else
                print_error "Homebrew not found. Please install Homebrew first: https://brew.sh/"
                exit 1
            fi
            ;;
        linux)
            if command -v apt &> /dev/null; then
                sudo apt update
                sudo apt install -y ffmpeg
            elif command -v dnf &> /dev/null; then
                sudo dnf install -y ffmpeg
            elif command -v yum &> /dev/null; then
                sudo yum install -y ffmpeg
            elif command -v pacman &> /dev/null; then
                sudo pacman -S --noconfirm ffmpeg
            else
                print_error "Could not detect package manager. Please install FFmpeg manually."
                exit 1
            fi
            ;;
        windows)
            if command -v winget &> /dev/null; then
                winget install ffmpeg
            elif command -v choco &> /dev/null; then
                choco install ffmpeg
            else
                print_error "Please install FFmpeg manually: https://ffmpeg.org/download.html"
                exit 1
            fi
            ;;
    esac
    
    print_success "FFmpeg installed successfully"
}

install_go() {
    print_step "Go not found. Installing automatically..."
    
    OS=$(detect_os)
    ARCH=$(detect_arch)
    
    case "$OS" in
        darwin)
            if command -v brew &> /dev/null; then
                brew install go
            else
                print_error "Homebrew not found. Please install Homebrew first: https://brew.sh/"
                exit 1
            fi
            ;;
        linux)
            local go_version="1.21.5"
            local go_arch="amd64"
            if [[ "$ARCH" == "arm64" ]]; then
                go_arch="arm64"
            fi
            
            print_step "Downloading Go $go_version..."
            wget -q "https://go.dev/dl/go$go_version.linux-$go_arch.tar.gz" -O /tmp/go.tar.gz
            sudo rm -rf /usr/local/go
            sudo tar -C /usr/local -xzf /tmp/go.tar.gz
            rm /tmp/go.tar.gz
            
            # Add to PATH
            if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
                echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
                echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
            fi
            export PATH=$PATH:/usr/local/go/bin
            ;;
        windows)
            if command -v winget &> /dev/null; then
                winget install GoLang.Go
            elif command -v choco &> /dev/null; then
                choco install golang
            else
                print_error "Please install Go manually: https://go.dev/dl/"
                exit 1
            fi
            ;;
    esac
    
    print_success "Go installed successfully"
}

# Check if git is installed and auto-install
install_git() {
    if ! command -v git &> /dev/null; then
        print_step "Git not found. Installing automatically..."
        
        OS=$(detect_os)
        
        case "$OS" in
            darwin)
                if command -v brew &> /dev/null; then
                    brew install git
                else
                    print_error "Homebrew not found. Please install Homebrew first."
                    exit 1
                fi
                ;;
            linux)
                if command -v apt &> /dev/null; then
                    sudo apt update
                    sudo apt install -y git
                elif command -v dnf &> /dev/null; then
                    sudo dnf install -y git
                elif command -v yum &> /dev/null; then
                    sudo yum install -y git
                elif command -v pacman &> /dev/null; then
                    sudo pacman -S --noconfirm git
                else
                    print_error "Could not install Git. Please install manually."
                    exit 1
                fi
                ;;
            windows)
                if command -v winget &> /dev/null; then
                    winget install Git.Git
                elif command -v choco &> /dev/null; then
                    choco install git
                else
                    print_error "Please install Git manually: https://git-scm.com/download/win"
                    exit 1
                fi
                ;;
        esac
        
        print_success "Git installed successfully"
    fi
}

# Check and install all dependencies
check_and_install_deps() {
    print_step "Checking dependencies..."
    
    # Git
    install_git
    
    # FFmpeg
    if ! command -v ffmpeg &> /dev/null; then
        install_ffmpeg
    else
        print_success "FFmpeg found"
    fi
    
    # Go
    if ! command -v go &> /dev/null; then
        install_go
        # Source bashrc to get go in PATH
        if [[ -f ~/.bashrc ]]; then
            source ~/.bashrc
        fi
        # Recheck if go is now available
        if ! command -v go &> /dev/null; then
            export PATH=$PATH:/usr/local/go/bin
            export PATH=$PATH:$(go env GOPATH 2>/dev/null || echo ~/go)/bin
        fi
    else
        print_success "Go found"
    fi
}

# Install from source
install_from_source() {
    local tmp_dir
    tmp_dir=$(mktemp -d)
    
    print_step "Cloning mediax repository..."
    git clone --depth 1 "https://github.com/$REPO.git" "$tmp_dir" 2>&1 | tail -n 5
    
    cd "$tmp_dir"
    
    print_step "Downloading dependencies..."
    go mod download 2>&1 | tail -n 3
    
    print_step "Building mediax..."
    go build -ldflags="-s -w" -o "$BINARY_NAME" . 2>&1 | tail -n 3
    
    # Install
    if [[ "$OS" == "windows" ]]; then
        local install_path="$HOME/AppData/Local/Programs/mediax"
        mkdir -p "$install_path"
        cp "$BINARY_NAME.exe" "$install_path/"
        print_info "Added to $install_path"
        print_info "Please add this directory to your PATH if not already done."
    else
        print_step "Installing to $INSTALL_DIR..."
        sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    cd - > /dev/null
    rm -rf "$tmp_dir"
}

# Check if mediax is already installed
check_existing() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        local existing_version
        existing_version=$("$BINARY_NAME" version 2>/dev/null | head -n 1 || echo "unknown")
        print_info "Found existing installation: $existing_version"
        read -p "Do you want to upgrade? [Y/n] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]] && [[ ! -z $REPLY ]]; then
            print_info "Installation cancelled."
            exit 0
        fi
    fi
}

# Main installation
main() {
    echo ""
    echo "███████╗███╗   ███╗███████╗██████╗ ██╗ █████╗ ██╗  ██╗"
    echo "██╔════╝████╗ ████║██╔════╝██╔══██╗██║██╔══██╗╚██╗██╔╝"
    echo "█████╗  ██╔████╔██║█████╗  ██║  ██║██║███████║ ╚███╔╝ "
    echo "██╔══╝  ██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║ ██╔██╗ "
    echo "███████╗██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║██╔╝ ██╗"
    echo "╚══════╝╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝"
    echo "          FFmpeg's Cooler Cousin - Installer"
    echo ""

    OS=$(detect_os)
    ARCH=$(detect_arch)
    print_info "Detected: $OS ($ARCH)"
    print_info "Repository: $REPO"
    echo ""

    check_existing
    check_and_install_deps
    echo ""

    print_step "Installing mediax..."
    install_from_source

    print_success "mediax installed successfully!"
    echo ""
    print_info "Try it now:"
    echo "  Interactive TUI:        mediax"
    echo "  List all verbs:         mediax verbs"
    echo "  Verb help:              mediax convert --help"
    echo "  Show version:           mediax version"
    echo ""
    print_info "Quick examples:"
    echo "  mediax probe video.mp4"
    echo "  mediax convert input.mov output.mp4"
    echo "  mediax compress large.mp4 small.mp4 --quality medium"
    echo ""
    print_info "Read the full docs: https://github.com/$REPO"
}

# Run main
main