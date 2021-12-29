# goto-command
 The ultimate way to move between folders in the command line

Goto is a command that can be used like cd, and also allows you to add specific path to move faster, this path can be used like abbreviation or a index number

It works by means of a compiled Go file (goto.bin) that returns the corresponding path based on the arguments passed as input. And passes it as an argument to an alias that uses cd on the command line to move to the specified path
# How to install?

## Use the automatically way:

**Note:** *The install.sh is only for linux 64 bits and 32 bits*

1. **Clone** repository: <br /> 
    ```bash 
    git clone https://github.com/Joacohbc/goto.git
    cd ./goto
    ``` 
2. **And use the install scripts**: <br />
    ```bash 
    sh ./install.sh
    ```
## Or do it yourself:

1. **Clone** this the repository and go there: <br />
    ```bash
    git clone https://github.com/Joacohbc/goto.git
    cd ./goto
    ```
2. **Compile** the code: <br />
    ```bash
    go build -o goto.bin ./main.go ./config.go 
    ```
3. **Create** the config dir: <br />
    ```bash
    mkdir $HOME/.config/goto/
    ```
4. **Move the files** to the config dir and go there: <br />
    ```bash
    cp -r ./* $HOME/.config/goto/
    cd $HOME/.config/goto/
    ```
5. **Add** the next file to your shell file(ex: .bashrc or .zshrc): <br />
    ```bash
    source $HOME/.config/goto/alias.sh >> {YOUR_SHELL_FILE} 
    ```
6. **Give execution permissions to bin files:**
    ```bash
    chmod +x $HOME/.config/goto/bin/*
    ```
7. To finish the instalation you need to change the GOTO_FILE VARIABLE in alias.sh <br />
    **Note:** *This if the GOTO_FILE is incorrect or the goto command doesn't work!*
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


# How to configure it?

The configuration file is created automatically. To add or remove fav directories
of your config file, you only need add/remove the block between "{}" in the .json

```json
{
  "Path": "{THE-PATH}", 
  "Short": "{THE-ABBREVIATION}", 
} 
```
And you need to add a "," after the "}" if it not the last of the list

**Note:** *You can do this with ```goto --add```*

# Usage:

## Help and version information
In the alias.sh there are more options besides the goto to move:
```bash
#Only return the path for the directory with quotes
goto -q /home/joaco

Output: "/home/joaco"
```

Also have ```goto -help``` to print help message and ```goto -v``` to print version   

## Move
To use the main function of goto:
```bash  
#Move to the destination directory
#home is the abreviation of /home/joaco/ in the config.json
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

## Add new path
To add a new path to the config file:
```bash
#The new path will be ~/Wallpaper/ and "w" is the abreviation 
goto --add="~/Wallpaper/,w" 

Output: The changes were applied successfully

#And try the new path 
goto w

Output: Go to: /home/joaco/Wallpaper/
```

## List paths
To list the path of the config file:
```bash
goto -l 

Output: 
0- Path: "/home/joaco", Short: "home"
1- Path: "/home/joaco/.config/goto/", Short: "conf"
2- Path: "/home/joaco/Wallpaper", Short: "w"
...
```

## Delete paths
To delete a path to the config file:
```bash
#I want to delete the path /home/joaco/Wallpaper
goto --del="/home/joaco/Wallpaper"  

Output: The changes were applied successfully
```