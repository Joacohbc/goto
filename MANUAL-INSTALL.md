# Manual Installation & Alias Explanation

## Use Goto to move in the CLI

If you use Goto with cd (e.g. with aliases) you have the ultimate way to move between folders on the command line. It is quick and easy to use and implement. It works via a compiled Go file that returns the corresponding path based on the arguments passed as input. And passes it as an argument to an alias that uses cd on the command line to move to the specified path.

```bash
goto() {
    # NOTE: $GOTO_FILE is the path to the compiled Go file

    OUTPUT=$("$GOTO_FILE" $@) # Execute the compiled Go file with the arguments passed to the functions

    # If $? (exit status of the last command) is 2, goto return a path corresponding to the abbreviation passed as argument
    # This indicates that the path exists and can be used to move to it 
    # (The fact that return 2 and not 0 it intensionally to know if the output is a path, an error or other output, such as list, delete/add/update message, etc)
    if [[ $? -eq 2 ]]; then 
        cd "$OUTPUT"   
        echo "Go to:" $OUTPUT
    elif [ $? -eq 1 ]; then # If $? == 1, output is an error message, so, print it and return 1 (error stats)
        echo "$OUTPUT"
        return 1
    else
        # If $? is not 1 or 2 (its probably 0), output is other functionality of goto (e.g: list of gpaths) so, 
        # print it and return 0 (success stats)
        echo "$OUTPUT" 
    fi
}
```

## Manual Installation

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
