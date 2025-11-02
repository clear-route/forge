# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

## Reporting a Vulnerability

We take the security of Forge seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please do NOT:

- Open a public GitHub issue for security vulnerabilities
- Discuss the vulnerability publicly until it has been addressed

### Please DO:

1. **Report via GitHub Security Advisories**: Use the [Security tab](https://github.com/entrhq/forge/security) in the GitHub repository
2. **Or email us**: If you prefer email, contact the maintainers via GitHub
3. **Provide details**: Include as much information as possible:
   - Type of vulnerability
   - Full paths of source file(s) related to the vulnerability
   - Location of the affected source code (tag/branch/commit or direct URL)
   - Step-by-step instructions to reproduce the issue
   - Proof-of-concept or exploit code (if possible)
   - Impact of the issue, including how an attacker might exploit it

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
- **Assessment**: We will send an assessment of the vulnerability within 7 days
- **Fix Timeline**: We will work on a fix and keep you updated on progress
- **Disclosure**: We will coordinate with you on the disclosure timeline
- **Credit**: We will credit you in the release notes (unless you prefer to remain anonymous)

## Security Best Practices

When using Forge in your applications:

1. **API Keys**: Never commit API keys or sensitive credentials to version control
2. **Input Validation**: Always validate user input before passing to agents
3. **Rate Limiting**: Implement rate limiting when exposing agent functionality via APIs
4. **Tool Execution**: Be cautious when allowing custom tool execution, especially in production
5. **Dependencies**: Keep dependencies up to date and monitor for security advisories
6. **Environment Variables**: Use environment variables for configuration, not hardcoded values

## Known Security Considerations

### LLM-Specific Risks

- **Prompt Injection**: Agents may be vulnerable to prompt injection attacks. Sanitize user inputs
- **Data Exposure**: Be careful about what data is included in prompts sent to LLM providers
- **Tool Abuse**: Custom tools should be carefully reviewed for security implications
- **Cost Control**: Implement cost controls to prevent abuse of LLM API calls

### Framework Security

- **Tool Execution**: Tools execute in the same process - ensure custom tools are trustworthy
- **Memory Management**: Conversation history may contain sensitive information
- **Error Messages**: Error messages may leak information about your system

## Disclosure Policy

When we receive a security bug report, we will:

1. Confirm the problem and determine affected versions
2. Audit code to find similar problems
3. Prepare fixes for all supported versions
4. Release patches as soon as possible

We appreciate your efforts to responsibly disclose your findings and will make every effort to acknowledge your contributions.