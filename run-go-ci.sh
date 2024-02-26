#!/usr/bin/env bash
set -euo pipefail

readarray -d '' ALL_GO_MODS < <(find . -type f -depth 2 -name "go.mod" -print0)

# Run go vet
for mod in "${ALL_GO_MODS[@]}"; do
  pushd "$(dirname "$mod")"
  go vet ./...
  go test ./...
  go test -race ./...
  popd
done
