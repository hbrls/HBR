#!/bin/bash
# build-mac.sh — 薄壳脚本，转发到 Makefile
# 无参数：构建 macOS 版本（make build-mac）
# 有参数：透传给 make，例如 ./build-mac.sh clean
set -e
cd "$(dirname "$0")"

if [ $# -eq 0 ]; then
    exec make build-mac
else
    exec make "$@"
fi
