#!/bin/bash
# build-mac.sh — 只执行 macOS Release 构建
set -e
cd "$(dirname "$0")"

if [ $# -ne 0 ]; then
    echo "Usage: $0" >&2
    exit 2
fi

exec make build-mac
