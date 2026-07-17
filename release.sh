#!/bin/bash
set -e

echo "开始编译下载器UI"
cd web && npm run build && cd ..
mkdir -p build

case "$(uname -s)" in
  MINGW*|MSYS*|CYGWIN*)
    EXT=".exe"; FFMPEG="ffmpeg.exe"
    FFMPEG_URL="https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" ;;
  Darwin)
    EXT=""; FFMPEG="ffmpeg"
    FFMPEG_URL="https://evermeet.cx/ffmpeg/getrelease/zip" ;;
  Linux)
    EXT=""; FFMPEG="ffmpeg"
    case "$(uname -m)" in
      aarch64|arm64) FFMPEG_URL="https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz" ;;
      *)             FFMPEG_URL="https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz" ;;
    esac ;;
  *) echo "不支持的平台: $(uname -s)"; exit 1 ;;
esac

# 缺少 ffmpeg 时按平台从网上拉取,解压后放到 build/ 与项目根目录(项目根做缓存,下次构建直接命中)。
# ffmpeg 仅作合并兜底(默认已用纯 Go remux),下载失败不阻断构建。
download_ffmpeg() {
  echo "未找到 $FFMPEG,尝试从网上下载: $FFMPEG_URL"
  local tmp="build/.ffmpeg_tmp"
  rm -rf "$tmp"; mkdir -p "$tmp"
  local archive="$tmp/archive"
  if ! curl -fL --progress-bar "$FFMPEG_URL" -o "$archive"; then
    echo "下载失败,请手动从 $FFMPEG_URL 下载 $FFMPEG 放到项目根目录后重试"
    rm -rf "$tmp"; return 1
  fi
  case "$FFMPEG_URL" in
    *.zip|*/zip) unzip -q "$archive" -d "$tmp" ;;
    *.tar.xz)    tar -xf "$archive" -C "$tmp" ;;
  esac
  local found
  found="$(find "$tmp" -type f -name "$FFMPEG" | head -n1)"
  if [ -z "$found" ]; then
    echo "解压后未找到 $FFMPEG,请手动从 $FFMPEG_URL 处理"
    rm -rf "$tmp"; return 1
  fi
  chmod +x "$found"
  cp "$found" "./$FFMPEG"   # 缓存到项目根目录
  cp "$found" build/
  rm -rf "$tmp"
  echo "已下载 ffmpeg -> build/$FFMPEG"
}

echo "开始编译 Wails 桌面应用"
wails build -tags desktop -o "desktop${EXT}"
mv "build/bin/desktop${EXT}" "build/desktop${EXT}"
rmdir "build/bin"

if [ -f "$FFMPEG" ]; then
  cp "$FFMPEG" build/
else
  download_ffmpeg || echo "提示: ffmpeg 为可选(仅作合并兜底,默认已用纯 Go remux),可忽略本次失败"
fi

echo "编译证书注册器"
go build -o "build/install_cert${EXT}" ./tool/install-cert

echo "编译证书卸载器"
go build -o "build/uninstall_cert${EXT}" ./tool/uninstall-cert

echo "拷贝 Chrome 扩展"
rm -rf build/chrome_ext
cp -r chrome_ext build/chrome_ext

echo "拷贝 README.md"
cp README.md build/README.md

echo "完成"
