# OpenAI Provider Test

This test program verifies the OpenAI provider implementation works correctly.

## Running Without API Key (Demo Mode)

```bash
go run cmd/test-provider/main.go
```

This shows the provider structure without making actual API calls.

## Running With API Key (Live Test)

### 1. Set your OpenAI API key

```bash
export OPENAI_API_KEY="sk-your-key-here"
```

### 2. (Optional) Set custom base URL

For Azure OpenAI, local models, or other OpenAI-compatible APIs:

```bash
export OPENAI_BASE_URL="https://your-resource.openai.azure.com"
```

Or for local models:

```bash
export OPENAI_BASE_URL="http://localhost:8080/v1"
```

### 3. Run the test

```bash
go run cmd/test-provider/main.go
```

## What The Live Test Does

1. **Creates Provider**: Initializes OpenAI provider with `gpt-4o-mini` model
2. **Shows Model Info**: Displays provider, model name, and capabilities
3. **Tests Streaming**: Sends a simple math question and streams the response chunk-by-chunk
4. **Tests Non-Streaming**: Sends another question using the `Complete()` convenience method
5. **Validates Results**: Ensures both methods work correctly

## Expected Output (Live Test)

```
✅ OpenAI API key found, running live test...

1. Creating OpenAI provider...
   ✓ Provider: openai
   ✓ Model: gpt-4o-mini
   ✓ Supports Streaming: true

2. Creating test conversation...
   ✓ System message set
   ✓ User message: What is 2+2? Just give the number.

3. Testing StreamCompletion()...
   Response: 4
   ✓ Received 1 chunks
   ✓ Full response: "4"

4. Testing Complete() (non-streaming)...
   ✓ Response role: assistant
   ✓ Response content: "6"

✅ All tests passed!

The OpenAI provider is working correctly.
Next steps: Implement DefaultAgent to wrap this provider.
```

## What This Proves

- ✅ Provider correctly connects to OpenAI API (or compatible APIs)
- ✅ Custom base URL support works (Azure OpenAI, local models, etc.)
- ✅ Streaming works and returns `StreamChunk` instances
- ✅ Non-streaming `Complete()` wrapper works
- ✅ Message format conversion works correctly
- ✅ Provider layer is decoupled from Agent events (returns simple chunks, not AgentEvents)

## Using With Different Providers

### Standard OpenAI
```bash
export OPENAI_API_KEY="sk-..."
go run cmd/test-provider/main.go
```

### Azure OpenAI
```bash
export OPENAI_API_KEY="your-azure-key"
export OPENAI_BASE_URL="https://your-resource.openai.azure.com"
go run cmd/test-provider/main.go
```

### Local OpenAI-Compatible API (e.g., LM Studio, Ollama with OpenAI compatibility)
```bash
export OPENAI_API_KEY="local"  # Many local APIs don't validate this
export OPENAI_BASE_URL="http://localhost:8080/v1"
go run cmd/test-provider/main.go
```

## Next Steps

Once this test passes, the next phase is to:
1. Implement `DefaultAgent` that wraps the provider
2. Convert `StreamChunk` → `AgentEvent` in the agent layer
3. Implement the full event loop with multi-turn support