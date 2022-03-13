##ADD THIS FILE TO .bashrc OR .zshrc WITH "SOURCE <PATH-OF-THIS-FILE>"   
GOTO_FILE="$XDG_CONFIG_HOME/goto/goto.bin"

#GOTO FUNC
goto() {
    OUTPUT=$("$GOTO_FILE" $@)

    #If the return "3", the program return a gpath successfully
    if [[ "$?" == "3" ]]; then

        cd "$OUTPUT"   

        echo "Go to:" $OUTPUT
        
    else 
        #If not a "3", can be 0 or 1
        echo "$OUTPUT"
    fi  
}

#cd is change by goto function
alias cd="goto"
