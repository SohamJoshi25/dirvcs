# DirVCS

**DirVCS** is a lightweight version control system for directory structures, inspired by Git but optimized for snapshots of file trees.

Designed to be fast and efficient for tracking large directory structures.

> Benchmarked: 13.6 GB • 231,946 files • 25,984 folders  
> Indexing Time: ~15.14 seconds

---

## Features

- Snapshot (persist) current directory state
- Compare any two snapshots or with the current working directory
- List, view, and remove old snapshots
- Export directory trees to JSON format
- Ignore specific files and folders using a .ignore file
- View internal logs of operations
- YAML-based configuration system
- Optimized performance for large-scale directory trees

---

## Installation

Requires **Go 1.20+**
```bash
    git clone https://github.com/SohamJoshi25/dirvcs
    cd dirvcs
    go build -o dirvcs
    ./dirvcs --help
```
---

## Command Reference

### init

Initialize a new DirVCS repository.
```bash
    dirvcs init
```
---

### persist

Create a snapshot of the current directory state.
```bash
    dirvcs persist -m "Initial snapshot"
```
---

### tree

List, view, or remove persisted snapshots.
```bash
    dirvcs tree --list
    dirvcs tree --index 1
    dirvcs tree --uuid <UUID>
    dirvcs tree --remove <UUID>
```
---

### changes

Compare differences between two snapshots, or between a snapshot and the current state.
```bash
    dirvcs changes
    dirvcs changes --old -o <UUID>
    dirvcs changes --old -o <UUID> --new -n <UUID>
    dirvcs changes --old-path <file1.gz> --new-path <file2.gz>
    dirvcs changes --export-simple <changelog.txt>
    dirvcs changes --export-verbose <changelog.json>
    dirvcs changes --print -p
    dirvcs changes --export -e
    dirvcs changes --verbose -v
```
Notes:
- At least one snapshot must exist before running `changes`.
- If both export and print are disabled, output defaults to terminal.

---

### ignore

Manage ignored files and directories.
```bash
    dirvcs ignore node_modules dist
    dirvcs ignore --list
    dirvcs ignore --remove node_modules
```
---

### logs

Display internal DirVCS operation logs.
```bash
    dirvcs logs
```
---

### export

Export a snapshot or current directory tree as JSON.

    dirvcs export --path -p ./tree.json
    dirvcs export --uuid -u <UUID> --path -p ./snapshot.json

---

### config

View or modify configuration settings.
```bash
    dirvcs config
    dirvcs config --set-key treelimit --set-value 50
```
---

## Configuration

Configuration is stored in a YAML file at:
```bash
    .dirvcs/config.yaml
```
Example:
```bash
    changes:
      export: false
      print: true
    indent: "├──"
    treelimit: 20
    verbose: false
```
---

## File Structure
```bash
    .dirvcs/
    ├── trees/
    │   ├── <uuid>.gz
    │   └── treelogs.json
    ├── logs.json
    ├── config.yaml
    └── .ignore
```
---

## License

MIT License © 2025 Soham Joshi

---

Built using Cobra (github.com/spf13/cobra) and Viper (github.com/spf13/viper).
`
