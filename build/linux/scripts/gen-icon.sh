#!/bin/sh
# Generates build/linux/icons/getdevice.png (512x512) from the Wails app icon
# at build/appicon.png. Uses sips on macOS; ImageMagick elsewhere.
set -e

# Repo root = script dir (build/linux/scripts) up three levels
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
SRC="$ROOT/build/appicon.png"
OUT_DIR="$ROOT/build/linux/icons"
OUT="$OUT_DIR/getdevice.png"

mkdir -p "$OUT_DIR"

if [ ! -f "$SRC" ]; then
    echo "ERROR: $SRC not found" >&2
    exit 1
fi

if command -v magick >/dev/null 2>&1; then
    magick "$SRC" -resize 512x512 "$OUT"
elif command -v convert >/dev/null 2>&1; then
    convert "$SRC" -resize 512x512 "$OUT"
elif command -v sips >/dev/null 2>&1; then
    cp "$SRC" "$OUT"
    sips -z 512 512 "$OUT" >/dev/null
else
    # Fall back to copying the original icon; most toolkits scale it fine.
    echo "WARN: no ImageMagick or sips found, copying 1024px icon as-is" >&2
    cp "$SRC" "$OUT"
fi

echo "Generated $OUT"
