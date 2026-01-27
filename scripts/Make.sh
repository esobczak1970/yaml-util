#!/bin/bash
# /scripts/Make.sh

# Ensure dependencies are installed
go get gopkg.in/yaml.v3
go mod tidy

# Create missing files if they do not exist
for file in "maxify/maxify.go" "maxify/maxify_test.go" "verbose/verbose.go" "verbose/verbose_test.go"; do
    if [ ! -f "$file" ]; then
        echo "Creating placeholder: $file"
        mkdir -p "$(dirname "$file")"
        touch "$file"
        echo "package $(basename "$(dirname "$file")")" > "$file"
        echo "" >> "$file"
    fi
done

echo "Setup complete!"
