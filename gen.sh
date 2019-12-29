#!/bin/sh

for typ in Int Int64 String Uint Uint64 Uintptr; do
	typlc=`echo "$typ" | tr A-Z a-z`

	echo "Generating $typlc.go"
	cat "type_go.txt" | sed -e "s/TypeName/$typ/g;s/typename/$typlc/g" >"$typlc.go"
done
