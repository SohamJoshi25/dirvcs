# DIRVCS

Link: https://github.com/SohamJoshi25/dirvcs

# Overview

## Introduction

- DIRVCS is a version management application similar to git which can be used to track any file or directory structure changes similar to git but only focusing structure on not focusing on content.
- DIRVCS is a Command Line Interface Application.
- DIRVCS creates a tree structure of your current directory recording its name, modification and a shallow hash which allow to compare between different versions.
- DIRVCS enables user to easily see the changed files, deleted files, created files on terminal along with exporting this changelog to file.
- Very Fast Speed : Benchmarked: 13.6 GB • 231,946 files • 25,984 folders , Indexing Time: ~15.14 seconds

## Objectives

- Make end user understand what changes occurred between last directory checkpoint and current directory.
- Enable End User to manage their files and folders.

## Motivation

- This project is a solution to a personal problem and the motivation is intrinsic.  I had lots of backups lying around and was unable to compare backup versions to each other efficiently. This application proved to be useful for tracking such changes.
- End User are any people with knowledge to use  a CLI app.

# Technical

- Used Golang as the primary programming language because CLI application Development has great support in Go.
- Used Cobra and Cobra CLI to create a Command Line Application with various flags like —help.
- Used Viper to Manage User Configuration for DIRVCS using YAML config file.
- Used JSON to maintain and save directory structure because JSON is best for storing nested data structures.

# Design and Architecture

Logic Flow

```sql
┌────────────────────────┐
│      CLI Interface     │◄────── User Input (e.g. dirvcs init, persist, compare)
└────────┬───────────────┘
         │
         ▼
┌────────────────────────┐
│    Command Parser      │── Dispatches to appropriate command handler
└────────┬───────────────┘
         ▼
┌────────────────────────┐
│   Command Handlers     │   (init, persist, compare, showtree, logs, etc.)
└────────┬───────────────┘
         ▼
┌────────────────────────┐
│ Version Tracker Engine │  ← Core logic to hash files, detect diffs
└────────┬───────────────┘
         ▼
┌────────────────────────┐
│  Metadata Storage      │  ← JSON / YAML 
└────────┬───────────────┘
         ▼
┌────────────────────────┐
│ File System Abstraction│  ← For reading/writing, recursive walking, ignoring patterns
└────────────────────────┘
```

Generated File Structure

```
    .dirvcs/
    ├── trees/
    │   ├── <uuid>.gz
    │   └── treelogs.json
    ├── logs.json
    ├── config.yaml
    └── .ignore
```

# Challenges and Problem Solving

- Saving Directory Structure to compare it for later use → Was Decided that JSON was best way to store file using Go’s marshalling data into JSON and saving it using `os.create()`
- Maintaining a log file for all changes → Implemented Custom Functions to append to log with timestamp.
- Comparing 2 Tree Snapshots from different paths and logging verbose → Recursive Tree Traversal
- Colorful Output on Terminal → Use of Certain Special Characters to change color of terminal.
- Implement  .*ignore*  file for preventing Some

# Project Impact

- The impact of DRIVCS was that it provided Git-like version tracking but for directory structure and metadata — without the complexity of full Git.
- This was helpful in projects where the structure mattered more than content, like deployments, config folders, or zip submissions.
- It helped me and some of my peers quickly snapshot and compare directory changes over time, making it easier to detect issues or unintended changes.

# Use Cases

- **Local backup verification** – Check if any files or folders have changed before syncing to cloud or external storage.
- **Configuration drift detection** – Track changes in system config directories like `/etc/nginx` or `/var/www`.
- **Student assignment tracking** – Capture directory structure snapshots to monitor progress or check for plagiarism.
- **Media/content folder auditing** – Detect structural changes in folders containing large binary files (images, videos, design assets).
- **Pre-deployment folder checks** – Ensure static site or app folder contents match expected versions before release.
- **Template/boilerplate integrity** – Track and compare versions of reused project templates over time.
- **Versioning non-Git projects** – Use in personal or small projects where Git is overkill or not feasible.
- **Automated directory monitoring** – Schedule regular snapshots via cron to detect unplanned structural changes.