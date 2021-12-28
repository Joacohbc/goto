# goto-command
 The ultimate way to move between folders

Goto is a command that can be used like cd, and also allows you to add specific path to move faster, this path can be used like abbreviation or a index number

# How to install?

**Use install.sh**
    ```bash 
    git clone https://github.com/Joacohbc/goto.git
    cd ./goto/
    sh ./install.sh
    ```
**or do it yourself**

1. **Clone** this the repository and go there <br />
    ```bash
    git clone https://github.com/Joacohbc/goto.git
    cd ./goto/
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
    cp ./* $HOME/.config/goto/
    cd $HOME/.config/goto/
    ```
5. **Add** the next file to your shell file: <br />
    ```bash
    source ./alias.sh >> {SHELL_FILE} 
    ```

6. To finish the instalation you need to change the GOTO_FILE VARIABLE in alias.sh <br />
    **THIS IF THE GOTO_FILE DON'T WORK!**

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