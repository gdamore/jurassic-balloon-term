#!/bin/sh

set -eu

output_dir=${1:-dist/web}

case "$output_dir" in
  ""|/|.)
    echo "refusing unsafe web output directory: $output_dir" >&2
    exit 2
    ;;
esac

module_dir=$(go list -f '{{.Dir}}' github.com/gdamore/tcell/v3)
webfiles_dir="$module_dir/webfiles"
wasm_exec="$(go env GOROOT)/lib/wasm/wasm_exec.js"

mkdir -p "$output_dir"
cp -f web/index.html web/site.css "$output_dir/"
cp -f "$webfiles_dir/tcell.js" "$webfiles_dir/termstyle.css" \
  "$webfiles_dir/beep.wav" "$output_dir/"
mkdir -p "$output_dir/ghostty-web"
chmod -R u+w "$output_dir/ghostty-web"
cp -Rf "$webfiles_dir/ghostty-web/." "$output_dir/ghostty-web/"
cp -f "$wasm_exec" "$output_dir/"
cp -f LICENSE THIRD_PARTY_NOTICES.md "$output_dir/"

GOOS=js GOARCH=wasm CGO_ENABLED=0 go build \
  -trimpath -ldflags='-s -w' \
  -o "$output_dir/main.wasm" ./cmd/jurassic-balloon-term
