#!/bin/bash

# Frontend Embed Script
# 将前端构建产物复制到 embed 目录

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR"
FRONTEND_DIR="$SCRIPT_DIR/../frontend"
EMBED_DIR="$BACKEND_DIR/embed/dist"

echo "📦 准备嵌入前端资源..."

if [ ! -d "$FRONTEND_DIR/dist" ]; then
    echo "❌ 前端未构建，请先运行: cd $FRONTEND_DIR && npm run build"
    exit 1
fi

rm -rf "$EMBED_DIR"
mkdir -p "$EMBED_DIR"

echo "📁 复制前端构建产物到 embed 目录..."
cp -r "$FRONTEND_DIR/dist/"* "$EMBED_DIR/"

echo "✅ 前端资源已嵌入到 $EMBED_DIR"
echo "📊 文件统计:"
echo "   - 文件数: $(find "$EMBED_DIR" -type f | wc -l)"
echo "   - 总大小: $(du -sh "$EMBED_DIR" | cut -f1)"
