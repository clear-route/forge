# Forge Product Documentation

This directory contains comprehensive product documentation for Forge, organized as a collection of feature-focused Product Requirement Documents (PRDs).

## Purpose

Product documentation serves as the **single source of truth** for understanding Forge's capabilities from a product perspective. Each document focuses on:

- **What** the feature does and **why** it exists
- **Who** uses it and **how** they benefit
- **User journeys** and interaction patterns
- **Product strategy** and competitive positioning
- **Success metrics** and validation criteria

## Documentation Structure

### Product Overview
- **[Product Definition](product-definition.md)**: Forge's mission, vision, and core product strategy

### Feature Documentation
Each major feature has its own PRD following a consistent template:

- Agent Loop & Orchestration
- Tool System
- Terminal UI (TUI)
- Context Management
- Auto-Approval System
- Git Integration
- Coding Tools
- Security & Workspace Isolation
- Slash Commands
- LLM Provider Abstraction

## How to Use This Documentation

### For Product Managers
- Understand feature vision and strategy
- Track success metrics and KPIs
- Plan roadmap and prioritization
- Communicate value propositions

### For Designers
- Understand user needs and pain points
- Learn user workflows and journeys
- Identify UX improvement opportunities
- Design new features consistently

### For Engineers
- Understand product requirements and rationale
- See how features connect to user value
- Find links to technical architecture (ADRs)
- Validate implementations against product intent

### For Users/Customers
- Discover available capabilities
- Learn best practices and workflows
- Understand product evolution
- Provide informed feedback

## Relationship to Other Documentation

```
Product Docs (docs/product/)     ← Product vision, user value, requirements
         ↓
         Links to
         ↓
Technical Docs (docs/adr/)       ← Implementation decisions, architecture
         ↓
         Links to
         ↓
Reference Docs (docs/reference/) ← APIs, configuration, technical details
         ↓
         Links to
         ↓
How-To Guides (docs/how-to/)     ← Practical implementation examples
```

## Template

All feature PRDs follow a standard template: **[template.md](template.md)**

This ensures consistency and completeness across all product documentation.

## Contributing

When adding a new feature to Forge:

1. Create a feature PRD using the template
2. Link to relevant ADRs for technical decisions
3. Update the Product Definition if the feature represents a strategic shift
4. Ensure user-facing documentation (how-to guides) references the PRD

## Maintenance

Product documentation should be:
- **Updated** when features evolve significantly
- **Reviewed** quarterly for accuracy
- **Versioned** to track major changes
- **Validated** against actual usage and feedback

---

**Last Updated**: 2025
**Maintained By**: Forge Product Team
