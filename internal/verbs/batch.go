package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Batch applies a verb to multiple files
type Batch struct{}

func (b *Batch) Name() string {
	return "batch"
}

func (b *Batch) Description() string {
	return "Apply a verb to multiple files in a directory"
}

func (b *Batch) Usage() string {
	return "mediax batch <verb> <input-pattern> <output-pattern> [--args \"--quality high\"]"
}

func (b *Batch) Execute(args []string, flags map[string]string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: %s", b.Usage())
	}

	verb := args[0]
	inputPattern := args[1]
	outputPattern := args[2]

	// Find matching files
	files, err := filepath.Glob(inputPattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %s", inputPattern)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files matching pattern: %s", inputPattern)
	}

	fmt.Printf("Batch processing %d files with verb '%s'...\n\n", len(files), verb)

	// Process each file
	for i, input := range files {
		// Generate output filename
		base := filepath.Base(input)
		ext := filepath.Ext(base)
		name := base[:len(base)-len(ext)]

		output := strings.ReplaceAll(outputPattern, "{name}", name)
		output = strings.ReplaceAll(output, "{ext}", ext)
		output = strings.ReplaceAll(output, "{index}", fmt.Sprintf("%03d", i+1))

		fmt.Printf("[%d/%d] Processing: %s -> %s\n", i+1, len(files), input, output)

		// Build command arguments
		cmdArgs := []string{verb, input, output}

		// Add extra args if provided
		if extraArgs := flags["args"]; extraArgs != "" {
			cmdArgs = append(cmdArgs, extraArgs)
		}

		// Execute the verb
		cmd := exec.Command(os.Args[0], cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", input, err)
		}

		fmt.Println()
	}

	fmt.Printf("Batch processing complete. Processed %d files.\n", len(files))
	return nil
}

func init() {
	Register(&Batch{})
}
