#!/bin/bash

# Define variables
GOFILE="main.go"
OUTPUT="nftfetch"
BIN_DIR="$HOME/bin"

# Ensure the bin directory exists
if [ ! -d "$BIN_DIR" ]; then
    mkdir -p "$BIN_DIR"
fi

# Build the Go file
go build -o "$OUTPUT" "$GOFILE"

# Check if the build was successful
if [ $? -ne 0 ]; then
    echo "Build failed"
    exit 1
fi

# Move the executable to the bin directory and overwrite if it exists
mv -f "$OUTPUT" "$BIN_DIR"

# Check if the move was successful
if [ $? -ne 0 ]; then
    echo "Move failed"
    exit 1
fi

echo "Build and move completed successfully"
