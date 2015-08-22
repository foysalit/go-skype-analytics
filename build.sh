#!/bin/sh
go build

if [ $? -eq 0 ]; then
	./sky $1 $2
else
	echo "Error!"
	exit
fi