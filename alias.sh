##ADD THIS FILE TO .bashrc OR .zshrc WITH "SOURCE <PATH-OF-THIS-FILE>"   
GOTO_FILE="$XDG_CONFIG_HOME/goto/goto64.bin"

#GOTO FUNC
goto() {
    args=`echo $@`
    
    DESTINATION=$("$GOTO_FILE" $args)

    #If the return isn't an error, put a bad argument, or the args have argument don't use the cd
    if [[ "$DESTINATION" != *"Error:"* ]] && [[ "$DESTINATION" != *"flag provided but not defined:"* ]] && [[ "$args" != "-"* ]]; then

        cd "$DESTINATION"   

        echo "Go to:" $DESTINATION 
        
    else 
        #If it is an error, print it
        echo "$DESTINATION"
    fi  
}

#cd is change by goto function
alias cd="goto"