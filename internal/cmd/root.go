package cmd

import (
	"context"
	"fmt"
	"io"
)

func Execute(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stdout)
		return 0
	}
	ctx := context.Background()
	switch args[0] {
	case "scan":
		if err := runScan(ctx, args[1:], stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		return 0
	case "issue":
		if err := runIssue(ctx, args[1:], stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		return 0
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "error: unknown command %q\n", args[0])
		printUsage(stderr)
		return 1
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "gh-flakefinder detects likely flaky GitHub Actions workflow/job runs.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  gh flakefinder scan [flags]")
	fmt.Fprintln(w, "  gh flakefinder issue [flags]")
}
