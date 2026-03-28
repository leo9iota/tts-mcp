# Documentation Hub

Welcome to the project documentation hub. This directory contains architectural documentation, decision logs, active tasks, and general project overviews.

## Directory Structure

```text
docs/
├── docs.md                  # This file
├── project/                 # High-level project overviews and definitions
│   └── 01-tts-mcp.md
├── adr/                     # Architectural Decision Records
│   └── _ADR.md              # Log of all decisions
├── rfc/                     # Requests for Comments (Technical designs)
│   └── _RFC.md              # Log of all designs
└── tasks/                   # Ephemeral task tracking and checklists
    ├── backend/
    │   ├── active/
    │   │   └── _Backend-Active-Tasks.md
    │   ├── archive/
    │   │   └── _Backend-Archived-Tasks.md
    │   └── backlog/
    │       └── _Backend-Backlog-Tasks.md
    └── ci/
        ├── active/
        │   └── _Ci-Active-Tasks.md
        ├── archive/
        │   └── _Ci-Archived-Tasks.md
        └── backlog/
            └── _Ci-Backlog-Tasks.md
```

## Task Board Workflow

Features and bugs are tracked dynamically using Markdown checklists.

| State       | Location                    | Description                                                 |
| :---------- | :-------------------------- | :---------------------------------------------------------- |
| **Backlog** | `tasks/{{domain}}/backlog/` | Unstarted tasks. Move to Active when starting work.         |
| **Active**  | `tasks/{{domain}}/active/`  | Tasks currently being worked on by a developer or AI agent. |
| **Archive** | `tasks/{{domain}}/archive/` | Completed or canceled tasks. Kept for audit trail.          |
