#!/bin/bash
set -e

echo "🔧 Membundling proot binary..."

# Salin proot dan proot-distro ke direktori jniLibs
mkdir -p app/src/main/jniLibs/armeabi-v7a
mkdir -p app/src/main/jniLibs/arm64-v8a
mkdir -p app/src/main/jniLibs/x86
mkdir -p app/src/main/jniLibs/x86_64

# Asumsi assets/proot dan assets/proot-distro sudah ada di repo
cp -r assets/proot/* app/src/main/jniLibs/
cp -r assets/proot-distro/* app/src/main/jniLibs/

# Pastikan binary memiliki izin eksekusi
find app/src/main/jniLibs -type f -name "proot*" -exec chmod +x {} \;

echo "✅ Proot binary berhasil dibundling!"