rm -f goto

#Build the code
go build -o goto ./main.go ./config.go 

#Define the config dir
CONFIG_DIR="$XDG_CONFIG_HOME"/goto/

#Create the config dir
mkdir -p $CONFIG_DIR

#Copy all of the repository to CONFIG_DIR
cp ./* $CONFIG_DIR

#Absolute path line /home/username/.bashrc, not ~/.bashrc
SHELL_FILE="<PUT-PATH>"

#Add the alias.sh to $SHELL_FILE
echo "" >> $SHELL_FILE #New line
echo "Aliases to use goto:" >> $SHELL_FILE
echo "" >> $SHELL_FILE #New line
echo "source $CONFIG_DIR/alias.sh" >> $SHELL_FILE

#Some advises:
echo "This almost complete, please change GOTO_FILE variable in $CONFIG_DIR/alias.sh to complete, IF THE CURRENT GOTO_FILE DON'T WORK!"
echo "If you want to add paths, use goto 1 to go the config dir and edit the config.json"