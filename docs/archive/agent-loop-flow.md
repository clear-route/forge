# Agent Loop Flow Diagram

## High-Level Agent Loop Flow

```mermaid
graph TD
    A[User Input] --> B[Build Dynamic System Prompt]
    B --> C[Assemble Messages]
    C --> D[Call LLM Provider]
    D --> E{Response Contains Tool Call?}
    
    E -->|Yes| F[Parse Tool Call]
    F --> G{Is Loop-Breaking Tool?}
    
    G -->|task_completion| H[Emit Results]
    G -->|ask_question| I[Prompt User]
    G -->|converse| J[Send Message]
    
    G -->|No| K[Execute Tool]
    K --> L[Add Result to History]
    L --> M{Max Iterations?}
    M -->|No| B
    M -->|Yes| N[Emit Error: Max Iterations]
    
    E -->|No| O[Emit No Tool Call Event]
    O --> P[Add Reminder Prompt]
    P --> C
    
    H --> Q[Exit Loop]
    I --> Q
    J --> Q
    N --> Q
```

## Detailed Message Assembly

```mermaid
graph LR
    A[Start Message Assembly] --> B[System Message]
    B --> C[Add User Instructions]
    C --> D[Add Agent Loop Prompt]
    D --> E[Add System Capabilities]
    E --> F[Add Chain of Thought]
    
    F --> G[Conversation History]
    G --> H[Previous User Messages]
    H --> I[Previous Assistant Messages]
    I --> J[Previous Tool Results]
    
    J --> K[Current Turn Message]
    K --> L[Add Tool Schemas]
    L --> M[Add Tool Calling Rules]
    M --> N[Add Planning Context]
    N --> O[Add RAG Context]
    O --> P[Add User Input]
    
    P --> Q[Final Message Array]
```

## Tool Execution Flow

```mermaid
graph TD
    A[LLM Response Received] --> B[Parse for tool XML]
    B --> C{Tool Call Found?}
    
    C -->|No| D[Emit NoToolCall Event]
    D --> E[Return to Loop]
    
    C -->|Yes| F[Extract Tool Data]
    F --> G[Validate JSON]
    G --> H{Valid?}
    
    H -->|No| I[Emit Error Event]
    I --> E
    
    H -->|Yes| J[Get server_name]
    J --> K{Server Type?}
    
    K -->|local| L[Lookup in Local Registry]
    K -->|mcp| M[Lookup in MCP Registry]
    
    L --> N{Tool Found?}
    M --> N
    
    N -->|No| O[Emit ToolNotFound Error]
    O --> E
    
    N -->|Yes| P[Emit ToolCall Event]
    P --> Q[Execute Tool]
    Q --> R{Success?}
    
    R -->|Yes| S[Emit ToolResult Event]
    R -->|No| T[Emit ToolResultError Event]
    
    S --> U{Loop-Breaking?}
    T --> E
    
    U -->|Yes| V[Exit Loop]
    U -->|No| W[Add to History]
    W --> E
```

## Memory Management Flow

```mermaid
graph TD
    A[New Message] --> B[Add to Memory]
    B --> C{Check Token Count}
    
    C -->|Under Limit| D[Continue]
    
    C -->|Over Limit| E[Prune Strategy]
    E --> F{Keep System Messages}
    F --> G[Remove Oldest User/Assistant]
    G --> H{Still Over Limit?}
    
    H -->|Yes| G
    H -->|No| D
    
    D --> I[Return Messages for Context]
```

## Event Emission Timeline

```mermaid
sequenceDiagram
    participant U as User
    participant A as Agent
    participant P as Provider
    participant T as Tool
    
    U->>A: User Input
    A->>A: LoopStart Event
    
    loop Agent Loop
        A->>A: LoopIteration Event
        A->>P: Call LLM
        P->>A: Stream Response
        A->>A: ThinkingStart Event
        A->>A: ThinkingContent Events
        A->>A: ThinkingEnd Event
        
        alt Tool Call
            A->>A: ToolCall Event
            A->>T: Execute Tool
            T->>A: Tool Result
            A->>A: ToolResult Event
            
            alt Loop-Breaking Tool
                A->>A: LoopEnd Event
                A->>U: Final Result
            else Continue Loop
                A->>A: Add to History
            end
        else No Tool Call
            A->>A: NoToolCall Event
            A->>A: Add Reminder
        end
    end
```

## Tool Call Parsing

```mermaid
graph TD
    A[LLM Response Text] --> B[Search for tool XML]
    B --> C{Found tool tags?}
    
    C -->|No| D[Return nil]
    
    C -->|Yes| E[Extract JSON between tags]
    E --> F[Parse JSON]
    F --> G{Valid JSON?}
    
    G -->|No| H[Return ParseError]
    
    G -->|Yes| I[Extract server_name]
    I --> J{server_name present?}
    
    J -->|No| K[Return MissingFieldError]
    
    J -->|Yes| L[Extract tool_name]
    L --> M{tool_name present?}
    
    M -->|No| K
    
    M -->|Yes| N[Extract arguments]
    N --> O[Return ToolCall struct]
```

## Prompt Building Process

```mermaid
graph LR
    A[PromptBuilder] --> B[Static Prompts]
    A --> C[Dynamic Content]
    
    B --> D[System Capabilities]
    B --> E[Agent Loop]
    B --> F[Chain of Thought]
    B --> G[Tool Calling Rules]
    
    C --> H[Tool Schemas]
    C --> I[MCP Tool Schemas]
    C --> J[Planning Context]
    C --> K[RAG Context]
    
    D --> L[Assemble System Message]
    E --> L
    F --> L
    
    H --> M[Assemble Turn Message]
    I --> M
    G --> M
    J --> M
    K --> M
    
    L --> N[Final Messages Array]
    M --> N
```

## Key Decision Points

### When to Exit Loop?

1. **task_completion** tool called → Agent finished task
2. **ask_question** tool called → Agent needs user input
3. **converse** tool called → Agent wants to chat
4. **Max iterations** reached → Safety limit
5. **User cancellation** → User stopped the agent
6. **Critical error** → Unrecoverable failure

### When to Continue Loop?

1. **Regular tool** executed successfully
2. **No tool call** but under max iterations
3. **Recoverable error** with retry logic

### When to Emit NoToolCall?

1. LLM response has no `<tool>` tags
2. Response is pure conversation/thinking
3. Malformed tool call (parsing failed)