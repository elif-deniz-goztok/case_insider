#!/bin/sh
# Install this as a git pre-commit hook:
#   cp .claude/hooks/pre-commit.sh .git/hooks/pre-commit
#   chmod +x .git/hooks/pre-commit

set -e

echo "→ gofmt check..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
  echo "✗ Unformatted files (run gofmt -w .):"
  echo "$UNFORMATTED"
  exit 1
fi

echo "→ go vet..."
go vet ./...

echo "→ go build..."
go build ./...

echo "→ go test..."
go test ./...

echo "✓ All checks passed"
