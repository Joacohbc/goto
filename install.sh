rm -f goto
go build -o goto ./main.go ./config.go 

CONFIG_DIR="$XDG_CONFIG_HOME"/goto/
mkdir -p $CONFIG_DIR

cp ./* $CONFIG_DIR

SHELL_FILE="<PUT-PATH>"

echo "source `pwd`/alias.sh" >> $SHELL_FILE
echo "Finish c:"