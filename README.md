# vibedetector

A CLI tool that detects AI coding assistant configuration files in your projects.

As "vibe coding" becomes mainstream, projects often accumulate configuration files for multiple AI coding tools. `vibedetector` scans a directory and reports which AI assistants have been configured, making it easy to understand a project's AI tooling setup at a glance.

## Installation

### Using Go

```bash
go install github.com/VacTube/vibedetector@latest
```

### From Source

```bash
git clone https://github.com/VacTube/vibedetector.git
cd vibedetector
go build -o vibedetector .
```

### Homebrew (coming soon)

```bash
brew install vactube/tap/vibedetector
```

## Usage

```bash
# Scan current directory
vibedetector

# Scan a specific directory
vibedetector /path/to/project

# Output as JSON
vibedetector -f json

# Output as compact list (just tool names)
vibedetector -f compact

# Output as table
vibedetector -f table

# List all supported tools
vibedetector -l

# Quiet mode - only return exit code (useful in scripts)
vibedetector -q && echo "AI tools detected"

# Show version
vibedetector -v
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | AI coding tools detected |
| 1 | No AI coding tools detected |
| 2 | Error (invalid directory, etc.) |

## Output Formats

### Plain (default)

```
AI coding tools detected in /path/to/project:

  Claude Code
    Anthropic's CLI for Claude
    https://claude.ai/code
    Files:
      [file] CLAUDE.md
      [dir ] .claude

  Cursor
    AI-powered code editor
    https://cursor.com
    Files:
      [file] .cursorrules
```

### JSON

```json
{
  "directory": "/path/to/project",
  "detected": true,
  "tools": [
    {
      "name": "Claude Code",
      "description": "Anthropic's CLI for Claude",
      "url": "https://claude.ai/code",
      "paths": [
        {"path": "CLAUDE.md", "type": "file"},
        {"path": ".claude", "type": "directory"}
      ]
    }
  ]
}
```

### Compact

```
Claude Code, Cursor, GitHub Copilot, Windsurf
```

### Table

```
Tool                   | Type | Path
-----------------------|------|-----
Claude Code            | file | CLAUDE.md
Claude Code            | dir  | .claude
Cursor                 | file | .cursorrules
```

## Supported Tools

| Tool | Files | Directories | URL |
|------|-------|-------------|-----|
| **Claude Code** | `CLAUDE.md` | `.claude/` | [claude.ai/code](https://claude.ai/code) |
| **Cursor** | `.cursorrules` | `.cursor/` | [cursor.com](https://cursor.com) |
| **Windsurf** | `.windsurfrules` | `.windsurf/` | [codeium.com/windsurf](https://codeium.com/windsurf) |
| **GitHub Copilot** | `.github/copilot-instructions.md` | — | [github.com/features/copilot](https://github.com/features/copilot) |
| **Aider** | `.aider.conf.yml`, `.aiderignore`, `CONVENTIONS.md` | `.aider/` | [aider.chat](https://aider.chat) |
| **Cline** | `.clinerules` | `.clinerules/` | [github.com/cline/cline](https://github.com/cline/cline) |
| **Zed** | — | `.zed/` | [zed.dev](https://zed.dev) |
| **Continue.dev** | — | `.continue/` | [continue.dev](https://continue.dev) |
| **Kiro** | — | `.kiro/` | [kiro.dev](https://kiro.dev) |
| **Gemini CLI** | `GEMINI.md`, `AGENT.md` | `.gemini/` | [developers.google.com/gemini-code-assist](https://developers.google.com/gemini-code-assist) |
| **AGENTS.md Standard** | `AGENTS.md` | — | [github.com/anthropics/agent-rules](https://github.com/anthropics/agent-rules) |
| **Bolt** | `.bolt` | `.bolt/` | [bolt.new](https://bolt.new) |
| **Replit Agent** | `.replit` | `.replit/` | [replit.com](https://replit.com) |
| **Codex CLI** | `codex.md` | `.codex/` | [github.com/openai/codex](https://github.com/openai/codex) |
| **Tabnine** | `.tabnine.json`, `tabnine.yaml` | `.tabnine/` | [tabnine.com](https://tabnine.com) |
| **Amazon Q Developer** | — | `.amazonq/`, `.q/` | [aws.amazon.com/q/developer](https://aws.amazon.com/q/developer/) |
| **Sourcegraph Cody** | `.cody.json`, `cody.json` | `.cody/` | [sourcegraph.com/cody](https://sourcegraph.com/cody) |
| **Augment Code** | — | `.augment/` | [augmentcode.com](https://augmentcode.com) |
| **Supermaven** | — | `.supermaven/` | [supermaven.com](https://supermaven.com) |

## Use Cases

### CI/CD Integration

Check if a project has AI coding configurations before running AI-assisted code review:

```bash
if vibedetector -q; then
  echo "AI coding tools detected - running AI review checks"
  # run additional checks
fi
```

### Project Auditing

Scan multiple projects to understand AI tool adoption:

```bash
for dir in ~/projects/*; do
  echo "=== $dir ==="
  vibedetector -f compact "$dir" 2>/dev/null || echo "None"
done
```

### JSON Processing with jq

```bash
# Get just the tool names
vibedetector -f json | jq -r '.tools[].name'

# Check if a specific tool is configured
vibedetector -f json | jq -e '.tools[] | select(.name == "Claude Code")' > /dev/null
```

### Pre-commit Hook

Add to `.git/hooks/pre-commit` to document AI tool usage:

```bash
#!/bin/bash
tools=$(vibedetector -f compact 2>/dev/null)
if [ -n "$tools" ]; then
  echo "AI coding tools in use: $tools"
fi
```

## Configuration File Details

### Claude Code

- **CLAUDE.md**: Project-level instructions and context for Claude
- **.claude/**: Directory containing commands, settings, and MCP configurations

### Cursor

- **.cursorrules**: Project-specific rules for AI behavior, code style, and patterns
- **.cursor/**: Directory with rules and additional configuration

### Windsurf

- **.windsurfrules**: Rules file similar to Cursor's format

### GitHub Copilot

- **.github/copilot-instructions.md**: Instructions that Copilot reads for project context

### Aider

- **.aider.conf.yml**: YAML configuration for Aider settings
- **.aiderignore**: Files to exclude from Aider's context
- **CONVENTIONS.md**: Coding conventions loaded via `--read` flag

### Cline

- **.clinerules/**: Directory with Markdown files for project guidelines

### Zed

- **.zed/**: Contains `settings.json` and rules for Zed's AI features

### Continue.dev

- **.continue/**: Contains `config.yaml` or `config.json` with model and rule configurations

### Kiro

- **.kiro/**: Contains steering files (`product.md`, `structure.md`, `tech.md`), specs, and agent configurations

### Gemini CLI

- **GEMINI.md** or **AGENT.md**: Context files for Gemini Code Assist
- **.gemini/**: Contains `settings.json` and additional configuration

### AGENTS.md Standard

- **AGENTS.md**: Proposed cross-tool standard for agent rules (supported by multiple tools)

## Contributing

Contributions are welcome! If you know of an AI coding tool that's missing, please open an issue or PR.

To add a new tool, edit the `tools` slice in `main.go`:

```go
{
    Name:        "New Tool",
    Description: "Description of the tool",
    URL:         "https://newtool.dev",
    Files:       []string{".newtoolrc", "newtool.config"},
    Directories: []string{".newtool"},
},
```

## License

MIT License - see [LICENSE](LICENSE) for details.
