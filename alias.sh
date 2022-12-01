##ADD THIS FILE TO .bashrc OR .zshrc WITH "SOURCE <PATH-OF-THIS-FILE>"   
GOTO_FILE="$XDG_CONFIG_HOME/goto/bin/goto.bin"

#GOTO FUNC
goto() {
    OUTPUT=$("$GOTO_FILE" $@)
    
    #If the return "2", the program return a gpath successfully
    if [ $? -eq 2 ]; then
        cd "$OUTPUT"   
        echo "Go to:" $OUTPUT
    elif [ $? -eq 1 ]; then # If error exit with status 1
        echo "$OUTPUT"
        return 1
    else
        echo "$OUTPUT"
    fi
}

#cd is change by goto function
alias cd="goto"
alias cdt="goto -t"
