# AGENTS.md

This file provides strict guidelines for all agents working on the Goto codebase. Follow these instructions carefully.

## Project Overview

Goto is a Go-based CLI tool for managing directory aliases. It is built using Cobra for the CLI interface and Sonic for JSON processing.

## Architecture

*   **Logic Separation**:
    *   `src/core/`: Contains all business logic. This layer must remain framework-agnostic.
    *   `src/cmd/`: Handles CLI interaction using Cobra. Do not place business logic here.
*   **Dependencies**:
    *   **NO NEW DEPENDENCIES**. Use only the existing dependencies defined in `go.mod`.

## Code Style & Standards

*   **Go Version**: Strictly use **Go 1.25**.
*   **Style**: Maintain the existing code style. Ensure consistency with the current codebase.

## Testing

*   **Strategy**: Focus on **end-to-end testing** located in the `tests/` directory.
*   **Execution**: Use the provided bash script for testing and coverage verification:
    ```bash
    ./check-coverage.sh
    ```

## Workflow

*   **Versioning**:
    *   Whenever source code in `src/` is modified, you **MUST** update the `VersionGoto` constant in `src/cmd/version.go`.
*   **Commits & PRs**:
    *   Write clear, descriptive commit messages and Pull Request descriptions.
    *   **NO EMOJIS** in commit messages or PR titles/descriptions.
