# Panduan OmniRoute

## Instalasi OmniRoute

```bash
git clone https://github.com/your-org/omniroute.git
cd omniroute
make build   # menghasilkan biner `omniroute`
./omniroute --config=config.yml
```

## Konfigurasi Contoh (`config.yml`)

```yaml
log_level: info
listen_addr: "0.0.0.0:3000"
providers:
  - name: openai
    base_url: "https://api.openai.com/v1"
    api_key_env: "OPENAI_API_KEY"
    weight: 1
  - name: anthropic
    base_url: "https://api.anthropic.com/v1"
    api_key_env: "ANTHROPIC_API_KEY"
    weight: 1
  - name: ollama
    base_url: "http://127.0.0.1:11434/v1"
    weight: 1
  - name: custom_omni
    base_url: "http://omniroute-server:3000/v1"
    weight: 2
```

## Menggunakan OmniRoute via Tailscale

1. **Aktifkan Tailscale** di host yang menjalankan OmniRoute (`tailscale up`).
2. **Aktifkan MagicDNS** di dashboard web Tailscale.
3. Di device Android, **instal aplikasi Tailscale** dan login ke Tailnet yang sama.
4. Aplikasi Android akan resolve `omniroute-server` melalui DNS Tailscale secara otomatis.

## Daftar Model (Discovery)

OmniRoute mendukung endpoint standar OpenAI `/v1/models`. Aplikasi Android Anda dapat:

* Fetch sekali saat provider disimpan.
* Cache lokal menggunakan Room.
* Tampilkan di dropdown model (`model_id`, `model_name`, `owned_by`).

Tombol "Refresh Models" memicu permintaan GET ulang.