package verbs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Loop repeats video/audio N times or infinitely
type Loop struct{}

func (l *Loop) Name() string {
	return "loop"
}

func (l *Loop) Description() string {
	return "Loop video or audio N times"
}

func (l *Loop) Usage() string {
	return "mediax loop <input> <output> [--times 3 | --duration 60 | --infinite]"
}

func (l *Loop) Execute(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s", l.Usage())
	}

	input, output := args[0], args[1]

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Determine looping method
	times := flags["times"]
	if times == "" {
		times = flags["count"]
	}

	duration := flags["duration"]
	infinite := flags["infinite"] != "" || flags["forever"] != ""

	var cmdArgs []string

	if infinite {
		// Use stream_loop -1 for infinite loop
		cmdArgs = []string{
			"-stream_loop", "-1",
			"-i", input,
			"-c", "copy",
			"-y",
			output,
		}
		fmt.Printf("Looping infinitely: %s -> %s\n", input, output)
	} else if times != "" {
		// Loop N times using stream_loop
		n, _ := strconv.Atoi(times)
		if n < 1 {
			n = 1
		}
		n-- // stream_loop adds N extra loops

		cmdArgs = []string{
			"-stream_loop", strconv.Itoa(n),
			"-i", input,
			"-c", "copy",
			"-y",
			output,
		}
		fmt.Printf("Looping %d times: %s -> %s\n", n+1, input, output)
	} else if duration != "" {
		// Loop until specific duration using filter
		cmdArgs = []string{
			"-stream_loop", "-1",
			"-i", input,
			"-t", duration,
			"-c", "copy",
			"-y",
			output,
		}
		fmt.Printf("Looping to duration %s: %s -> %s\n", duration, input, output)
	} else {
		// Default: loop 3 times
		cmdArgs = []string{
			"-stream_loop", "2",
			"-i", input,
			"-c", "copy",
			"-y",
			output,
		}
		fmt.Printf("Looping 3 times (default): %s -> %s\n", input, output)
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	Register(&Loop{})
}
