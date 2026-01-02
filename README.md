## GEL

### An Agentic Version Control System

A Git-compatible version control system built from scratch in Go, designed to pioneer AI-enhanced developer workflows.

---

## Documentation

### Core Documents

- **[AGENTS.md](./AGENTS.md)** - Project overview, architecture, and AI vision
- **[ROADMAP.md](./docs/ROADMAP.md)** - Implementation phases and progress

### Technical Specifications

- **[CONFIGURATION_MANAGEMENT.md](./docs/CONFIGURATION_MANAGEMENT.md)** - How to handle user identity and settings (
  comprehensive guide)
- **[CONFIGURATION_QUICK_START.md](./docs/CONFIGURATION_QUICK_START.md)** - Quick reference for implementing config
  system
- **[TESTING_STRATEGY.md](./docs/TESTING_STRATEGY.md)** - Testing approach and guidelines
- **[VALIDATION_STRATEGY.md](./docs/VALIDATION_STRATEGY.md)** - Validation rules and patterns

### Project Planning

- **[PROJECT_ASSESSMENT_AND_AI_FEATURES.md](./docs/PROJECT_ASSESSMENT_AND_AI_FEATURES.md)** - AI features and career
  impact analysis
- **[CAREER_GUIDANCE_7_MONTH_ROADMAP.md](./docs/CAREER_GUIDANCE_7_MONTH_ROADMAP.md)** - Career development plan

---

## Quick Start

```bash
# Initialize a repository
gel init

# Configure your identity (required for commits)
gel config --global user.name "Your Name"
gel config --global user.email "your@email.com"

# Stage a file
gel add file.txt

# Create a tree from staged files
gel write-tree

# Create a commit
gel commit-tree <tree-hash> -m "Initial commit"
```

---

## Current Status

**Implemented** (Phases 1-2 + partial Phase 3):

- ✅ Repository initialization
- ✅ Object storage (blobs, trees, commits)
- ✅ Staging area (index)
- ✅ Basic commands: `init`, `hash-object`, `cat-file`, `add`, `update-index`, `ls-files`, `write-tree`, `read-tree`,
  `ls-tree`, `commit-tree`

**In Progress**:

- 🚧 Configuration management system
- 🚧 High-level `commit` command
- 🚧 Comprehensive testing

**Planned**:

- 📋 References and branches (Phase 4-5)
- 📋 Status and diff (Phase 6)
- 📋 AI-enhanced features (intelligent commit messages, code review, conflict resolution)

---

## Architecture

Gel follows clean architecture principles with clear layer separation:

```
cmd/        - CLI commands (Cobra)
vcs/        - Service layer (business logic)
domain/     - Domain models (Blob, Tree, Commit, Index)
storage/    - Storage abstraction (filesystem operations)
core/       - Utilities (hashing, compression, validation)
```

See [AGENTS.md](./AGENTS.md) for detailed architecture documentation.

---

## License

See [LICENSE](./LICENSE) file

