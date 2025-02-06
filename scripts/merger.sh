#!/bin/bash

if [ $# -lt 2 ]; then
    echo "Usage: $0 <directory> <output_file>"
    exit 1
fi

DIR="$1"
OUTPUT_FILE="$2"

ls "$DIR" | sort | while read FILE; do
    if [ -f "$DIR/$FILE" ]; then
        cat "$DIR/$FILE" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    fi
done

echo "Merged files into $OUTPUT_FILE"
