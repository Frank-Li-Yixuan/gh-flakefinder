package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/detect"
	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/githubapi"
	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/report"
)

type scanOptions struct {
	repo          string
	days          int
	limit         int
	workflow      string
	format        string
	minConfidence float64
	fixture       string
	verbose       bool
}

func runScan(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	opts, err := parseScanFlags("scan", args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	data, repo, err := loadScanData(ctx, opts, stderr)
	if err != nil {
		return err
	}
	suspects := detect.FindSuspects(data, detect.Options{
		MinConfidence: opts.minConfidence,
		Workflow:      opts.workflow,
	})
	return writeReport(stdout, opts.format, repo, opts.days, suspects)
}

func parseScanFlags(name string, args []string, stderr io.Writer) (scanOptions, error) {
	opts := scanOptions{
		days:          30,
		limit:         300,
		format:        "table",
		minConfidence: 0.60,
	}
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.repo, "repo", "", "repository in owner/repo format")
	fs.IntVar(&opts.days, "days", opts.days, "number of recent days to scan")
	fs.IntVar(&opts.limit, "limit", opts.limit, "maximum workflow runs to fetch")
	fs.StringVar(&opts.workflow, "workflow", "", "workflow name or id to scan")
	fs.StringVar(&opts.format, "format", opts.format, "output format: table, json, markdown")
	fs.Float64Var(&opts.minConfidence, "min-confidence", opts.minConfidence, "minimum confidence threshold")
	fs.StringVar(&opts.fixture, "fixture", "", "path to local fixture JSON")
	fs.BoolVar(&opts.verbose, "verbose", false, "print API diagnostics")
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	if !validFormat(opts.format) {
		return opts, fmt.Errorf("unsupported format %q; expected table, json, or markdown", opts.format)
	}
	if opts.days <= 0 {
		return opts, fmt.Errorf("--days must be greater than zero")
	}
	if opts.limit <= 0 {
		return opts, fmt.Errorf("--limit must be greater than zero")
	}
	return opts, nil
}

func loadScanData(ctx context.Context, opts scanOptions, stderr io.Writer) (githubapi.DataSet, string, error) {
	if opts.fixture != "" {
		data, err := githubapi.LoadFixture(opts.fixture)
		if err != nil {
			return githubapi.DataSet{}, "", err
		}
		data = data.WithLimit(opts.limit)
		repo := opts.repo
		if repo == "" {
			repo = data.Repository
		}
		if repo == "" {
			repo = "fixture/repo"
		}
		return data, repo, nil
	}

	repo := opts.repo
	if repo == "" {
		inferred, err := inferRepo(ctx)
		if err != nil {
			return githubapi.DataSet{}, "", err
		}
		repo = inferred
	}
	client := githubapi.GHClient{Stderr: stderr}
	data, err := client.Fetch(ctx, repo, githubapi.FetchOptions{
		Days:     opts.days,
		Limit:    opts.limit,
		Workflow: opts.workflow,
		Verbose:  opts.verbose,
	})
	if err != nil {
		return githubapi.DataSet{}, "", err
	}
	return data, repo, nil
}

func writeReport(stdout io.Writer, format string, repo string, days int, suspects []detect.Suspect) error {
	switch format {
	case "table":
		_, err := io.WriteString(stdout, report.Table(repo, days, suspects))
		return err
	case "json":
		out, err := report.JSON(repo, suspects)
		if err != nil {
			return err
		}
		_, err = io.WriteString(stdout, out)
		return err
	case "markdown":
		_, err := io.WriteString(stdout, report.Markdown(repo, suspects))
		return err
	default:
		return fmt.Errorf("unsupported format %q; expected table, json, or markdown", format)
	}
}

func validFormat(format string) bool {
	return format == "table" || format == "json" || format == "markdown"
}
