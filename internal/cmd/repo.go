package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func inferRepo(ctx context.Context) (string, error) {
	if repo, err := inferRepoFromGit(ctx); err == nil && repo != "" {
		return repo, nil
	}
	if repo, err := inferRepoFromGH(ctx); err == nil && repo != "" {
		return repo, nil
	}
	return "", fmt.Errorf("could not infer repository; pass --repo owner/repo")
}

func inferRepoFromGit(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return parseGitHubRemote(strings.TrimSpace(string(out)))
}

func inferRepoFromGH(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	repo := strings.TrimSpace(string(out))
	if strings.Count(repo, "/") != 1 {
		return "", fmt.Errorf("unexpected gh repo view output %q", repo)
	}
	return repo, nil
}

func parseGitHubRemote(remote string) (string, error) {
	remote = strings.TrimSuffix(remote, ".git")
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^https://[^/]+/([^/]+)/([^/]+)$`),
		regexp.MustCompile(`^git@[^:]+:([^/]+)/([^/]+)$`),
		regexp.MustCompile(`^ssh://git@[^/]+/([^/]+)/([^/]+)$`),
	}
	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(remote)
		if len(matches) == 3 {
			return matches[1] + "/" + matches[2], nil
		}
	}
	return "", fmt.Errorf("could not parse git remote %q", remote)
}
