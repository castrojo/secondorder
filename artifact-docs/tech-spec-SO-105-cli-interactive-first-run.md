# Tech Spec: CLI Interactive First-Run Configuration (SO-105)

## Overview
This document specifies the technical implementation of an interactive CLI configuration flow for `secondorder`. When a user runs the application for the first time in a new directory, they will be prompted to select a team template and a default agent runner.

## Goals
- Improve onboarding by making configuration options visible.
- Prevent silent creation of agents with hidden defaults.
- Allow skipping the prompt in non-interactive environments (CI) or when flags are provided.

## Technical Requirements

### 1. TTY Detection
The interactive prompt must only trigger if `stdin` is a TTY.
Implementation: Use `golang.org/x/term` package.
```go
import "golang.org/x/term"
// ...
isTTY := term.IsTerminal(int(os.Stdin.Fd()))
```

### 2. First-Run Condition
The prompt should be shown if and only if ALL of the following are true:
1. The database is empty (`len(agents) == 0`).
2. Stdin is a TTY.
3. Neither `-t` (template) nor `-m` (model/runner) flags were explicitly provided. (If only one is provided, prompt for the other).

### 3. Prompt Logic in `main.go`

#### Flag Tracking
Modify the argument parsing loop in `main.go` to track if flags were explicitly set by the user, distinguishing them from default or environment variable values.

#### `promptFirstRun` Function
Add a function `promptFirstRun(database *db.DB, templateProvided, modelProvided bool) (string, string)` to `cmd/secondorder/main.go`.

**Logic:**
1. Check `database.ListAgents()`. If `len > 0`, return existing values.
2. If `!isTTY`, return default values.
3. If `!templateProvided`:
   - Display the list of templates:
     - 1. **startup** - Founding team: CEO, Engineer, Product, Designer, QA, DevOps
     - 2. **dev-team** - Engineering-focused: leads, backend, frontend, QA
     - 3. **saas** - SaaS product: growth, product, engineering, support
     - 4. **agency** - Agency delivery: account, PM, developers, QA
     - 5. **enterprise** - Larger org with multiple team leads and specialists
     - 6. **blank** - No agents, configure manually
   - Read input. Default to `1` (startup) on Enter.
   - Re-prompt once on invalid input, then use default.
4. If template is "blank", skip the runner prompt and return `("blank", "")`.
5. If `!modelProvided`:
   - Display the list of runners:
     - 1. **claude** - Claude Code (default)
     - 2. **gemini** - Google Gemini
     - 3. **codex** - OpenAI Codex
   - Read input. Default to `1` (claude) on Enter.
   - Re-prompt once on invalid input, then use default.
6. Return the selected `(templateName, defaultModel)`.

### 4. Integration
- Call `promptFirstRun` after opening the database but before `applyStartupTemplate`.
- Update `applyStartupTemplate` to handle `templateName == "blank"` by doing nothing.
- Use `bufio.NewScanner(os.Stdin)` or `bufio.NewReader(os.Stdin)` for reading user input.

### 5. Confirmation Echo
Before starting the server, if a selection was made, print:
`Starting with template=<template> runner=<runner>`

## Dependencies
- `golang.org/x/term` (to be added to `go.mod`)

## Success Criteria
- First run with no flags -> Interactive prompt shown.
- Subsequent runs -> No prompt.
- `secondorder -t saas` -> No template prompt, only runner prompt (if no `-m`).
- `secondorder -t saas -m gemini` -> No prompts.
- `cat config.txt | secondorder` -> No prompts.
