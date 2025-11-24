---
description: Instructions for working with the MRBuild (Mono Repo Build) project
applyTo: '**'
---

## Project Overview

**MRBuild** is a cross-platform CLI tool written in Go that detects changes in monorepo folders and triggers appropriate builds for affected projects. It was inspired by NX's "affected" command and designed to solve GitHub + Azure DevOps integration challenges with PR validation builds.

### Purpose

- **Problem Solved**: GitHub PR validation builds with Azure DevOps in monorepos cannot selectively run required builds for only changed projects
- **Solution**: MRBuild analyzes git diffs, identifies affected projects, and spawns only necessary builds in parallel
- **Key Benefit**: Single required build in CI/CD that intelligently runs only affected project builds

## Architecture

### Technology Stack

- **Language**: Go 1.19+
- **CLI Framework**: Cobra (command structure) + Viper (configuration)
- **Concurrency**: Worker pool pattern using `gammazero/workerpool`
- **Logging**: Logrus with structured logging
- **Config Formats**: YAML, JSON, or TOML

### Project Structure

```
cmd/                    - CLI command definitions (root, affected, version)
internal/
  affected/            - Core logic for detecting and building affected projects
  config/              - Configuration parsing and management
  models/              - Data models (App, SpawnBuild)
  util/                - Utility functions (file ops, command building)
  constants/           - Application constants
docs/                  - AsciiDoc documentation
build/                 - Build tooling (eirctl, GitHub actions)
testing/integration/   - Integration tests
```

### Key Components

1. **Affected Detection** (`internal/affected/affected.go`)
   - Runs `git --no-pager diff --name-only <branch>` to get changed files
   - Matches changed files against project patterns (regex)
   - Creates `SpawnBuild` objects for affected projects
   - Executes builds using worker pool for concurrency

2. **Configuration** (`internal/config/`)
   - `Config`: Main config struct with Input and Self sections
   - `Project`: Defines name, folder, patterns, build command, env vars, and order
   - Supports environment variable overrides with `MRBUILD_` prefix
   - Default comparison branch: `main`

3. **Models** (`internal/models/`)
   - `App`: Application-level objects (logger, workers)
   - `SpawnBuild`: Encapsulates build execution details (name, directory, command, env, order)

## Configuration File

Projects are defined in `mrbuild.yaml` (or `.json`/`.toml`) at repo root:

```yaml
projects:
  - name: project-name
    folder: src/project-folder        # Relative path from repo root
    patterns:                          # Regex patterns to match changed files
      - ".*\\.tf"
      - ".*\\.go"
    env:                               # Environment variables for build
      STAGE: production
    build:
      cmd: taskctl build               # Command to execute
      folder: .                        # Where to run (. = config dir, empty = project folder)
    order: 1                           # Execution order (lower runs first)
```

### Configuration Attributes

- **name**: Project identifier
- **folder**: Path to project code from repo root
- **patterns**: Array of regex patterns appended to folder for matching (escape `\` as `\\`)
- **env**: Key-value map of environment variables passed to build process
- **build.cmd**: Shell command to execute for builds
- **build.folder**: Working directory (`.` = config file location, empty = project folder)
- **order**: Integer for build sequencing (projects sorted by this value)

## CLI Commands

### `mrbuild affected`

Main command that detects and builds affected projects.

**Flags:**
- `-c, --config`: Path to config file (default: `./mrbuild.yaml`)
- `--branch`: Branch to compare against (default: `main`, env: `MRBUILD_BRANCH`)
- `--datafile`: Path to file with git output (for testing)
- `--workers`: Number of concurrent workers (default: 1, env: `MRBUILD_WORKERS`)
- `--ignore`: Comma-delimited list of projects to skip
- `--dryrun`: Preview without executing builds
- `--cmdlog`: Log all executed commands to `cmdlog.txt`
- `-l, --loglevel`: Logging level (info, debug, trace, etc.)
- `-f, --logformat`: Output format (text or json)
- `--logfile`: Write logs to file

**Example:**
```bash
mrbuild affected -c ./mrbuild.yaml --workers 3 --ignore project1,project2
```

### `mrbuild version`

Displays version information.

## Key Implementation Details

### Change Detection Flow

1. **Get Changed Files**: Execute git diff or read from datafile/stdin
2. **Match Projects**: For each project, test patterns against changed file list
3. **Build Spawns**: Create `SpawnBuild` for each affected project
4. **Sort by Order**: Sort spawns by `Order` field (ascending)
5. **Execute Builds**: Submit to worker pool with configured concurrency

### Worker Pool Pattern

- Uses `gammazero/workerpool` for concurrent build execution
- Workers configured via `--workers` flag (forced to 1 in dryrun mode)
- Each build runs in isolated goroutine with dedicated stdout/stderr
- `StopWait()` blocks until all builds complete

### Input Modes

1. **Git Command** (default): Runs `git --no-pager diff --name-only <branch>`
2. **Datafile**: Reads from `--datafile` path (for testing/debugging)
3. **Pipe**: Reads from stdin if detected via `util.IsInputFromPipe()`

### Sorting Projects by Order

As of the latest implementation, `SpawnBuild` items are sorted by the `Order` field:

```go
sort.Slice(spawns, func(i, j int) bool {
    return spawns[i].Order < spawns[j].Order
})
```

This ensures projects execute in the specified sequence (lower order values run first).

## Development Guidelines

### Code Style

- Use structured logging with `logrus.WithFields()`
- Follow Go standard naming conventions
- Leverage Viper for configuration binding
- Write unit tests for utilities (`*_test.go`)
- Use AsciiDoc for documentation

### Adding New Projects to Config

When defining projects:
- Use specific regex patterns to avoid false positives
- Set `order` field for build dependencies (e.g., infrastructure before apps)
- Test patterns with `--datafile` and `--dryrun` flags
- Use environment variables for stage/region configuration

### Testing Changes

1. **Dry Run**: Test detection logic without executing builds
   ```bash
   mrbuild affected --dryrun --loglevel debug
   ```

2. **Datafile Testing**: Use captured git output
   ```bash
   git --no-pager diff --name-only main > changes.txt
   mrbuild affected --datafile changes.txt --dryrun
   ```

3. **Integration Tests**: Located in `testing/integration/`

### Build Process

- Uses `eirctl` (custom build tool) via `eirctl _compile`
- Compilation script: `build/scripts/Invoke-Compile.ps1`
- Output binary: `mrbuild` (platform-specific)

## Common Operations

### Adding a New Command

1. Create command file in `cmd/` (e.g., `cmd/newcmd.go`)
2. Define `cobra.Command` with Use, Short, Long, Run
3. Register in `init()` with `rootCmd.AddCommand(newCmd)`
4. Bind flags with Viper in `init()`

### Modifying Configuration Schema

1. Update structs in `internal/config/` with `mapstructure` tags
2. Update unmarshaling in `cmd/root.go` `preRun()`
3. Document in `docs/usage.adoc`
4. Update example YAML in documentation

### Enhancing Affected Logic

Core logic in `internal/affected/affected.go`:
- `getFiles()`: Retrieves changed file list
- `getProjects()`: Matches patterns and creates SpawnBuilds
- `Run()`: Orchestrates detection and execution

## Important Notes

- **Platform**: Cross-platform single binary (Windows, Linux, macOS)
- **Git Dependency**: Requires `git` in PATH for default operation
- **Regex Escaping**: In YAML, escape backslashes (`\\.cs` for `.cs` extension)
- **Environment Variables**: Prefix with `MRBUILD_`, replace `.` with `_` (e.g., `MRBUILD_LOG_LEVEL`)
- **Worker Pool**: Always set to 1 in dryrun mode for sequential output
- **Config Location**: Recommended at repo root, but can be anywhere via `-c` flag

## Troubleshooting

- **No Projects Built**: Check regex patterns match changed file paths exactly
- **All Projects Built**: Verify patterns aren't too broad (e.g., `.*` matches everything)
- **Order Not Working**: Ensure `order` field is set in config and builds are being sorted
- **Git Errors**: Confirm git is installed and repo has valid remote branches
- **Worker Timeouts**: Increase workers or check build command hangs

## Related Documentation

- Full docs in `docs/` directory (AsciiDoc format)
- `docs/overview.adoc`: Architecture and process flow
- `docs/usage.adoc`: Complete configuration reference
- `docs/reasons.adoc`: Problem statement and design rationale
- `docs/examples/`: Usage examples (datafile, piping, env vars)
