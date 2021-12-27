# goto-command
 The ultimate way to move between folders

Goto is a command that can be used like cd, and also allows you to specify add path to move faster, this path can be used like abbreviation or a index number

# How to install?

1. **Clone** this the repository and go there <br />
    ```bash
    git clone https://github.com/Joacohbc/goto.git <br />
    cd ./goto/
    ```

2. **Build** the bin: <br />
    ```bash
    go build -o goto ./main.go ./config.go 
    ```
3. **Create** the config dir: <br />
    ```bash
    mkdir /home/{username}/.config/goto/
    ```

4. **Move the files** to the config dir and go there: <br />
    ```bash
    cp ./* /home/{username}/.config/goto/ <br />
    cd /home/{username}/.config/goto/
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