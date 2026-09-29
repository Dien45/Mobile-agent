# Agentik Android – Alpine Linux + Hermes AI

**Catatan**: Aplikasi ini dibangun di luar Play Store. Instal langsung dari rilis terbaru.

## 🚀 Instalasi

1. **Clone repo** (atau unduh APK):
   ```bash
   curl -L -o agentik.apk $(curl -s https://api.github.com/repos/your-org/android-agentik/releases/latest | jq -r '.assets[] | select(.name | endswith(".apk")) .browser_download_url')
   ```
2. **Izinkan instalasi dari sumber yang tidak dikenal** (Pengaturan → Keamanan).
3. **Jalankan** – wizard akan mengunduh Alpine Linux melalui proot-distro.

## 📋 Fitur Utama

* **Shell Alpine Linux** langsung di ponsel (tanpa root).
* **Hermes Agent** dengan provider LLM yang dapat diganti (OpenAI, Groq, Anthropic, Ollama, **OmniRoute**).
* **File manager visual** untuk Alpine FS.
* **Thema biru pastel** yang ramah pengguna.

## 🔧 Konfigurasi Provider (OmniRoute + Tailscale)

1. Buka **Settings → Provider → Tambah Custom**.
2. **Nama**: `OmniRoute (Tailscale)`
3. **Base URL**: `http://omniroute-server:3000/v1`  *(MagicDNS OmniRoute Anda)*
4. **API Key**: token API dari server OmniRoute (atau biarkan kosong jika tidak butuh auth).
5. Simpan → aplikasi akan resolve `omniroute-server` melalui Tailscale dan mengambil daftar model.

## 📄 Dokumentasi

* `docs/deployment.md` – panduan CI/CD lengkap, skrip, dan tips pemecahan masalah.
* `docs/omniroute.md` – panduan instalasi OmniRoute, konfigurasi, dan integrasi Tailscale.

## 🐞 Masalah & Bantuan

Buka Issues di repo ini. Sertakan:
* Tangkapan layar log error.
* Versi Android / build APK.
* Tangkapan layar layar pengaturan provider.

---

## 🎨 Kredit Estetika

* **Material 3** (Surface, Color System)
* **Warna Biru Pastel** – Primary `#B3C7F7`, Surface `#F5F7FF`
* **JetBrains Mono** (code font) + **Noto Sans** (UI)

---

## 📜 Lisensi

MIT License – Hakcipta (c) 2024‑2025 your-org