#!/bin/bash
set -e

echo "开始编译下载器ui"
cd web && npm run build && cd ..
mkdir -p build

case "$(uname -s)" in
  MINGW*|MSYS*|CYGWIN*) EXT=".exe"; FFMPEG="ffmpeg.exe" ;;
  Darwin)               EXT="";     FFMPEG="ffmpeg" ;;
  Linux)                EXT="";     FFMPEG="ffmpeg" ;;
  *) echo "不支持的平台: $(uname -s)"; exit 1 ;;
esac

echo "开始编译下载器"
go build -o "build/go_video${EXT}"

if [ -f "$FFMPEG" ]; then
  cp "$FFMPEG" build/
else
  echo "提示: 未找到 $FFMPEG,请通过包管理器安装 (brew/apt/choco install ffmpeg)"
fi

echo "编译证书注册器"
go build -o "build/install_cert${EXT}" ./cmd/proxy

echo "拷贝 Chrome 扩展"
rm -rf build/chrome_ext
cp -r chrome_ext build/chrome_ext

echo "运行 go_video${EXT} 启动服务,访问 http://localhost:8080" > build/README.txt
echo "完成"
