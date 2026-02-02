package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const version = "1.0.0"

// Tool represents an AI coding assistant configuration.
type Tool struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Files       []string `json:"-"`
	Directories []string `json:"-"`
}

// Detection represents a detected configuration file or directory.
type Detection struct {
	Tool        *Tool  `json:"-"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"type"`
}

// tools contains all supported AI coding tools and their configuration files.
var tools = []Tool{
	{
		Name:        "Claude Code",
		Description: "Anthropic's CLI for Claude",
		URL:         "https://claude.ai/code",
		Files:       []string{"CLAUDE.md"},
		Directories: []string{".claude"},
	},
	{
		Name:        "Cursor",
		Description: "AI-powered code editor",
		URL:         "https://cursor.com",
		Files:       []string{".cursorrules"},
		Directories: []string{".cursor"},
	},
	{
		Name:        "Windsurf",
		Description: "Codeium's AI IDE",
		URL:         "https://codeium.com/windsurf",
		Files:       []string{".windsurfrules"},
		Directories: []string{".windsurf"},
	},
	{
		Name:        "GitHub Copilot",
		Description: "GitHub's AI pair programmer",
		URL:         "https://github.com/features/copilot",
		Files:       []string{".github/copilot-instructions.md"},
		Directories: []string{},
	},
	{
		Name:        "Aider",
		Description: "AI pair programming in your terminal",
		URL:         "https://aider.chat",
		Files:       []string{".aider.conf.yml", ".aiderignore", "CONVENTIONS.md"},
		Directories: []string{".aider"},
	},
	{
		Name:        "Cline",
		Description: "AI coding assistant for VS Code",
		URL:         "https://github.com/cline/cline",
		Files:       []string{".clinerules"},
		Directories: []string{".clinerules"},
	},
	{
		Name:        "Zed",
		Description: "Zed editor AI configuration",
		URL:         "https://zed.dev",
		Files:       []string{},
		Directories: []string{".zed"},
	},
	{
		Name:        "Continue.dev",
		Description: "Open-source AI code assistant",
		URL:         "https://continue.dev",
		Files:       []string{},
		Directories: []string{".continue"},
	},
	{
		Name:        "Kiro",
		Description: "AWS agentic AI IDE",
		URL:         "https://kiro.dev",
		Files:       []string{},
		Directories: []string{".kiro"},
	},
	{
		Name:        "Gemini CLI",
		Description: "Google's Gemini Code Assist",
		URL:         "https://developers.google.com/gemini-code-assist",
		Files:       []string{"GEMINI.md", "AGENT.md"},
		Directories: []string{".gemini"},
	},
	{
		Name:        "AGENTS.md Standard",
		Description: "Proposed cross-tool agent rules standard",
		URL:         "https://github.com/anthropics/agent-rules",
		Files:       []string{"AGENTS.md"},
		Directories: []string{},
	},
	{
		Name:        "Bolt",
		Description: "StackBlitz AI full-stack development",
		URL:         "https://bolt.new",
		Files:       []string{".bolt"},
		Directories: []string{".bolt"},
	},
	{
		Name:        "Replit Agent",
		Description: "Replit's AI coding agent",
		URL:         "https://replit.com",
		Files:       []string{".replit"},
		Directories: []string{".replit"},
	},
	{
		Name:        "Codex CLI",
		Description: "OpenAI's coding agent CLI",
		URL:         "https://github.com/openai/codex",
		Files:       []string{"codex.md"},
		Directories: []string{".codex"},
	},
	{
		Name:        "Tabnine",
		Description: "AI code completion assistant",
		URL:         "https://tabnine.com",
		Files:       []string{".tabnine.json", "tabnine.yaml"},
		Directories: []string{".tabnine"},
	},
	{
		Name:        "Amazon Q Developer",
		Description: "AWS AI coding assistant",
		URL:         "https://aws.amazon.com/q/developer/",
		Files:       []string{},
		Directories: []string{".amazonq", ".q"},
	},
	{
		Name:        "Sourcegraph Cody",
		Description: "Sourcegraph's AI coding assistant",
		URL:         "https://sourcegraph.com/cody",
		Files:       []string{".cody.json", "cody.json"},
		Directories: []string{".cody"},
	},
	{
		Name:        "Augment Code",
		Description: "Enterprise AI coding assistant",
		URL:         "https://augmentcode.com",
		Files:       []string{},
		Directories: []string{".augment"},
	},
	{
		Name:        "Supermaven",
		Description: "AI code completion with large context",
		URL:         "https://supermaven.com",
		Files:       []string{},
		Directories: []string{".supermaven"},
	},
}

func detectTools(directory string) ([]Detection, error) {
	var detections []Detection

	for i := range tools {
		tool := &tools[i]

		// Check for configuration files
		for _, file := range tool.Files {
			path := filepath.Join(directory, file)
			info, err := os.Stat(path)
			if err == nil && !info.IsDir() {
				detections = append(detections, Detection{
					Tool:        tool,
					Path:        path,
					IsDirectory: false,
				})
			}
		}

		// Check for configuration directories
		for _, dir := range tool.Directories {
			path := filepath.Join(directory, dir)
			info, err := os.Stat(path)
			if err == nil && info.IsDir() {
				detections = append(detections, Detection{
					Tool:        tool,
					Path:        path,
					IsDirectory: true,
				})
			}
		}
	}

	return detections, nil
}

func formatPlain(detections []Detection, directory string) string {
	if len(detections) == 0 {
		return fmt.Sprintf("No AI coding tool configurations detected in %s", directory)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("AI coding tools detected in %s:\n\n", directory))

	// Group by tool
	toolsFound := make(map[string][]Detection)
	var toolOrder []string

	for _, d := range detections {
		if _, exists := toolsFound[d.Tool.Name]; !exists {
			toolOrder = append(toolOrder, d.Tool.Name)
		}
		toolsFound[d.Tool.Name] = append(toolsFound[d.Tool.Name], d)
	}

	slices.Sort(toolOrder)

	for _, toolName := range toolOrder {
		toolDetections := toolsFound[toolName]
		tool := toolDetections[0].Tool

		sb.WriteString(fmt.Sprintf("  %s\n", tool.Name))
		sb.WriteString(fmt.Sprintf("    %s\n", tool.Description))
		sb.WriteString(fmt.Sprintf("    %s\n", tool.URL))
		sb.WriteString("    Files:\n")

		for _, d := range toolDetections {
			relPath, _ := filepath.Rel(directory, d.Path)
			icon := "file"
			if d.IsDirectory {
				icon = "dir "
			}
			sb.WriteString(fmt.Sprintf("      [%s] %s\n", icon, relPath))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// JSONOutput represents the JSON output format.
type JSONOutput struct {
	Directory string         `json:"directory"`
	Detected  bool           `json:"detected"`
	Tools     []JSONToolData `json:"tools"`
}

// JSONToolData represents a tool in JSON output.
type JSONToolData struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Paths       []JSONPathData `json:"paths"`
}

// JSONPathData represents a path in JSON output.
type JSONPathData struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func formatJSON(detections []Detection, directory string) string {
	output := JSONOutput{
		Directory: directory,
		Detected:  len(detections) > 0,
		Tools:     []JSONToolData{},
	}

	// Group by tool
	toolsFound := make(map[string][]Detection)
	var toolOrder []string

	for _, d := range detections {
		if _, exists := toolsFound[d.Tool.Name]; !exists {
			toolOrder = append(toolOrder, d.Tool.Name)
		}
		toolsFound[d.Tool.Name] = append(toolsFound[d.Tool.Name], d)
	}

	slices.Sort(toolOrder)

	for _, toolName := range toolOrder {
		toolDetections := toolsFound[toolName]
		tool := toolDetections[0].Tool

		toolData := JSONToolData{
			Name:        tool.Name,
			Description: tool.Description,
			URL:         tool.URL,
			Paths:       []JSONPathData{},
		}

		for _, d := range toolDetections {
			relPath, _ := filepath.Rel(directory, d.Path)
			pathType := "file"
			if d.IsDirectory {
				pathType = "directory"
			}
			toolData.Paths = append(toolData.Paths, JSONPathData{
				Path: relPath,
				Type: pathType,
			})
		}

		output.Tools = append(output.Tools, toolData)
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	return string(data)
}

func formatCompact(detections []Detection, _ string) string {
	if len(detections) == 0 {
		return "No AI coding tools detected"
	}

	toolNames := make(map[string]bool)
	for _, d := range detections {
		toolNames[d.Tool.Name] = true
	}

	var names []string
	for name := range toolNames {
		names = append(names, name)
	}
	slices.Sort(names)

	return strings.Join(names, ", ")
}

func formatTable(detections []Detection, directory string) string {
	if len(detections) == 0 {
		return "No AI coding tools detected"
	}

	var sb strings.Builder
	sb.WriteString("Tool                   | Type | Path\n")
	sb.WriteString("-----------------------|------|-----\n")

	// Sort detections
	slices.SortFunc(detections, func(a, b Detection) int {
		if a.Tool.Name != b.Tool.Name {
			return strings.Compare(a.Tool.Name, b.Tool.Name)
		}
		return strings.Compare(a.Path, b.Path)
	})

	for _, d := range detections {
		toolName := d.Tool.Name
		if len(toolName) > 22 {
			toolName = toolName[:22]
		}
		toolName = fmt.Sprintf("%-22s", toolName)

		typeStr := "file"
		if d.IsDirectory {
			typeStr = "dir "
		}

		relPath, _ := filepath.Rel(directory, d.Path)
		sb.WriteString(fmt.Sprintf("%s | %s | %s\n", toolName, typeStr, relPath))
	}

	return sb.String()
}

func listSupportedTools() string {
	var sb strings.Builder
	sb.WriteString("Supported AI coding tools:\n\n")

	// Sort tools by name
	sortedTools := make([]Tool, len(tools))
	copy(sortedTools, tools)
	slices.SortFunc(sortedTools, func(a, b Tool) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	for _, tool := range sortedTools {
		sb.WriteString(fmt.Sprintf("  %s\n", tool.Name))
		sb.WriteString(fmt.Sprintf("    %s\n", tool.Description))
		sb.WriteString(fmt.Sprintf("    URL: %s\n", tool.URL))

		if len(tool.Files) > 0 {
			sb.WriteString(fmt.Sprintf("    Files: %s\n", strings.Join(tool.Files, ", ")))
		}
		if len(tool.Directories) > 0 {
			sb.WriteString(fmt.Sprintf("    Directories: %s\n", strings.Join(tool.Directories, "/, ")+"/"))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func main() {
	formatFlag := flag.String("f", "plain", "Output format: plain, json, compact, table")
	formatLongFlag := flag.String("format", "plain", "Output format: plain, json, compact, table")
	listFlag := flag.Bool("l", false, "List all supported tools")
	listLongFlag := flag.Bool("list", false, "List all supported tools")
	quietFlag := flag.Bool("q", false, "Quiet mode - only return exit code")
	quietLongFlag := flag.Bool("quiet", false, "Quiet mode - only return exit code")
	versionFlag := flag.Bool("v", false, "Show version")
	versionLongFlag := flag.Bool("version", false, "Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `vibedetector - Detect AI coding assistant configuration files

Usage: vibedetector [options] [directory]

Detects configuration files for various "vibe coding" tools including:
Claude Code, Cursor, Windsurf, GitHub Copilot, Aider, Cline, Zed,
Continue.dev, Kiro, Gemini CLI, and more.

Options:
  -f, --format string   Output format: plain, json, compact, table (default "plain")
  -l, --list            List all supported tools and their configuration files
  -q, --quiet           Quiet mode - only return exit code (0 if detected, 1 if not)
  -v, --version         Show version

Arguments:
  directory             Directory to scan (default: current directory)

Examples:
  vibedetector                    # Scan current directory
  vibedetector /path/to/project   # Scan specific directory
  vibedetector -f json            # Output as JSON
  vibedetector -f compact         # Just list tool names
  vibedetector -l                 # List all supported tools
  vibedetector -q && echo "Found" # Use in scripts

`)
	}

	flag.Parse()

	// Handle version
	if *versionFlag || *versionLongFlag {
		fmt.Printf("vibedetector %s\n", version)
		os.Exit(0)
	}

	// Handle list
	if *listFlag || *listLongFlag {
		fmt.Print(listSupportedTools())
		os.Exit(0)
	}

	// Determine format
	format := *formatFlag
	if *formatLongFlag != "plain" {
		format = *formatLongFlag
	}

	// Determine quiet mode
	quiet := *quietFlag || *quietLongFlag

	// Get directory
	directory := "."
	if flag.NArg() > 0 {
		directory = flag.Arg(0)
	}

	// Resolve to absolute path
	absDir, err := filepath.Abs(directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Check if directory exists
	info, err := os.Stat(absDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Directory does not exist: %s\n", absDir)
		os.Exit(2)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: Not a directory: %s\n", absDir)
		os.Exit(2)
	}

	// Detect tools
	detections, err := detectTools(absDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Handle quiet mode
	if quiet {
		if len(detections) > 0 {
			os.Exit(0)
		}
		os.Exit(1)
	}

	// Format and output
	var output string
	switch format {
	case "json":
		output = formatJSON(detections, absDir)
	case "compact":
		output = formatCompact(detections, absDir)
	case "table":
		output = formatTable(detections, absDir)
	default:
		output = formatPlain(detections, absDir)
	}

	fmt.Println(output)

	if len(detections) > 0 {
		os.Exit(0)
	}
	os.Exit(1)
}
