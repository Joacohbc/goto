# goto-command
 The ultimate way to  between folders

Goto is a command that can be used like cd, and also allows you to specify add path to move faster, this path can be used like abbreviation or a index number

# How to install?

1. Clone this the repository and go there

2. Build the bin:
    go build -o goto ./main.go ./config.go 

3. Create the config dir:  
    mkdir /home/username/.config/goto/

4. Move the files to the config dir and go there: 
    cp ./* /home/username/.config/goto/
    cd /home/username/.config/goto/

5. To finish add the next file to your shell file
    source ./alias.sh >> <SHELL_FILE>