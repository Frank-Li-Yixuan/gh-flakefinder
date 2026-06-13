package main

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBinaryFixtureScanAndHelp(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gh-flakefinder")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	runCommand(t, "go", "build", "-o", binary, ".")
	fixture := filepath.FromSlash("../../testdata/workflow_runs.json")

	jsonOut := runCommand(t, binary, "scan", "--fixture", fixture, "--format", "json")
	var parsed struct {
		SuspectCount int               `json:"suspect_count"`
		Suspects     []json.RawMessage `json:"suspects"`
	}
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		t.Fatalf("scan JSON did not parse: %v\n%s", err, jsonOut)
	}
	if parsed.SuspectCount != 4 || len(parsed.Suspects) != 4 {
		t.Fatalf("expected 4 suspects, got count=%d len=%d\n%s", parsed.SuspectCount, len(parsed.Suspects), jsonOut)
	}

	emptyOut := runCommand(t, binary, "scan", "--fixture", fixture, "--min-confidence", "0.91", "--format", "json")
	parsed = struct {
		SuspectCount int               `json:"suspect_count"`
		Suspects     []json.RawMessage `json:"suspects"`
	}{}
	if err := json.Unmarshal([]byte(emptyOut), &parsed); err != nil {
		t.Fatalf("empty JSON did not parse: %v\n%s", err, emptyOut)
	}
	if parsed.SuspectCount != 0 || parsed.Suspects == nil || len(parsed.Suspects) != 0 {
		t.Fatalf("expected empty suspects array, got count=%d nil=%v len=%d\n%s", parsed.SuspectCount, parsed.Suspects == nil, len(parsed.Suspects), emptyOut)
	}

	markdown := runCommand(t, binary, "scan", "--fixture", fixture, "--format", "markdown")
	if !strings.Contains(markdown, "https://github.com/owner/repo/actions/runs/1001/job/5001") {
		t.Fatalf("markdown evidence link missing:\n%s", markdown)
	}

	runCommand(t, binary, "scan", "--help")
	runCommand(t, binary, "issue", "--help")
}

func runCommand(t *testing.T, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, string(out))
	}
	return string(out)
}
