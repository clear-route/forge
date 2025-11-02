# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records (ADRs) for the Forge project.

## What is an ADR?

An Architecture Decision Record (ADR) is a document that captures an important architectural decision made along with its context and consequences.

## Why ADRs?

ADRs help us:
- **Document the "why"** behind important decisions
- **Provide context** for future contributors and maintainers
- **Track evolution** of the architecture over time
- **Prevent revisiting** decisions without understanding past context
- **Enable knowledge sharing** across the team

## ADR Format

We follow the format proposed by Michael Nygard with these sections:

1. **Title**: Short descriptive name
2. **Status**: Proposed, Accepted, Deprecated, Superseded
3. **Context**: What is the issue we're facing?
4. **Decision**: What is the change we're proposing/have made?
5. **Consequences**: What becomes easier or harder as a result?

## Naming Convention

ADRs are numbered sequentially and named with the pattern:

```
NNNN-title-with-dashes.md
```

For example:
- `0001-use-xml-for-tool-calls.md`
- `0002-implement-streaming-responses.md`

## Creating a New ADR

1. Copy the template: `cp template.md NNNN-your-decision.md`
2. Increment the number from the last ADR
3. Fill in all sections
4. Submit a PR for review
5. Update status after decision is made

## ADR Lifecycle

```
Proposed → Accepted → [Deprecated or Superseded]
```

- **Proposed**: Under discussion
- **Accepted**: Decision has been made and implemented
- **Deprecated**: No longer recommended but not replaced
- **Superseded**: Replaced by another ADR (link to new ADR)

## Index

<!-- Keep this list updated when adding new ADRs -->

| Number | Title | Status |
|--------|-------|--------|
| [0001](0001-record-architecture-decisions.md) | Record Architecture Decisions | Accepted |
| [0002](0002-xml-format-for-tool-calls.md) | Use XML Format for Tool Calls | Accepted |
| [0003](0003-provider-abstraction-layer.md) | Provider Abstraction Layer | Accepted |
| [0004](0004-agent-content-processing.md) | Agent-Level Content Processing | Accepted |
| [0005](0005-channel-based-agent-communication.md) | Channel-Based Agent Communication | Accepted |
| [0006](0006-self-healing-error-recovery.md) | Self-Healing Error Recovery with Circuit Breaker | Accepted |
| [0007](0007-memory-system-design.md) | Memory System Design | Accepted |

## Resources

- [Architecture Decision Records](https://adr.github.io/)
- [Documenting Architecture Decisions by Michael Nygard](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions)