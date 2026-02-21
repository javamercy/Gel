# Gel Command Roadmap

## Phase 1: Core Foundation

_Basic repository setup and object storage_

- [x] **init** - Initialize a new Gel repository
- [x] **hash-object** - Compute object hash and optionally store in object database
- [x] **cat-file** - Display contents of repository objects
- [x] **config** - Get and set repository options

## Phase 2: Staging Area

_Index management for tracking file changes_

- [x] **update-index** - Register file contents to the index (plumbing)
- [x] **ls-files** - Show information about files in the index
- [x] **add** - Add file contents to the staging area

## Phase 3: Tree Objects

_Directory structure representation_

- [x] **write-tree** - Create a tree object from the current index
- [x] **read-tree** - Read tree information into the index
- [x] **ls-tree** - List the contents of a tree object

## Phase 4: Commits

_Recording changes to the repository_

- [x] **commit-tree** - Create a commit object from a tree (plumbing)
- [x] **commit** - Record changes to the repository
- [x] **log** - Show commit logs

## Phase 5: References

_Managing branches and HEAD_

- [x] **update-ref** - Update the object name stored in a ref (plumbing)
- [x] **symbolic-ref** - Read, modify and delete symbolic refs
- [x] **branch** - List, create, or delete branches
- [x] **switch** - Switch branches
- [x] **restore** - Restore working tree files from index/commit

## Phase 6: Status & Diff

_Understanding repository state_

- [x] **status** - Show working tree status (staged, modified, untracked)
- [x] **diff** - Show changes between commits, index, and working tree
- [x] **show** - Show commit details and diff

## Phase 7: History Navigation

_Moving through commit history_

- [ ] **checkout** - Detached HEAD mode, checkout specific commits
- [ ] **reset** - Reset current HEAD to a specified state
- [ ] **reflog** - Manage reflog information

## Phase 8: Undoing Changes

_Reverting and cleaning_

- [ ] **rm** - Remove files from working tree and index

## Phase 9: Merging & Rebasing

_Combining branches_

- [ ] **merge** - Join two or more development histories

## Phase 11: Remote Operations

_Distributed version control_

- [ ] **remote** - Manage tracked repositories
- [ ] **clone** - Clone a repository
- [ ] **fetch** - Download objects and refs from remote
- [ ] **push** - Update remote refs
- [ ] **pull** - Fetch and integrate with remote

## Phase 12: Advanced Features

_Power user operations_

- [ ] **gc** - Cleanup and optimize repository

---

## Phase 13: AI-Powered Commands

_Intelligent version control assistance_

### Commit Assistance

- [ ] **gel ai commit** - Generate commit message from staged changes
- [ ] **gel ai amend** - Suggest improvements to last commit message

### Code Understanding

- [ ] **gel ai explain** - Explain what a commit or diff does
- [ ] **gel ai summarize** - Summarize changes between two refs
- [ ] **gel ai review** - AI code review of staged changes

### Smart Operations

- [ ] **gel ai resolve** - Suggest merge conflict resolutions
- [ ] **gel ai changelog** - Generate changelog from commit history
- [ ] **gel ai pr** - Generate pull request description

### Repository Intelligence

- [ ] **gel ai search** - Natural language search through history
- [ ] **gel ai patterns** - Detect code patterns and anti-patterns
- [ ] **gel ai hotspots** - Identify frequently changed files
