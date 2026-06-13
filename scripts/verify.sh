#!/usr/bin/env bash
set -euo pipefail

GO_BIN="${GO:-go}"
GO_PATH="$(command -v "$GO_BIN")"
GOFMT_BIN="${GOFMT:-$(dirname "$GO_PATH")/gofmt}"
if [[ ! -x "$GOFMT_BIN" ]]; then
  GOFMT_BIN="gofmt"
fi

run() {
  echo "==> $*"
  "$@"
}

tracked_go_files() {
  mapfile -t files < <(git ls-files -- '*.go')
  if [[ "${#files[@]}" -eq 0 ]]; then
    mapfile -t files < <(find . -name '*.go' -not -path './.git/*' -print)
  fi
  printf '%s\n' "${files[@]}"
}

echo "==> gofmt check"
mapfile -t go_files < <(tracked_go_files)
if [[ "${#go_files[@]}" -gt 0 ]]; then
  unformatted="$("$GOFMT_BIN" -l "${go_files[@]}")"
  if [[ -n "$unformatted" ]]; then
    echo "$unformatted"
    echo "gofmt check failed" >&2
    exit 1
  fi
fi

run "$GO_BIN" test ./...
run "$GO_BIN" vet ./...
run "$GO_BIN" build -o gh-flakefinder ./cmd/gh-flakefinder

echo "==> fixture JSON smoke test"
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format json | grep '"suspect_count": 4'

echo "==> fixture Markdown smoke test"
./gh-flakefinder scan --fixture testdata/workflow_runs.json --format markdown | grep 'https://github.com/owner/repo/actions/runs/1001/job/5001'

echo "==> empty-result JSON contract"
./gh-flakefinder scan --fixture testdata/workflow_runs.json --min-confidence 0.91 --format json | grep '"suspects": \[\]'

run ./gh-flakefinder scan --help
run ./gh-flakefinder issue --help
