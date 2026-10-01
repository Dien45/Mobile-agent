# Panduan Deployment

## Alur Kerja CI/CD

Lihat `.github/workflows/android-build.yml` di root repo.

### Workflow Utama: `android-build.yml`

- **Trigger**: push ke `main` atau pull request
- **Jobs**:
  1. `build`: Build APK Android
  2. `omniroute-build`: (Opsional) Build binary OmniRoute

### Langkah-langkah dalam Job `build`:

1. Checkout kode
2. Setup JDK 21 (temurin)
3. Jalankan `scripts/bootstrap_proot.sh` (membundling proot binary)
4. Build APK dengan `./gradlew app:assembleDebug --stacktrace`
5. Jalankan tes unit (opsional): `./gradlew app:testDebugUnitTest`
6. Upload APK sebagai artifact
7. Jika push ke `main`: buat GitHub Release dengan `softprops/action-gh-release`

### Job `omniroute-build` (Opsional):

- Checkout submodule `omniroute`
- Setup Go 1.22
- Jalankan `make build` di direktori `omniroute`
- Upload binary sebagai artifact

## Distribusi

### APK Langsung

README menyertakan badge "Latest Release" yang mengarah ke release GitHub terbaru.

### Link Download Manual

```bash
curl -L -o agentik.apk $(curl -s https://api.github.com/repos/your-org/android-agentik/releases/latest | jq -r '.assets[] | select(.name | endswith(".apk")) .browser_download_url')
```

## Persyaratan Build Lokal

```bash
# Install JDK 21
sudo apt-get install openjdk-21-jdk

# Install Gradle wrapper
./gradlew wrapper --gradle-version 8.5

# Build APK (termasuk proot binary)
./gradlew app:assembleDebug
```

## Persiapan Release

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
# GitHub Actions akan membuat release secara otomatis.
```

## Tips Pemecahan Masalah

- **Gagal build proot binary**: Pastikan `scripts/bootstrap_proot.sh` memiliki izin eksekusi dan assets ada.
- **Gagal resolve hostname Tailscale**: Verifikasi `Network Access` diaktifkan di aplikasi Tailscale Android.
- **API key tidak terbaca**: Pastikan secret di GitHub Actions sesuai nama (`OPENAI_API_KEY`, dsb.).