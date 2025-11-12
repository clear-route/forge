package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type CommitInfo struct {
	Hash    string
	Message string
}

type CommitMessageGenerator struct {
	llmClient LLMClient
}

type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

func NewCommitMessageGenerator(llmClient LLMClient) *CommitMessageGenerator {
	return &CommitMessageGenerator{
		llmClient: llmClient,
	}
}

func (g *CommitMessageGenerator) Generate(ctx context.Context, workingDir string, files []string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("no files to commit")
	}

	diff, err := getDiff(workingDir, files)
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}

	prompt := buildCommitPrompt(diff, files)
	message, err := g.llmClient.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	return strings.TrimSpace(message), nil
}

func getDiff(workingDir string, files []string) (string, error) {
	args := append([]string{"diff", "--cached"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff failed: %w, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func buildCommitPrompt(diff string, files []string) string {
	var sb strings.Builder

	sb.WriteString("Generate a conventional commit message for these changes.\\n\\n")
	sb.WriteString("Format: <type>(<scope>): <description>\\n")
	sb.WriteString("Types: feat, fix, docs, style, refactor, test, chore\\n\\n")

	sb.WriteString("Files changed:\\n")
	for _, file := range files {
		sb.WriteString(fmt.Sprintf("- %s\\n", file))
	}

	sb.WriteString("\\nDiff:\\n")
	sb.WriteString(truncateDiff(diff, 3000))

	sb.WriteString("\\n\\nGenerate ONLY the commit message (one line), nothing else.")

	return sb.String()
}

func truncateDiff(diff string, maxChars int) string {
	if len(diff) <= maxChars {
		return diff
	}
	return diff[:maxChars] + "\\n... (diff truncated)"
}

func StageFiles(workingDir string, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("no files to stage")
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = workingDir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

func CreateCommit(workingDir, message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git commit failed: %w, stderr: %s", err, stderr.String())
	}

	hash, err := getLatestCommitHash(workingDir)
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}

	return hash, nil
}

func getLatestCommitHash(workingDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git rev-parse failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
