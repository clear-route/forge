package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// AskQuestionTool is a loop-breaking tool that allows the agent to ask
// the user a clarifying question when additional information is needed
// to complete the task.
type AskQuestionTool struct{}

// NewAskQuestionTool creates a new ask question tool
func NewAskQuestionTool() *AskQuestionTool {
	return &AskQuestionTool{}
}

// Name returns the tool's identifier
func (t *AskQuestionTool) Name() string {
	return "ask_question"
}

// Description returns a description of what this tool does
func (t *AskQuestionTool) Description() string {
	return "Ask the user a clarifying question when you need additional information to complete the task. " +
		"Use this when you genuinely need user input to proceed. " +
		"The question should be clear and specific about what information you need."
}

// Schema returns the JSON schema for the tool's arguments
func (t *AskQuestionTool) Schema() map[string]interface{} {
	return BaseToolSchema(
		map[string]interface{}{
			"question": map[string]interface{}{
				"type":        "string",
				"description": "A clear, specific question asking for the information needed to proceed with the task.",
			},
			"suggestions": map[string]interface{}{
				"type":        "array",
				"description": "Optional list of 2-4 suggested answers to help the user respond quickly.",
				"items": map[string]interface{}{
					"type": "string",
				},
				"minItems": 0,
				"maxItems": 4,
			},
		},
		[]string{"question"},
	)
}

// Execute runs the tool and returns the question for the user
func (t *AskQuestionTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Question    string   `json:"question"`
		Suggestions []string `json:"suggestions,omitempty"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("invalid arguments for ask_question: %w", err)
	}

	if args.Question == "" {
		return "", fmt.Errorf("question cannot be empty")
	}

	// Format the question with suggestions if provided
	result := args.Question
	if len(args.Suggestions) > 0 {
		result += "\n\nSuggested answers:"
		for i, suggestion := range args.Suggestions {
			result += fmt.Sprintf("\n%d. %s", i+1, suggestion)
		}
	}

	return result, nil
}

// IsLoopBreaking returns true because this tool terminates the agent loop
// and waits for user input
func (t *AskQuestionTool) IsLoopBreaking() bool {
	return true
}
