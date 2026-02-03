<div align="center">
  <h1>Goto</h1>
  
  <p>
    <strong>Navigate faster, not harder.</strong>
  </p>
  
  <p>
    A lightning-fast, minimalist <strong>Path Manager</strong> CLI tool designed to supercharge your terminal workflow. 
    Alias your most-used directories, index them, and jump between folders instantly—leaving long paths in the past.
  </p>

</div>

## 🛠️ Built With

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#) [![JSON](https://img.shields.io/badge/JSON-000?logo=json&logoColor=fff)](#) [![Bash](https://img.shields.io/badge/Bash-4EAA25?logo=gnubash&logoColor=fff)](#) [![SonarQube Cloud](https://img.shields.io/badge/SonarQube%20Cloud-126ED3?logo=sonarqubecloud&logoColor=fff)](#) [![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?logo=github-actions&logoColor=white)](#) [![Linux](https://img.shields.io/badge/Linux-FCC624?logo=linux&logoColor=black)](#) [![macOS](https://img.shields.io/badge/macOS-000000?logo=apple&logoColor=F0F0F0)](#) [![Windows](https://custom-icon-badges.demolab.com/badge/Windows-0078D6?logo=windows11&logoColor=white)](#)

## Installation

**Download** the latest release from [releases](https://github.com/Joacohbc/goto/releases/latest).

### Download Binary (Recommended)

You can download the binary directly using `curl` or `wget`. Choose the command for your OS and architecture:

#### Linux

##### **AMD64**

```bash
curl -L -o goto https://github.com/Joacohbc/goto/releases/latest/download/goto-linux-amd64
chmod +x goto
```

##### **ARM64**

```bash
curl -L -o goto https://github.com/Joacohbc/goto/releases/latest/download/goto-linux-arm64
chmod +x goto
```

### Run Init

Once downloaded and made executable, run init to set up aliases automatically:

```bash
./goto init
```

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



