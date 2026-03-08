#!/bin/bash

echo "开始编译下载器ui"
cd web && npm run build && cd ..
mkdir -p build
echo "开始编译下载器"
go build -o build/go_video.exe
cp ffmpeg.exe build/

echo "编译证书注册器"
go build -o build/install_cert.exe cmd/proxy/main.go

echo "运行 go_video.exe 启动服务，访问 http://localhost:8080" > build/README.txt
echo "完成"