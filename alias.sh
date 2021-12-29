##ADD THIS FILE TO .bashrc OR .zshrc WITH "SOURCE <PATH-OF-THIS-FILE>"   
GOTO_FILE="/home/joaco/Archivos/Colegio y Estudio/Z-Proyectos/Go/goto/goto.test"

goto() {
    args=`echo $@`
    
    DESTINATION=$("$GOTO_FILE" $args)

    #If the return isn't an error
    if [[ "$DESTINATION" != *"Error:"* ]] || [[ "$DESTINATION" != *"flag provided but not defined:"* ]]; then

        cd "$DESTINATION"   

        echo "Go to:" $DESTINATION 
        
    else 
        #If it is an error, print it
        echo "$DESTINATION"
    fi  
}