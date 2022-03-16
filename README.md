# goto-command

Goto is a path manager that allows you to add a specific path with an identifier to move faster, this path can be used as an abbreviation or an index number. If you use Goto with cd (e.g. with aliases) you have the ultimate way to move between folders on the command line.

Quick and easy to use and install

It works via a compiled Go file (goto.bin) that returns the corresponding path based on the arguments passed as input. And passes it as an argument to an alias that uses cd on the command line to move to the specified path.

## How to install?

### Use the automatically way

**Note:** *The install.sh is only for linux 64 bits and 32 bits*

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

    **Note:** *The step 7 only if the $GOTO_FILE (variable) is incorrect or the goto command doesn't work!*
7. To finish the instalation you need to change the GOTO_FILE VARIABLE in alias.sh

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

## How to configure it?

The configuration file is created automatically. To add or remove fav directories
of your config file, you only need add/remove the block between "{}" in the .json

```json
{
    "path": "{THE-PATH}", 
    "abbreviation": "{THE-ABBREVIATION}", 
} 
```

And you need to add a "," after the "}" if it not the last of the list

**Note:** *You can do this with ```goto -add``` option*

## Usage

### Move

To use the main function of goto:

```bash  
#Move to the destination directory
#"home" is the abreviation of /home/joaco/ in the config.json
goto home

Output: Go to: /home/joaco

#"0" is the number index of the /home/joaco/ in the config.json
goto 0

Output: Go to: /home/joaco
```

Or also you can use goto like cd, use a complete/relative path:

```bash  
goto /home/joaco/.config/goto

Output: Go to: /home/joaco/.config/goto
```

**Note**: *goto always give priority to the abbreviation and index over a path in the current directory. If in the current working directory exists a directory named "scripts" and you put "script" goto search first if "script" is abbreviation and after search if a valid path*

### Add new path

To add a new path:

```bash
#The new path will be ~/Wallpaper/ and "w" is the abreviation 
goto -add -path ~/Wallpaper/ -abbv w 

Output: The changes were applied successfully

#And try the new path 
goto w

Output: Go to: /home/joaco/Wallpaper/
```

### List paths

To list the path of the config file:

```bash
goto -list

Output: 
0- Path: "/home/joaco", Abbreviation: "h"
1- Path: "/home/joaco/.config/goto/", Abbreviation: "config"
2- Path: "/home/joaco/Wallpaper", Abbreviation: "w"
...
```

### Delete paths

To delete a path:

```bash
#I want to delete the path /home/joaco/Wallpaper
goto -del -path /home/joaco/Wallpaper  

Output: The changes were applied successfully
```

### Modify paths

To modify the abreviation of the path:

```bash
#I want to modify the path /home/joaco/.config/goto/
#The new abreviation will be "conf"
goto -modify -path /home/joaco/.config/goto/ -abbv conf  

Output: The changes were applied successfully
```

**Note:** *The path will be exactly the same taht the path in the config file, so you should use the same path* ```goto -list```

### Backup and Restore

To make a backup of the configuration file

```bash
goto -backup

Output: Backup complete
```

To make a restore of the configuration file from a backup

```bash
goto -restore

Output: Restore complete
```

### Extras

More options besides the goto to move:

```bash
#Return a path with quotes, you need to specify a abreviation, a number of index or a directory 
goto -q 

Output: "/home/joaco"
```

Also have ```goto -help``` to print help message and ```goto -v``` to print version  

## IMPORTANT

**If you want to use only cd, not the alias of he goto function, you should use:**

```bash
    #This use the commnad cd and not the alias
    \cd ~/home
```
