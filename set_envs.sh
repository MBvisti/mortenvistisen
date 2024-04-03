#!/bin/bash

# Specify the path to the .env.staging file
env_file=".env.prod"

# Check if the file exists
if [ ! -f "$env_file" ]; then
    echo "Error: $env_file not found."
    exit 1
fi

# Loop over each line in the file
while IFS= read -r line; do
    # Print each line (replace this with your desired action)
    echo "Processing line: $line"
    read -r keyName paramValue <<< "$line"

    echo "Split fields:"
    echo "Field 1: $keyName"
    echo "Field 3: $paramValue"

    flyctl secrets set $keyName=$paramValue
done < "$env_file"

