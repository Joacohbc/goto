##ADD THIS FILE TO .bashrc OR .zshrc WITH "SOURCE <ABSOLUTE-PATH-OF-THIS-FILE>"   
# GOTO_FILE="<ABSOLUTE-PATH-OF-THIS-FILE>"
GOTO_FILE="$XDG_CONFIG_HOME/goto/goto.bin"

gotop() {
    args=`echo $@`
    DESTINATION=$("$GOTO_FILE" --path="$args")
    echo $DESTINATION 
}

gotoc() {
    args=`echo $@`
    DESTINATION=$("$GOTO_FILE" --path="$args")
    echo "\"$DESTINATION\""
}

goto-add() {
    echo `"$GOTO_FILE" --add="$1,$2"`
}

goto() {
    args=`echo $@`
    
    DESTINATION=$("$GOTO_FILE" --path="$args")

    #If the return isn't an error
    if [[ $DESTINATION != *"Error:"* ]]; then

        cd "$DESTINATION"   

        echo "Go to:" $DESTINATION 
        
    else 
        #If it is an error, print it
        echo "$DESTINATION"
    fi
    
}

#Only return the path for the directory
alias gotop="gotop"

#Only return the path for the directory with ""
alias gotoc="gotoc"

#Move to the destination directory
alias goto='goto'

#Add new Path
alias goto-add='goto-add'

#Show help message
alias gotoh="\"$GOTO_FILE\" --help"

#Show version information
alias gotov="\"$GOTO_FILE\" --version"