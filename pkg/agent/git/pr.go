package git

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type PRContent struct {
	Title       string
	Description string
}

type PRGenerator struct {
	llmClient LLMClient
}

func NewPRGenerator(llmClient LLMClient) *PRGenerator {
	return &PRGenerator{
		llmClient: llmClient,
	}
}

func (g *PRGenerator) Generate(
	ctx context.Context,
	commits []CommitInfo,
	diffSummary string,
	baseBranch string,
	headBranch string,
	customTitle string,
) (*PRContent, error) {
	prompt := g.buildPRPrompt(commits, diffSummary, baseBranch, headBranch, customTitle)

	response, err := g.llmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PR content: %w", err)
	}

	content := parsePRContent(response)

	if customTitle != "" {
		content.Title = customTitle
	}

	return content, nil
}

func (g *PRGenerator) buildPRPrompt(
	commits []CommitInfo,
	diffSummary string,
	base, head string,
	customTitle string,
) string {
	var sb strings.Builder

	sb.WriteString("Generate a pull request title and description based on the following information.\n\n")

	if customTitle != "" {
		sb.WriteString(fmt.Sprintf("User-provided title: %s\n", customTitle))
		sb.WriteString("(You may use this or generate a better one based on the changes)\n\n")
	}

	sb.WriteString(fmt.Sprintf("Base: %s -> Head: %s\n\n", base, head))

	sb.WriteString("Commits:\n")
	for _, commit := range commits {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", commit.Hash, commit.Message))
	}

	sb.WriteString("\nMaterial Changes (from git diff):\n")
	sb.WriteString(diffSummary)

	sb.WriteString("\n\nYou MUST respond with a valid JSON object in this exact format:\n")
	sb.WriteString("{\n")
	sb.WriteString(`  "title": "concise, actionable PR title"` + ",\n")
	sb.WriteString(`  "description": "## Summary\n\n<what changed and why>\n\n## Changes\n\n- <key changes>\n\n## Testing\n\n<how to verify>"` + "\n")
	sb.WriteString("}\n\n")
	sb.WriteString("The description should be in markdown format with sections: Summary, Changes, and Testing.\n")
	sb.WriteString("Respond ONLY with the JSON object, no other text.")

	return sb.String()
}

func parsePRContent(response string) *PRContent {
	// Try to find JSON in the response (handle markdown code blocks)
	jsonStr := response

	// Remove markdown code fences if present
	if idx := strings.Index(response, "```json"); idx != -1 {
		jsonStr = response[idx+7:] // Skip "```json"
		if endIdx := strings.Index(jsonStr, "```"); endIdx != -1 {
			jsonStr = jsonStr[:endIdx]
		}
	} else if idx := strings.Index(response, "```"); idx != -1 {
		jsonStr = response[idx+3:]
		if endIdx := strings.Index(jsonStr, "```"); endIdx != -1 {
			jsonStr = jsonStr[:endIdx]
		}
	}

	// Find JSON object boundaries
	start := strings.Index(jsonStr, "{")
	end := strings.LastIndex(jsonStr, "}")

	if start != -1 && end != -1 && end > start {
		jsonStr = jsonStr[start : end+1]
	}

	// Try to unmarshal JSON
	var content PRContent
	if err := json.Unmarshal([]byte(strings.TrimSpace(jsonStr)), &content); err == nil {
		return &content
	}

	// Fallback: if JSON parsing fails, return empty content
	// The overlay will show "Pull Request Preview" and just the commits/changes
	return &PRContent{
		Title:       "",
		Description: "",
	}
}

func DetectBaseBranch(workingDir string) (string, error) {
	baseBranches := []string{"main", "master", "develop"}

	currentBranch, err := getCurrentBranch(workingDir)
	if err != nil {
		return "", err
	}

	for _, base := range baseBranches {
		cmd := exec.Command("git", "rev-parse", "--verify", base)
		cmd.Dir = workingDir
		if err := cmd.Run(); err != nil {
			continue
		}

		cmd = exec.Command("git", "merge-base", base, currentBranch)
		cmd.Dir = workingDir
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			return base, nil
		}
	}

	return "", fmt.Errorf("could not detect base branch")
}

func getCurrentBranch(workingDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func GetCommitsSinceBase(workingDir, base, head string) ([]CommitInfo, error) {
	cmd := exec.Command("git", "log", "--format=%h|%s", fmt.Sprintf("%s..%s", base, head))
	cmd.Dir = workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get commits: %w, stderr: %s", err, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	commits := make([]CommitInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			commits = append(commits, CommitInfo{
				Hash:    parts[0],
				Message: parts[1],
			})
		}
	}

	return commits, nil
}

func GetDiffSummary(workingDir, base, head string) (string, error) {
	cmd := exec.Command("git", "diff", "--stat", fmt.Sprintf("%s...%s", base, head))
	cmd.Dir = workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get diff stats: %w, stderr: %s", err, stderr.String())
	}

	stats := stdout.String()

	cmd = exec.Command("git", "diff", fmt.Sprintf("%s...%s", base, head))
	cmd.Dir = workingDir

	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get diff: %w, stderr: %s", err, stderr.String())
	}

	diffPreview := truncateDiff(stdout.String(), 5000)

	return fmt.Sprintf("Files Changed:\n%s\n\nCode Changes:\n%s", stats, diffPreview), nil
}
