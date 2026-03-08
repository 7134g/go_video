#!/bin/bash

cd web && npm run build && cd ..
mkdir -p build
go build -o build/go_video.exe
cp ffmpeg.exe build/
echo "运行 go_video.exe 启动服务，访问 http://localhost:8080" > build/README.txt
