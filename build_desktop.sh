#!/bin/bash
set -e

echo "=== 编译桌面版前端 ==="
cd web && npm run build && cd ..

echo "=== 构建 Wails 桌面应用 ==="
wails build -tags desktop -o desktop.exe

echo "=== 桌面版构建完成 ==="
echo "产物路径: desktop.exe"
