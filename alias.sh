##ADD THIS FILE TO .bashrc OR .zshrc WITH "SOURCE <PATH-OF-THIS-FILE>"   
GOTO_FILE="<ABSOLUTE-PATH-OF-THIS-FILE>"

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

alias goto='goto'
alias gotoh="\"$GOTO_FILE\" --help"
alias gotov="\"$GOTO_FILE\" --version"