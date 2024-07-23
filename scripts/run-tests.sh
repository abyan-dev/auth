#!/bin/bash

output_dir="reports"
output_file="${output_dir}/coverage.out"

mkdir -p "$output_dir"

echo "mode: set" > "$output_file"

for dir in $(find . -type d ! -path "./pkg/server*" ! -path "./cmd/api*" ! -path "./pkg/response*" ! -path "./pkg/middleware*"); do
    if ls "$dir"/*_test.go &> /dev/null; then
        go test -coverprofile=tmp_coverage.out "$dir"
        if [ -f tmp_coverage.out ]; then
            tail -n +2 tmp_coverage.out >> "$output_file"
            rm tmp_coverage.out
        fi
    fi
done

go tool cover -html="$output_file" -o "${output_dir}/coverage.html"
