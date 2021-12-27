# goto-command
 The ultimate way to move between folders

Goto is a command that can be used like cd, and also allows you to specify add path to move faster, this path can be used like abbreviation or a index number

# How to install?

1. **Clone** this the repository and go there

2. ***Build** the bin: <br />
_____go build -o goto ./main.go ./config.go 

3. **Create** the config dir: <br />
_____mkdir /home/username/.config/goto/

4. **Move the files** to the config dir and go there: <br />
_____cp ./* /home/username/.config/goto/
_____cd /home/username/.config/goto/

5. To finish **add** the next file to your shell file: <br />
_____source ./alias.sh >> {SHELL_FILE}

The configuration file is created automatically. To add or remove fav directories
of your config file, you only need add/remove the block between "{}" in the json

{
  "Path": "<YOU-PATH>",
  "Short": "<YOU-ABBREVIATION>",
}

And you need to add a "," after the "}" if it not the last of the list