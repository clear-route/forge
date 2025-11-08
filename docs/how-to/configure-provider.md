# How to Configure LLM Providers

Step-by-step guide to configuring and using LLM providers in Forge.

## Overview

Providers connect your agent to LLM services like OpenAI, Anthropic, or local models. This guide covers provider setup, configuration, and advanced usage.

**Time to complete:** 10 minutes

**What you'll learn:**
- Set up OpenAI provider
- Use OpenAI-compatible providers
- Configure provider options
- Handle provider errors
- Use streaming responses

---

## Prerequisites

- Go 1.21 or later
- API key for your chosen provider
- Basic understanding of environment variables

---

## Step 1: Get an API Key

### OpenAI

1. Visit [OpenAI Platform](https://platform.openai.com/)
2. Sign up or log in
3. Navigate to API Keys section
4. Click "Create new secret key"
5. Copy and save the key (starts with `sk-`)

### Alternative Providers

**Anyscale:**
1. Visit [Anyscale](https://www.anyscale.com/)
2. Sign up for endpoints
3. Get API key from dashboard

**Together AI:**
1. Visit [Together AI](https://www.together.ai/)
2. Create account
3. Generate API key

---

## Step 2: Set Environment Variable

Store your API key securely in environment variables:

### macOS/Linux

```bash
# Add to ~/.bashrc or ~/.zshrc
export OPENAI_API_KEY="sk-your-key-here"

# Or set for current session
export OPENAI_API_KEY="sk-your-key-here"
```

### Windows (PowerShell)

```powershell
$env:OPENAI_API_KEY="sk-your-key-here"
```

### Windows (Command Prompt)

```cmd
set OPENAI_API_KEY=sk-your-key-here
```

---

## Step 3: Create Basic Provider

### Using OpenAI

```go
package main

import (
    "log"
    "os"
    
    "github.com/yourusername/forge/pkg/provider/openai"
)

func main() {
    // Get API key from environment
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY not set")
    }
    
    // Create provider with default settings
    provider := openai.NewProvider("gpt-4", apiKey)
    
    // Use provider with agent...
}
```

### Choosing a Model

```go
// GPT-4 Turbo - Best for complex reasoning
provider := openai.NewProvider("gpt-4-turbo-preview", apiKey)

// GPT-4 - High quality
provider := openai.NewProvider("gpt-4", apiKey)

// GPT-3.5 Turbo - Fast and economical
provider := openai.NewProvider("gpt-3.5-turbo", apiKey)
```

---

## Step 4: Configure Provider Options

### Set Temperature

Controls randomness (0.0 = deterministic, 2.0 = very random):

```go
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithTemperature(0.7), // Balanced (default)
)

// For code generation (consistent)
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithTemperature(0.2),
)

// For creative writing (varied)
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithTemperature(1.2),
)
```

### Set Maximum Tokens

Limit response length:

```go
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithMaxTokens(1000), // Max 1000 tokens in response
)
```

### Set Timeout

Control how long to wait for responses:

```go
import "time"

provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithTimeout(60 * time.Second), // 60 second timeout
)
```

### Combine Options

```go
provider := openai.NewProvider(
    "gpt-4-turbo-preview",
    apiKey,
    openai.WithTemperature(0.7),
    openai.WithMaxTokens(2000),
    openai.WithTimeout(120*time.Second),
)
```

---

## Step 5: Use OpenAI-Compatible Providers

Many services offer OpenAI-compatible APIs. Use them by setting a custom base URL.

### Anyscale

```go
provider := openai.NewProvider(
    "meta-llama/Llama-2-70b-chat-hf",
    apiKey,
    openai.WithBaseURL("https://api.endpoints.anyscale.com/v1"),
)
```

### Together AI

```go
provider := openai.NewProvider(
    "mistralai/Mixtral-8x7B-Instruct-v0.1",
    apiKey,
    openai.WithBaseURL("https://api.together.xyz/v1"),
)
```

### LocalAI (Self-hosted)

```go
provider := openai.NewProvider(
    "gpt-3.5-turbo", // Model name in LocalAI
    "not-needed",    // LocalAI may not need API key
    openai.WithBaseURL("http://localhost:8080/v1"),
)
```

### Ollama (Local models)

```go
provider := openai.NewProvider(
    "llama2",
    "not-needed",
    openai.WithBaseURL("http://localhost:11434/v1"),
    openai.WithTimeout(120*time.Second), // Local models can be slower
)
```

---

## Step 6: Use Streaming Responses

Stream responses for better UX:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/yourusername/forge/pkg/provider/openai"
    "github.com/yourusername/forge/pkg/message"
)

func main() {
    provider := openai.NewProvider("gpt-4", apiKey)
    
    messages := []message.Message{
        message.System("You are a helpful assistant"),
        message.User("Write a short story about a robot"),
    }
    
    ctx := context.Background()
    stream, err := provider.Stream(ctx, messages)
    if err != nil {
        log.Fatal(err)
    }
    
    // Display chunks as they arrive
    for chunk := range stream {
        if chunk.Error != nil {
            log.Printf("Stream error: %v", chunk.Error)
            break
        }
        
        // Print chunk immediately (creates typewriter effect)
        fmt.Print(chunk.Content)
    }
    fmt.Println() // New line at end
}
```

---

## Step 7: Handle Provider Errors

### Basic Error Handling

```go
response, err := provider.Complete(ctx, messages)
if err != nil {
    log.Printf("Provider error: %v", err)
    return
}
```

### Detailed Error Handling

```go
import "errors"

response, err := provider.Complete(ctx, messages)
if err != nil {
    var providerErr *provider.ProviderError
    if errors.As(err, &providerErr) {
        switch providerErr.StatusCode {
        case 401:
            log.Fatal("Invalid API key")
        case 429:
            log.Println("Rate limited, waiting...")
            time.Sleep(time.Second)
            // Retry...
        case 500:
            log.Println("Server error, retrying...")
            // Retry with backoff...
        default:
            log.Printf("API error %d: %s", 
                providerErr.StatusCode, 
                providerErr.Message)
        }
    }
    return
}
```

### Retry Logic

```go
func callWithRetry(provider *openai.Provider, messages []message.Message) (*Response, error) {
    ctx := context.Background()
    maxRetries := 3
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        response, err := provider.Complete(ctx, messages)
        
        if err == nil {
            return response, nil
        }
        
        // Check if retryable
        var providerErr *provider.ProviderError
        if errors.As(err, &providerErr) && providerErr.Retryable {
            waitTime := time.Second * time.Duration(1<<uint(attempt))
            log.Printf("Attempt %d failed, waiting %v", attempt+1, waitTime)
            time.Sleep(waitTime)
            continue
        }
        
        return nil, err // Non-retryable error
    }
    
    return nil, fmt.Errorf("failed after %d attempts", maxRetries)
}
```

---

## Advanced Configuration

### Custom HTTP Client

```go
import "net/http"

customClient := &http.Client{
    Timeout: 90 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}

provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithHTTPClient(customClient),
)
```

### Proxy Configuration

```go
import "net/url"

proxyURL, _ := url.Parse("http://proxy.example.com:8080")
customClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithHTTPClient(customClient),
)
```

### Request Headers

```go
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithHeaders(map[string]string{
        "X-Custom-Header": "value",
        "Organization":    "org-id",
    }),
)
```

---

## Configuration Patterns

### Development Configuration

Fast iteration, verbose logging:

```go
provider := openai.NewProvider(
    "gpt-3.5-turbo",      // Cheaper model
    apiKey,
    openai.WithTemperature(0.7),
    openai.WithTimeout(30*time.Second),
    openai.WithMaxTokens(1000),
)
```

### Production Configuration

Robust, optimized:

```go
provider := openai.NewProvider(
    "gpt-4-turbo-preview",
    apiKey,
    openai.WithTemperature(0.5),
    openai.WithMaxTokens(2000),
    openai.WithTimeout(60*time.Second),
    openai.WithRetries(3),
)
```

### Cost-Optimized Configuration

Minimize API costs:

```go
provider := openai.NewProvider(
    "gpt-3.5-turbo",      // Cheapest model
    apiKey,
    openai.WithTemperature(0.3),
    openai.WithMaxTokens(500), // Limit response length
)
```

### High-Quality Configuration

Best results regardless of cost:

```go
provider := openai.NewProvider(
    "gpt-4",              // Best model
    apiKey,
    openai.WithTemperature(0.7),
    openai.WithMaxTokens(4000),
    openai.WithTimeout(120*time.Second),
)
```

---

## Testing Provider Configuration

### Test Basic Connection

```go
func testProvider(provider *openai.Provider) error {
    messages := []message.Message{
        message.User("Say 'hello'"),
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    response, err := provider.Complete(ctx, messages)
    if err != nil {
        return fmt.Errorf("provider test failed: %w", err)
    }
    
    if response.Content == "" {
        return fmt.Errorf("empty response")
    }
    
    fmt.Println("Provider test successful!")
    return nil
}
```

### Test Streaming

```go
func testStreaming(provider *openai.Provider) error {
    messages := []message.Message{
        message.User("Count to 5"),
    }
    
    ctx := context.Background()
    stream, err := provider.Stream(ctx, messages)
    if err != nil {
        return err
    }
    
    chunks := 0
    for chunk := range stream {
        if chunk.Error != nil {
            return chunk.Error
        }
        chunks++
        fmt.Print(chunk.Content)
    }
    
    if chunks == 0 {
        return fmt.Errorf("no chunks received")
    }
    
    fmt.Printf("\nReceived %d chunks\n", chunks)
    return nil
}
```

---

## Troubleshooting

### "Invalid API Key" Error

**Problem:** Authentication fails

**Solutions:**
1. Check API key is correct
2. Verify environment variable is set: `echo $OPENAI_API_KEY`
3. Ensure no extra spaces in key
4. Check key hasn't been revoked

```go
apiKey := os.Getenv("OPENAI_API_KEY")
if apiKey == "" {
    log.Fatal("OPENAI_API_KEY not set")
}
fmt.Printf("Using API key: %s...%s\n", apiKey[:7], apiKey[len(apiKey)-4:])
```

### Timeout Errors

**Problem:** Requests time out

**Solutions:**
1. Increase timeout
2. Use faster model
3. Reduce max tokens
4. Check network connection

```go
// Increase timeout
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithTimeout(120*time.Second),
)
```

### Rate Limit Errors

**Problem:** Too many requests

**Solutions:**
1. Implement retry with backoff
2. Reduce request frequency
3. Upgrade API plan
4. Use multiple API keys (if allowed)

```go
// Implement exponential backoff
for attempt := 0; attempt < 5; attempt++ {
    response, err := provider.Complete(ctx, messages)
    if err == nil {
        break
    }
    
    if isRateLimitError(err) {
        wait := time.Second * time.Duration(1<<uint(attempt))
        time.Sleep(wait)
        continue
    }
    
    return err
}
```

### Connection Errors

**Problem:** Can't connect to API

**Solutions:**
1. Check internet connection
2. Verify base URL is correct
3. Check firewall/proxy settings
4. Test with curl:

```bash
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
```

---

## Best Practices

### 1. Store Keys Securely

```go
// ✅ Good: Environment variable
apiKey := os.Getenv("OPENAI_API_KEY")

// ❌ Bad: Hardcoded
apiKey := "sk-hardcoded-key-here"
```

### 2. Set Reasonable Timeouts

```go
// ✅ Good: Appropriate timeout
openai.WithTimeout(60 * time.Second)

// ❌ Bad: Too short or too long
openai.WithTimeout(1 * time.Second)  // Too short
openai.WithTimeout(10 * time.Minute) // Too long
```

### 3. Choose Right Model

```go
// ✅ Good: Match model to task
// Simple: gpt-3.5-turbo
// Complex: gpt-4

// ❌ Bad: Always use most expensive
provider := openai.NewProvider("gpt-4", apiKey) // For simple tasks
```

### 4. Handle Errors Properly

```go
// ✅ Good: Detailed error handling
response, err := provider.Complete(ctx, messages)
if err != nil {
    var providerErr *provider.ProviderError
    if errors.As(err, &providerErr) {
        // Handle specific error
    }
    return err
}

// ❌ Bad: Ignore errors
response, _ := provider.Complete(ctx, messages)
```

### 5. Use Streaming for Long Responses

```go
// ✅ Good: Stream long content
if generatingLongContent {
    stream, _ := provider.Stream(ctx, messages)
    // Display chunks
}

// ❌ Bad: Complete for long responses
response, _ := provider.Complete(ctx, messages)
// User waits for entire response
```

---

## Complete Example

Full provider configuration with error handling:

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/yourusername/forge/pkg/provider/openai"
    "github.com/yourusername/forge/pkg/message"
)

func main() {
    // Get API key
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable not set")
    }
    
    // Create provider with options
    provider := openai.NewProvider(
        "gpt-4-turbo-preview",
        apiKey,
        openai.WithTemperature(0.7),
        openai.WithMaxTokens(2000),
        openai.WithTimeout(60*time.Second),
    )
    
    // Test provider
    if err := testProvider(provider); err != nil {
        log.Fatal(err)
    }
    
    // Use with agent...
    fmt.Println("Provider configured successfully!")
}

func testProvider(provider *openai.Provider) error {
    messages := []message.Message{
        message.System("You are a helpful assistant"),
        message.User("Say hello and confirm you're working"),
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    response, err := provider.Complete(ctx, messages)
    if err != nil {
        return handleProviderError(err)
    }
    
    fmt.Printf("Provider response: %s\n", response.Content)
    return nil
}

func handleProviderError(err error) error {
    var providerErr *provider.ProviderError
    if errors.As(err, &providerErr) {
        switch providerErr.StatusCode {
        case 401:
            return fmt.Errorf("invalid API key")
        case 429:
            return fmt.Errorf("rate limited - wait and retry")
        case 500, 503:
            return fmt.Errorf("server error - try again later")
        default:
            return fmt.Errorf("API error %d: %s", 
                providerErr.StatusCode, providerErr.Message)
        }
    }
    return err
}
```

---

## Next Steps

- Read [Configuration Reference](../reference/configuration.md) for all options
- Learn [Error Handling](../reference/error-handling.md) best practices
- See [Performance Guide](../reference/performance.md) for optimization
- Check [API Reference](../reference/api-reference.md) for Provider interface

You're now ready to configure providers for any use case!