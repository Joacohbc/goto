echo "You have go installed?(y/any)"
read op

if [ "$op" = "y" ]; then

    GOTO_BIN="goto.bin"

    rm -f $GOTO_BIN

    #Build the code
    go build -o $GOTO_BIN ./main.go ./config.go 

    if [ $? -eq 0 ]; then
        echo "Compaling successfully"
    else
        echo "Compaling failed"
    fi

else
    #Chek the ARCHITECTURE of the system, 32 bit or 64 bit
    ARCHITECTURE=`uname -m`
    
    if [ "$ARCHITECTURE" = "x86_64" ]; then
        GOTO_BIN="bin/goto64.bin"
    else
        GOTO_BIN="bin/goto32.bin"
    fi

fi

#if $XDG_CONFIG_HOME is empty
if [ -z "$XDG_CONFIG_HOME" ]; then
    XDG_CONFIG_HOME="$HOME/.config"
fi

#Define the config dir
CONFIG_DIR="$XDG_CONFIG_HOME"/goto/

#Create the config dir
mkdir -p $CONFIG_DIR

if [ $? -eq 0 ]; then
    echo "Config dir created successfully"
else
    echo "Config dir created failed"
fi

#Absolute path line /home/username/.bashrc, not ~/.bashrc
# SHELL_FILE="<PUT-PATH>"

while true; do
    echo "Enter the ABSOLUTE PATH of you shell configure file: "
    read SHELL_FILE

    if [ -f "$SHELL_FILE" ]; then
        #Add the alias.sh to $SHELL_FILE
        echo "" >> $SHELL_FILE #New line
        echo "#Aliases to use goto:" >> $SHELL_FILE
        echo "source "$CONFIG_DIR"alias.sh" >> $SHELL_FILE
        echo "" >> $SHELL_FILE #New line

        break
    else
        echo "$SHELL_FILE doesn't exist" 
    fi
done 

#Number of linees of the text
num=$(wc -l alias.sh | cut -d " " -f 1)

#Remainder of 3 first lines of orginal text(the lines of GOTO_FILE variable)
result=$(expr $num - 3)

#Copy all of the repository to CONFIG_DIR
cp -r ./* $CONFIG_DIR

#Give excute permission to the bin files
chmod +x "$CONFIG_DIR/bin/*"

if [ $? -eq 0 ]; then
    echo "All files copied successfully"
else
    echo "All files copied failed"
fi

#---------------------------------------#
#CHANGE CONFIG_FILE VARIABLE IN ALIAS.SH#
#---------------------------------------#

#Add the GOTO_FILE variable
aliasFile=""$CONFIG_DIR"alias.sh.new"

#Put the advise menssages and GOTO_FILE variable in the alias
echo "##ADD THIS FILE TO .bashrc OR .zshrc WITH \"SOURCE <ABSOLUTE-PATH-OF-THIS-FILE>\"" > $aliasFile
echo "#GOTO_FILE=\"<ABSOLUTE-PATH-OF-THIS-FILE>\"" >> $aliasFile
echo "GOTO_FILE=\"$CONFIG_DIR/$GOTO_BIN\"" >> $aliasFile
echo "" >> $aliasFile

#Text without 3 first lines
tail -n$result ./alias.sh >> $aliasFile

#Delete the current alias.sh of the config file
rm ""$CONFIG_DIR"alias.sh"

#Change the name of alias.sh.new to alias.sh of the config file
mv $aliasFile ""$CONFIG_DIR"alias.sh"

#Some advises:
echo "This almost complete, please restart the terminal and check if all work correctly"
echo "If you want to add paths, use goto 1 to go the config dir and edit the config.json"
