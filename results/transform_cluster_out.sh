#!/bin/bash

if [ $# -ne 1 ]; then
	echo "Incorrect number of arguments: $# (expected 1)"
	exit 86; 
fi

if [ ! -f $1 ]; then
	echo "File \"$1\" does not exist"
	exit 33;
fi

name=$( echo $1 | cut -d '.' -f 1 )
egrep '^([0-9]+.*)|(Ellapsed.*)' "$1" | sed 's/Ellapsed Time Since Last:  //' | tr '\n' ',' | sed -r 's/,([0-9]+_)/\n\1/g' | tr '_' ',' > "$name.csv"

