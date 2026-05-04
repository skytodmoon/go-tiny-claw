#!/bin/bash

# Function to print usage
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -s, --source PATTERN    Source file pattern (e.g., '*_test.go')"
    echo "  -t, --target PATTERN    Target file pattern (e.g., '*.go')"
    echo "  -r, --recursive         Recursively process subdirectories"
    echo "  -d, --dry-run           Dry run (show changes without applying)"
    echo "  -h, --help              Show this help message"
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--source)
            SOURCE_PATTERN="$2"
            shift 2
            ;;
        -t|--target)
            TARGET_PATTERN="$2"
            shift 2
            ;;
        -r|--recursive)
            RECURSIVE=true
            shift
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Check required arguments
if [[ -z "$SOURCE_PATTERN" || -z "$TARGET_PATTERN" ]]; then
    echo "Error: Source and target patterns are required."
    usage
fi

# Find files
if [[ "$RECURSIVE" = true ]]; then
    find_cmd="find . -name \"$SOURCE_PATTERN\""
else
    find_cmd="find . -maxdepth 1 -name \"$SOURCE_PATTERN\""
fi

# Process files
echo "Batch renaming files:"
eval "$find_cmd" | while read -r file; do
    dir=$(dirname "$file")
    base=$(basename "$file")
    new_base=$(echo "$base" | sed "s/$SOURCE_PATTERN/$TARGET_PATTERN/")
    new_file="$dir/$new_base"

    if [[ "$DRY_RUN" = true ]]; then
        echo "Would rename: $file -> $new_file"
    else
        echo "Renaming: $file -> $new_file"
        mv "$file" "$new_file"
    fi
done