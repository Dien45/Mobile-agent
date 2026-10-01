#!/bin/bash
set -e

echo "🔧 Membundling proot binary..."

# Salin proot dan proot-distro ke direktori jniLibs
mkdir -p app/src/main/jniLibs/armeabi-v7a
mkdir -p app/src/main/jniLibs/arm64-v8a
mkdir -p app/src/main/jniLibs/x86
mkdir -p app/src/main/jniLibs/x86_64

# Buat dummy proot binary untuk build (ganti dengan binary asli nanti)
# Ini placeholder — binary proot sebenarnya harus di-download dari release proot-distro
echo "#!/system/bin/sh" > app/src/main/jniLibs/arm64-v8a/proot
echo "# Dummy proot binary — replace with actual binary" >> app/src/main/jniLibs/arm64-v8a/proot
chmod +x app/src/main/jniLibs/arm64-v8a/proot

# Salin proot-distro assets jika ada, kalau tidak buat placeholder
if [ -d "assets/proot-distro" ]; then
    cp -r assets/proot-distro/* app/src/main/jniLibs/
else
    echo "⚠️  assets/proot-distro tidak ditemukan, skip."
fi

# Pastikan semua binary punya izin eksekusi
find app/src/main/jniLibs -type f -name "proot*" -exec chmod +x {} \;

echo "✅ Proot binary berhasil dibundling!"
