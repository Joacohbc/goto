# Goto 2.0

Goto is a "Path Manager" that allows you to add a specific path with an identifier, this path can be used as an abbreviation or an index number. Those path are automatically save in a json file, the goto-paths (*gpaths*) files. From these files can add, update, delete and list paths and abbreviations. A *gpath* consists of a Path and an Abbreviation to identify the path. A example of a *gpath* in the goto-paths file:

```json
{
    "path": "/home/user", 
    "abbreviation": "home", 
} 
```

## Use Goto to move in the CLI

If you use Goto with cd (e.g. with aliases) you have the ultimate way to move between folders on the command line. It is quick and easy to use and implement . It works via a compiled Go file that returns the corresponding path based on the arguments passed as input. And passes it as an argument to an alias that uses cd on the command line to move to the specified path.

## How to install?

### Use the automatically way

**Note:** *The install.sh is only for Linux x64 bits*

1. **Clone** repository:  

    ```bash
    git clone https://github.com/Joacohbc/goto.git
    cd ./goto
    ```

2. **And use the install scripts**:

    ```bash
    sh ./install.sh
    ```

### Or do it yourself

1. **Clone** this the repository and go there:

    ```bash
    git clone https://github.com/Joacohbc/goto.git
    cd ./goto
    ```

2. **Compile** the code:

    ```bash
    go build -o goto.bin ./src/*.go  
    ```

3. **Create** the config dir:

    ```bash
    mkdir $HOME/.config/goto/
    ```

4. **Move the files** to the config dir and go there:

    ```bash
    cp -r ./* $HOME/.config/goto/
    cd $HOME/.config/goto/
    ```

5. **Add** the next file to your shell file(ex: .bashrc or .zshrc):

    ```bash
    source $HOME/.config/goto/alias.sh >> {YOUR_SHELL_FILE} 
    ```

6. **Give execution permissions to bin files:**

    ```bash
    chmod +x $HOME/.config/goto/bin/*
    ```

    **Note:** *Step 7 only if the $GOTO_FILE (variable) is incorrect or the goto command doesn't work!*
7. To finish the installation you need to change the GOTO_FILE VARIABLE in alias.sh

    ```bash
    #Use your fav text editor: nano, vi, vim, nvim, etc
    vim $HOME/.config/goto/alias.sh
    ```

    In the editor:

    ```bash
    ##ADD THIS FILE TO .bashrc OR .zshrc WITH "SOURCE <ABSOLUTE-PATH-OF-THIS-FILE>"   
    # GOTO_FILE="<ABSOLUTE-PATH-OF-THIS-FILE>"
    GOTO_FILE="$XDG_CONFIG_HOME/goto/goto.bin" #<-- Here put the absolute path of the goto.bin ($HOME/.config/goto/goto.bin)
    ```

## Usage

### Move (cd aliases)

To use the main function of goto:

```bash  
# Move to the destination directory
# "home" is the abbreviation of /home/user
goto home

Output: Go to: /home/user/

# You also can use "0" (that is the default index of the /home/user)
goto 0

Output: Go to: /home/user/
```

Or also you can use goto like cd, use a complete/relative path:

```bash  
goto /home/user/.config/goto

Output: Go to: /home/user/.config/goto
```

**Note**: *goto always gives priority to the abbreviation and index over a path in the current directory. If in the current working directory exists a directory named "scripts" and you put "scripts" goto search first if "scripts" is abbreviation and after search if a valid path*

### Add new path

To add a new *gpath* require a Path and a Abbreviation:

```bash
# This command add the current directory(the "Path") to the gpaths file with
# the abbreviation "currentDir"
goto add --current -abbv currentDir

#The same that 
goto add -c -a currentDir

# To specify the "Path" and "Abbreviation" use:
goto add --path ~/Documents -abbv docs

#The same that 
goto add -p ~/Documents -a docs
```

### List paths

To list all *gpath* of the *gpaths* file:

```bash
goto list

Output: 
0- Path: "/home/user", Abbreviation: "h"
1- Path: "/home/user/.config/goto/", Abbreviation: "config"
2- Path: "/home/user/Documents", Abbreviation: "docs"
...
```

You also can get a specific line of the gpaths file using the following flags:

```bash
# -p to indicate the abbreviation (also you can use -c if you want to pass the current directory)
goto list -p ~/Documents

Output:
2- Path: "/home/user/Documents", Abbreviation: "docs"

# -a to indicate the abbreviation
goto list -a docs

Output:
2- Path: "/home/user/Documents", Abbreviation: "docs"
```

### Delete paths

To delete a *gpath* require a Path or a Abbreviation:

```bash
#I want to delete the path /home/user/Documents
goto delete --path /home/user/Documents

Output: The changes were applied successfully

#You also can use the Abbreviation or the Index
goto delete --abbv docs

#Or to delete the gpath in the index 2
goto delete --indx 2

Output: The changes were applied successfully
```

### Modify paths

To update a *gpath* you can use 9 modes to update, each mode needs two args, the first to identify the goto-path and the second specific to what is to be updated.

The 9 modes are:

- A "Path" and a new "Path" (path-path)
- A "Path" and a new "Abbreviation" (path-abbv)
- A "Path" and a new "Indx" (path-indx)
- A "Abbreviation" and a new "Path" (abbv-path)
- A "Abbreviation" and a new "Abbreviation" (abbv-path)
- A "Abbreviation" and a new "Indx" (abbv-indx)
- A "Index" and a new "Path" (indx-path)
- A "Index" and a new "Abbreviation" (indx-abbv)
- A "Index" and a new "Index" (indx-indx)

```bash
# Update the home of the user using the path to identify the gpath
goto update path-path --path /home/myuser --new /home/mynewuser

# Or "h" the default abbreviation to home directory 
goto update abbv-path --abbv h --new /home/mynewuser

Output: The changes were applied successfully

# Change the abbreviation of the home
goto update path-abbv --path /home/mynewuser --new home

# Or if you want to update the abbreviation of the home
goto update abbv-abbv --abbv h --new home

Output: The changes were applied successfully
```

### Backup and Restore

To make a backup of the configuration file

```bash
# Made a backup of goto-paths in the config directory
goto backup

Output: Backup complete

# If you want to specify the output path
goto backup -o /the/path/file.json.backup

Output: Backup complete
```

To make a restore of the configuration file from a backup

```bash
# Do a restore of goto-paths from a backup in the config directory
goto restore

Output: Restore complete

# If you want to specify the input path
goto restore -i /the/path/file.json.backup

Output: Restore complete
```

## Temporal gpaths

If you want to add a gpath, but only for a while (until shutdown, for example) you can use the temporary flags (-t) which do the adding, deleting, updating and listing of gpaths in/from a temporary gpath file. The temporary gpath file is deleted on reboot.

```bash
# To add you can use exactly the same command to add a normal gpath, with the -t
goto add -t --current -abbv currentDir

# For a temporal gpaths you have to use temporal flag(-t / --temporal)
goto currentDir

Output: Error: the Path "currentDir" do not exist

# You have to use -t to gpaths in the temporal gpath file
goto -t currentDir
```

## Extras

More options besides the goto to move:

```bash
# Return a path with quotes, you need to specify a abbreviation, a number of index or a directory 
goto -q home

Output: "/home/user"

# You can use the Index
goto -q 0

Output: "/home/user"

# Return a path without spaces (" " -> "\ ") you need to specify a abbreviation, a number of index or a directory 
goto -s example

Output: "/home/user/Dir\ with \ Spaces"
```

## IMPORTANT TO KNOW

**If you want to use only cd, not the alias of he goto function, you should use:**

```bash
#This use the command cd and not the alias
\cd ~/Documents
```

Other thing important is that **goto returns 2 in the exist ($? = 2)** when returning the path, this only happens to get the path (default use). This is because to detect in the alias when the user want to move:

```bash
# Without the alias (from the goto binary directly)
goto h
echo $?

Output: 2
```
