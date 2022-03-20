##Colores##
ESC=$(printf '\033')
RESET="${ESC}[0m"
RED="${ESC}[31m"
GREEN="${ESC}[32m"

##Funciones con colores##
greenprint() { printf "${GREEN}%s${RESET}\n" "$1"; }
redprint() { printf "${RED}%s${RESET}\n" "$1"; }

error() {
    echo "$(redprint "$1")"
}

exito() {
    echo "$(greenprint "$1")"
}

echo "You have go installed?(y or any)"
read op

if [ "$op" = "y" ]; then

    GOTO_BIN="goto.bin"

    rm -f $GOTO_BIN

    #Build the code
    go build -o $GOTO_BIN ./*.go

    if [ $? -eq 0 ]; then
        exito "Compaling successfully"
    else
        error "Compaling failed"
        exit 1
    fi

else
    GOTO_BIN="bin/goto.bin"
fi

#if $XDG_CONFIG_HOME is empty
if [ -z "$XDG_CONFIG_HOME" ]; then
    XDG_CONFIG_HOME="$HOME/.config"
    echo "The variable \$XDG_CONFIG_HOME was created with \"\$HOME/.config\""
fi

#Define the config dir
CONFIG_DIR="$XDG_CONFIG_HOME"/goto/

#Create the config dir
mkdir -p $CONFIG_DIR

if [ $? -eq 0 ]; then
    exito "Config dir created successfully"
else
    error "Config dir couldn't be created"
    exit 1
fi

#Copy all of the repository to CONFIG_DIR
cp -r ./* $CONFIG_DIR

if [ $? -eq 0 ]; then
    exito "All files copied successfully"
else
    error "Files couldn't be copied"
    exit 1
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
        error "\"$SHELL_FILE\" doesn't exist" 
    fi
done 

#Give excute permission to the bin file
chmod +x ""$CONFIG_DIR"$GOTO_BIN"

if [ $? -eq 0 ]; then
    exito "Permission added successfully"
else
    error "The permission couldn't be added"
fi

#---------------------------------------#
#CHANGE CONFIG_FILE VARIABLE IN ALIAS.SH#
#---------------------------------------#

#Add the GOTO_FILE variable
aliasFile=""$CONFIG_DIR"alias.sh.new"

#Put the advise menssages and GOTO_FILE variable in the alias
echo "##ADD THIS FILE TO .bashrc OR .zshrc WITH \"SOURCE <ABSOLUTE-PATH-OF-THIS-FILE>\"" > $aliasFile
echo "#GOTO_FILE=\"<ABSOLUTE-PATH-OF-THIS-FILE>\"" >> $aliasFile
echo "GOTO_FILE=\""$CONFIG_DIR"$GOTO_BIN\"" >> $aliasFile
echo "" >> $aliasFile

#Number of lines of the text
num=$(wc -l alias.sh | cut -d " " -f 1)

#Remainder of 3 first lines of orginal text(the lines of GOTO_FILE variable)
result=$(expr $num - 3)

#Text without 3 first lines
tail -n$result ./alias.sh >> $aliasFile

#Delete the current alias.sh of the config file
rm ""$CONFIG_DIR"alias.sh"

#Change the name of alias.sh.new to alias.sh of the config file
mv $aliasFile ""$CONFIG_DIR"alias.sh"

#Some advises:
echo "This almost complete, please restart the terminal and check if all work correctly"