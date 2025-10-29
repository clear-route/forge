# CLI Chat with Thinking Mode

This example demonstrates the Forge agent framework's thinking content support. The agent shows its reasoning process before providing answers, similar to models like OpenAI's o1 or Claude with extended thinking.

## How It Works

The example uses a system prompt that instructs the model to wrap its reasoning in `<thinking>` tags:

```
<thinking>
Step-by-step reasoning...
</thinking>

Final answer here.
```

The framework automatically:
1. **Parses thinking tags** from the LLM stream using [`ThinkingParser`](../../pkg/llm/parser/thinking.go:1)
2. **Emits separate events** for thinking vs message content
3. **Displays thinking** in `[Thinking...]` blocks (when enabled)
4. **Saves only the answer** to conversation history (thinking is ephemeral)

## Running the Example

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="sk-..."

# Run the chat
go run examples/cli-chat-thinking/main.go
```

## Example Interaction

```
=== Chat with Thinking Mode ===
The AI will show its reasoning process in [Thinking...] blocks
Try asking: 'What is 15 * 23? Show your work.'

You: What is 15 * 23? Show your work.

[Thinking...]
Let me break this down step by step:
- First, I'll multiply 15 × 20 = 300
- Then, I'll multiply 15 × 3 = 45
- Finally, I'll add them: 300 + 45 = 345
[Done thinking]