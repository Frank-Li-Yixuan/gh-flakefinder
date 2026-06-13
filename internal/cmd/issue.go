package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os/exec"

	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/detect"
	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/report"
)

type issueOptions struct {
	scanOptions
	dryRun bool
	create bool
	label  string
}

func runIssue(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	opts, err := parseIssueFlags(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	data, repo, err := loadScanData(ctx, opts.scanOptions, stderr)
	if err != nil {
		return err
	}
	suspects := detect.FindSuspects(data, detect.Options{
		MinConfidence: opts.minConfidence,
		Workflow:      opts.workflow,
	})
	body := report.Markdown(repo, suspects)
	if !opts.create {
		_, err := fmt.Fprintf(stdout, "%s\nDry run only. No issue was created.\n", body)
		return err
	}
	if repo == "" {
		return fmt.Errorf("--repo is required when using --create")
	}
	return createIssue(ctx, repo, "Flaky CI suspects", body, opts.label, stdout)
}

func parseIssueFlags(args []string, stderr io.Writer) (issueOptions, error) {
	opts := issueOptions{}
	opts.days = 30
	opts.limit = 300
	opts.format = "markdown"
	opts.minConfidence = 0.60
	opts.dryRun = true
	fs := flag.NewFlagSet("issue", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.repo, "repo", "", "repository in owner/repo format")
	fs.IntVar(&opts.days, "days", opts.days, "number of recent days to scan")
	fs.IntVar(&opts.limit, "limit", opts.limit, "maximum workflow runs to fetch")
	fs.StringVar(&opts.workflow, "workflow", "", "workflow name or id to scan")
	fs.StringVar(&opts.format, "format", opts.format, "output format; issue uses markdown")
	fs.Float64Var(&opts.minConfidence, "min-confidence", opts.minConfidence, "minimum confidence threshold")
	fs.StringVar(&opts.fixture, "fixture", "", "path to local fixture JSON")
	fs.BoolVar(&opts.verbose, "verbose", false, "print API diagnostics")
	fs.BoolVar(&opts.dryRun, "dry-run", true, "print the issue body without creating an issue")
	fs.BoolVar(&opts.create, "create", false, "create a GitHub issue")
	fs.StringVar(&opts.label, "label", "", "label to add when creating an issue")
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	if opts.format != "markdown" {
		return opts, fmt.Errorf("unsupported format %q for issue; expected markdown", opts.format)
	}
	if opts.days <= 0 {
		return opts, fmt.Errorf("--days must be greater than zero")
	}
	if opts.limit <= 0 {
		return opts, fmt.Errorf("--limit must be greater than zero")
	}
	return opts, nil
}

func createIssue(ctx context.Context, repo string, title string, body string, label string, stdout io.Writer) error {
	args := []string{
		"api",
		"--method", "POST",
		fmt.Sprintf("/repos/%s/issues", repo),
		"-f", "title=" + title,
		"-f", "body=" + body,
	}
	if label != "" {
		args = append(args, "-f", "labels[]="+label)
	}
	cmd := exec.CommandContext(ctx, "gh", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("create issue failed: %s", stderr.String())
		}
		return fmt.Errorf("create issue failed: %w", err)
	}
	var response struct {
		HTMLURL string `json:"html_url"`
	}
	if err := json.Unmarshal(out, &response); err == nil && response.HTMLURL != "" {
		_, err = fmt.Fprintf(stdout, "Created issue: %s\n", response.HTMLURL)
		return err
	}
	_, err = stdout.Write(out)
	return err
}
