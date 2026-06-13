$ErrorActionPreference = "Stop"

$Go = if ($env:GO) { $env:GO } else { "go" }
$GoCommand = Get-Command $Go -ErrorAction Stop
$GoBin = Split-Path -Parent $GoCommand.Source
$GoFmt = Join-Path $GoBin "gofmt.exe"
if (!(Test-Path $GoFmt)) {
  $GoFmt = "gofmt"
}

function Invoke-Step {
  param(
    [Parameter(Mandatory = $true)][string]$Label,
    [Parameter(Mandatory = $true)][scriptblock]$Command
  )
  Write-Host "==> $Label"
  & $Command
}

function Get-GoFiles {
  $tracked = git ls-files -- '*.go'
  if ($tracked) {
    return $tracked
  }
  return Get-ChildItem -Recurse -Filter '*.go' |
    Where-Object { $_.FullName -notmatch '\\.git\\' } |
    ForEach-Object { Resolve-Path -Relative $_.FullName }
}

Invoke-Step "gofmt check" {
  $goFiles = @(Get-GoFiles)
  if ($goFiles.Count -gt 0) {
    $unformatted = & $GoFmt -l @goFiles
    if ($unformatted) {
      $unformatted | ForEach-Object { Write-Error "gofmt needed: $_" }
      throw "gofmt check failed"
    }
  }
}

Invoke-Step "go test ./..." { & $Go test ./... }
Invoke-Step "go vet ./..." { & $Go vet ./... }
Invoke-Step "go build -o gh-flakefinder.exe ./cmd/gh-flakefinder" { & $Go build -o gh-flakefinder.exe ./cmd/gh-flakefinder }

Invoke-Step "fixture JSON smoke test" {
  $json = & .\gh-flakefinder.exe scan --fixture testdata/workflow_runs.json --format json
  $parsed = $json | ConvertFrom-Json
  if ($parsed.suspect_count -ne 4) { throw "expected suspect_count 4, got $($parsed.suspect_count)" }
}

Invoke-Step "fixture Markdown smoke test" {
  $markdown = (& .\gh-flakefinder.exe scan --fixture testdata/workflow_runs.json --format markdown) -join "`n"
  if ($markdown -notmatch 'https://github.com/owner/repo/actions/runs/1001/job/5001') {
    throw "markdown evidence link missing"
  }
}

Invoke-Step "empty-result JSON contract" {
  $json = (& .\gh-flakefinder.exe scan --fixture testdata/workflow_runs.json --min-confidence 0.91 --format json) -join "`n"
  if ($json -notmatch '"suspects": \[\]') { throw 'expected "suspects": []' }
}

Invoke-Step "scan help exits 0" { & .\gh-flakefinder.exe scan --help > $null }
Invoke-Step "issue help exits 0" { & .\gh-flakefinder.exe issue --help > $null }
