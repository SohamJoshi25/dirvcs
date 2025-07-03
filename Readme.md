# DirVCS

**DirVCS** is a lightweight version control system for directory structures, inspired by Git but tailored for snapshots of file trees.

It allows you to **persist**, **compare**, **export**, **prune**, and **track changes** in directory states easily from the command line.

---

## 🚀 Features

- Snapshot (persist) current directory state
- Compare any two snapshots or snapshot vs working directory
- List and remove old snapshots
- Export directory trees to JSON
- Manage ignored files/folders
- View internal logs
- Easy YAML-based configuration

---

## 📸 Screenshots

> _Add your screenshots here_

\`\`\`
📁 Placeholder for:
- Snapshot comparison output
- Tree structure export
- Example command usage
\`\`\`

---

## 🛠️ Installation

> Requires **Go 1.20+**

```bash
git clone https://github.com/SohamJoshi25/dirvcs
cd dirvcs
go build -o dirvcs
./dirvcs --help
```

---

## 📚 Commands

### `init`

```bash
dirvcs init
```

Initialize a new DirVCS repository in the current directory.

---

### `persist`

```bash
dirvcs persist -m "Initial snapshot"
```

Persist the current directory as a versioned snapshot.

---

### `tree`

```bash
dirvcs tree --list
dirvcs tree --index 1
dirvcs tree --uuid <UUID>
dirvcs tree --remove <UUID>
```

Manage and view persisted directory trees.

---

### `changes`

```bash
dirvcs changes
dirvcs changes --old <UUID>
dirvcs changes --old <UUID> --new <UUID>
```

Compare directory changes between snapshots or vs current state.

---

### `ignore`

```bash
dirvcs ignore node_modules dist
dirvcs ignore --list
dirvcs ignore --remove node_modules
```

Add, list, or remove ignored paths from the \`.ignore\` file.

---

### `logs`

```bash
dirvcs logs
```

View internal operation logs (persist, delete, compare, etc).

---

### `export`

```bash
dirvcs export --path ./tree.json
dirvcs export --uuid <UUID> --path ./snapshot.json
```

Export a directory tree snapshot to JSON.

---

### `config`

```bash
dirvcs config
dirvcs config --set-key treelimit --set-value 50
```

View or update config settings like \`treelimit\`.

---

## ⚙️ Configuration

DirVCS uses a YAML config file located at:

```
$HOME/.dirvcs/config.yaml
```

Example:

```yaml
treelimit: 20
autoCompress: true
```

You can override values using:

```bash
dirvcs config --set-key treelimit --set-value 50
```

---

## 📁 File Structure

```
.dirvcs/
├── snapshots/
├── logs/
├── config.yaml
└── .ignore
```

---

## 📄 License

MIT License © 2025 [Soham Joshi]

---

Built as a side project for backup managment in Go using [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper).
`