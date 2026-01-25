# Goto

Goto is a "Path Manager" that lets you alias directories with shorter names (abbreviations) or index numbers for quick navigation. These aliases are saved in a JSON file (`goto-paths`).

## Installation

1. **Download** the latest release from [releases](https://github.com/Joacohbc/goto/releases/latest).
2. **Make it executable**: `chmod +x goto`
3. **Initialize**: `./goto init` (creates config, generates aliases, updates shell rc).
4. **Restart terminal**.

See [MANUAL-INSTALL.md](MANUAL-INSTALL.md) for manual setup.

## Usage

### Navigation

```bash
goto home      # Go to path with abbreviation "home"
goto 0         # Go to path at index 0
goto /tmp      # Like regular cd
```

*Note: Abbreviations and indices take precedence over local directory names. Use `-d` to force directory navigation.*

### Manage Paths

**Add Path**
```bash
goto add-path ./ currentDir      # Add current dir as "currentDir"
goto add-path ~/Documents docs   # Add specific path
```

**List Paths**
```bash
goto list
# Output: 0 - "/home/user" - h
```

**Search**
```bash
goto search -a docs    # Search by abbreviation
goto search -p ~/Docs  # Search by path
```

**Delete Path**
```bash
goto delete --path ~/Documents
goto delete --abbv docs
goto delete --indx 2
```

**Modify Path**
Update entries using `goto update <mode>`. Modes combine the *identifier* and the *target* to change (e.g., `path-abbv` means identify by path, update abbreviation).

Shortcuts: `pp` (path-path), `pa` (path-abbv), `pi` (path-indx), `ap` (abbv-path), `aa`, `ai`, `ip`, `ia`, `ii`.

```bash
# Update path (identify by abbreviation 'h')
goto update ap -a h -n /new/path

# Rename abbreviation (identify by path)
goto update pa -p /current/path -n newname
```

### Self-Update
`goto update-goto`

### Backup & Restore
```bash
goto backup [-o file.json]
goto restore [-i file.json]
```

### Temporary Session
Use `-t` flag for temporary paths (cleared on reboot).
```bash
goto add-path -t ./ temp
goto -t temp
```

### Extras
*   `goto -q home` : Return quoted path.
*   `goto -s home` : Return path with escaped spaces.
*   `\cd ~/Documents` : Bypass alias to use standard `cd`.
