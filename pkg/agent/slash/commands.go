package slash

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/entrhq/forge/pkg/agent/git"
)

type Command struct {
	Name string
	Arg  string
}

type Handler struct {
	workingDir      string
	tracker         *git.ModificationTracker
	commitGenerator *git.CommitMessageGenerator
	prGenerator     *git.PRGenerator
	cancelFunc      context.CancelFunc
}

func NewHandler(
	workingDir string,
	tracker *git.ModificationTracker,
	commitGen *git.CommitMessageGenerator,
	prGen *git.PRGenerator,
	cancelFunc context.CancelFunc,
) *Handler {
	return &Handler{
		workingDir:      workingDir,
		tracker:         tracker,
		commitGenerator: commitGen,
		prGenerator:     prGen,
		cancelFunc:      cancelFunc,
	}
}

func Parse(input string) (*Command, bool) {
	trimmed := strings.TrimSpace(input)
	if !strings.HasPrefix(trimmed, "/") {
		return nil, false
	}

	parts := strings.SplitN(trimmed[1:], " ", 2)
	cmd := &Command{
		Name: parts[0],
	}

	if len(parts) > 1 {
		cmd.Arg = strings.TrimSpace(parts[1])
	}

	return cmd, true
}

func (h *Handler) Execute(ctx context.Context, cmd *Command) (string, error) {
	switch cmd.Name {
	case "help":
		return h.handleHelp(), nil
	case "stop":
		return h.handleStop(), nil
	case "commit":
		return h.handleCommit(ctx, cmd.Arg)
	case "pr":
		return h.handlePR(ctx, cmd.Arg)
	default:
		return "", fmt.Errorf("unknown command: /%s", cmd.Name)
	}
}

func (h *Handler) handleHelp() string {
	return "Available Commands:\n" +
		"/help - Show help\n" +
		"/stop - Stop operation\n" +
		"/commit [msg] - Create commit\n" +
		"/pr [title] - Create PR\n"
}

func (h *Handler) handleStop() string {
	if h.cancelFunc != nil {
		h.cancelFunc()
		return "Stopped"
	}
	return "Nothing to stop"
}

func (h *Handler) handleCommit(ctx context.Context, customMessage string) (string, error) {
	files := h.tracker.GetModified()
	if len(files) == 0 {
		return "", fmt.Errorf("no files to commit")
	}

	if err := git.StageFiles(h.workingDir, files); err != nil {
		return "", err
	}

	var message string
	var err error

	if customMessage == "" {
		message, err = h.commitGenerator.Generate(ctx, h.workingDir, files)
		if err != nil {
			return "", err
		}
	} else {
		message = customMessage
	}

	hash, err := git.CreateCommit(h.workingDir, message)
	if err != nil {
		return "", err
	}

	h.tracker.Clear()
	return fmt.Sprintf("Commit %s: %s", hash, message), nil
}

func (h *Handler) handlePR(ctx context.Context, customTitle string) (string, error) {
	base, err := git.DetectBaseBranch(h.workingDir)
	if err != nil {
		return "", err
	}

	head, err := h.getCurrentBranch()
	if err != nil {
		return "", err
	}

	commits, err := git.GetCommitsSinceBase(h.workingDir, base, head)
	if err != nil {
		return "", err
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits for PR")
	}

	diffSummary, err := git.GetDiffSummary(h.workingDir, base, head)
	if err != nil {
		return "", err
	}

	prContent, err := h.prGenerator.Generate(ctx, commits, diffSummary, base, head, customTitle)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("PR (%s -> %s)\n\n", head, base))
	result.WriteString(fmt.Sprintf("Title: %s\n\n", prContent.Title))
	result.WriteString(prContent.Description)

	return result.String(), nil
}

func (h *Handler) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = h.workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func ShouldIntercept(input string) bool {
	_, ok := Parse(input)
	return ok
}
