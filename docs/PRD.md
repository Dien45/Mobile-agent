# Product Requirements Document — Arka Agent

| Atribut | Nilai |
|---|---|
| Status | Draft v0.6 untuk review |
| Tanggal | 4 Oktober 2026 |
| Nama produk | **Arka** |
| Tipe produk | Personal AI agent, local-first, web-first |
| Implementasi utama | Go |
| Lisensi | Belum diputuskan |

## 1. Ringkasan

Arka adalah personal AI agent yang berjalan sebagai runtime lokal ringan, dapat memakai model LLM dari berbagai provider, mengeksekusi tools di komputer pengguna, mengingat informasi lintas sesi, mempelajari prosedur sebagai skill, dan diakses melalui antarmuka web lokal. Browser hanya menjadi presentation layer; agent runtime, storage, credentials, memory, skills, dan seluruh eksekusi tool tetap berada di mesin pengguna.

Inspirasi perilakunya adalah pola agent Hermes: loop penggunaan tool, persistent memory, reusable skills, pencarian sesi lama, delegasi subagent, dan automasi terjadwal. Arka sengaja tidak menyediakan integrasi WhatsApp, Telegram, Discord, atau kanal pesan lain. Produk berfokus pada web UI yang dilayani oleh daemon lokal. Arka harus merupakan implementasi mandiri dari nol—bukan salinan kode, nama, prompt, atau identitas Hermes.

Produk tahap pertama berfokus pada inti yang kecil tetapi benar: **local daemon + web UI + agent loop + tools + persistence + safety**. Setelah repository di-clone, satu installer harus menyiapkan dependency, membangun aplikasi, menginisialisasi direktori data, dan menyediakan command untuk menjalankan produk tanpa setup development manual.

## 2. Masalah

Chatbot biasa memiliki beberapa kelemahan:

1. Tidak dapat melakukan tindakan nyata secara konsisten.
2. Kehilangan konteks saat sesi berakhir.
3. Mengulang cara kerja yang sama karena tidak menyimpan prosedur yang berhasil.
4. Terikat pada satu provider/model.
5. Tidak transparan ketika menjalankan perintah berisiko.
6. Banyak agent web mengirim storage atau eksekusi ke cloud, sementara pengguna membutuhkan UI web dengan kontrol lokal penuh.
7. Instalasi agent sering mengharuskan pengguna merangkai runtime, frontend, database, dan dependency secara manual.

Arka mengatasi masalah tersebut dengan runtime agent yang persisten, provider-agnostic, tool-driven, dapat diaudit, dan aman secara default.

## 3. Visi produk

> Satu agent pribadi dengan kenyamanan aplikasi web dan kontrol penuh aplikasi lokal: ringan, mudah dipasang, tumbuh bersama pengguna, dan bertindak dengan aman.

### Prinsip produk

1. **Narrow core, extensible edges** — inti hanya berisi orkestrasi universal; kemampuan khusus hidup sebagai tool, skill, atau plugin.
2. **Local-only execution** — sesi, konfigurasi, memory, skills, artifacts, audit log, dan tool execution berada di daemon lokal; browser tidak menjadi tempat eksekusi agent. Tidak ada cloud sync bawaan.
3. **Human control** — aksi destruktif, privilege escalation, dan akses sensitif membutuhkan kebijakan atau persetujuan.
4. **Observable** — pengguna dapat melihat tool yang dipanggil, argumen, hasil, durasi, dan perubahan file.
5. **Provider-agnostic** — mendukung API OpenAI-compatible terlebih dahulu, lalu adapter native.
6. **Progressive disclosure** — skill dan memory dimuat sesuai kebutuhan agar context window hemat.
7. **Useful learning, not uncontrolled self-modification** — agent boleh membuat draft memory/skill, tetapi tidak mengubah binary atau security policy sendiri.
8. **No hidden-reasoning dependency** — produk tidak meminta atau menyimpan chain-of-thought privat model. Keputusan dijelaskan lewat plan singkat, tindakan, hasil, dan ringkasan alasan yang dapat diaudit.

## 4. Pengguna sasaran

### Persona utama: developer/power user

- Menginginkan kenyamanan web UI tetapi semua data dan tool tetap lokal.
- Dapat memakai terminal untuk instalasi awal, tetapi pekerjaan sehari-hari dilakukan dari browser.
- Ingin agent mengedit file, menjalankan command, mencari web, dan mengotomasi pekerjaan.
- Memerlukan kontrol atas model, biaya, data, dan izin.
- Menjalankan agent di laptop, VPS kecil, atau container.

### Persona sekunder: operator/pemilik bisnis kecil

- Ingin automasi terjadwal dari dashboard web lokal.
- Membutuhkan agent yang mengingat SOP dan preferensi.
- Tidak ingin mengelola stack berat.

## 5. Sasaran dan metrik

### Sasaran MVP

1. Setelah clone, satu command installer menyiapkan seluruh komponen dan pengguna dapat membuka web UI serta mulai chat dalam kurang dari 5 menit setelah konfigurasi provider.
2. Agent menyelesaikan task multi-step dengan loop tool tanpa restart proses.
3. Sesi dan memory bertahan setelah aplikasi ditutup.
4. Semua tool call tervalidasi, dibatasi timeout, dan tercatat di audit log.
5. Command berisiko meminta approval secara konsisten.
6. Model/provider dapat diganti tanpa mengubah agent core.
7. Menutup tab browser tidak menghentikan task; daemon lokal tetap menjadi pemilik lifecycle agent dan tool.
8. Tidak ada session, memory, credential, atau tool execution yang dipindahkan ke server web eksternal.

### Metrik keberhasilan

| Metrik | Target MVP |
|---|---:|
| Local daemon siap menerima request | < 1 detik, di luar migration pertama |
| Web UI first contentful load lokal | < 1 detik pada build production |
| Idle RSS target | < 80 MB |
| Keberhasilan resume sesi setelah restart | 100% pada integration test |
| Tool-call schema validation | 100% tool call |
| Audit coverage untuk tool execution | 100% tool call |
| Task success pada benchmark internal MVP | ≥ 80% |
| Persetujuan sebelum aksi berisiko | 100% test policy |
| Crash recovery tanpa korupsi DB | 100% test fault-injection utama |

Metrik latency LLM dan kualitas jawaban dilaporkan terpisah per provider/model.

## 6. Ruang lingkup

### 6.1 MVP (v0.1)

- Local daemon yang bind ke loopback (`127.0.0.1`/`::1`) secara default.
- Web UI responsif bertema **claymorphism** sebagai interface utama dengan streaming response, live tool events, dan web terminal.
- Chat, session browser dengan rename/delete, provider/model picker, memory viewer, skill viewer, tool activity, web terminal, approval dialog, settings, dan diagnostics dasar di web.
- CLI kecil untuk install/lifecycle/diagnostics (`start`, `stop`, `status`, `doctor`) dan automation noninteraktif; CLI chat bukan interface utama.
- Setup wizard melalui web serta config noninteraktif melalui command/environment variable.
- Frontend production assets di-embed ke binary Go agar runtime tidak membutuhkan Node.js atau web server terpisah.
- Installer repository untuk Unix-like (`./install.sh`: macOS, Linux, Android Termux, dan distro Linux dalam Android PRoot) serta Windows (`./install.ps1`) yang idempotent dan memverifikasi dependency/checksum.
- Pada mesin pengguna yang sudah dapat melakukan clone, installer hanya boleh mengasumsikan shell/PowerShell, Git, dan akses jaringan. Bila versi Go yang sesuai belum tersedia, installer mengunduh toolchain Go resmi yang pinned ke build cache lokal, memverifikasi checksum, lalu membangun Arka tanpa meminta pengguna memasang dependency secara manual.
- Platform target MVP: Windows, macOS, Linux desktop/server **x64/amd64 dan arm64**, Android native melalui Termux, dan distro Linux Android melalui PRoot. Installer mendeteksi platform/runtime dan memilih path, shell, lifecycle, browser-open command, serta capability yang kompatibel.
- `make install` sebagai jalur developer alternatif, bukan syarat bagi end user.
- Provider presets untuk OpenAI, Anthropic, Google Gemini, OpenRouter, xAI, Groq, Mistral, Azure OpenAI, dan Ollama/local.
- Provider **Custom OpenAI-compatible** dengan base URL, API key, custom headers, dan path override; jalur ini menjadi dukungan resmi untuk OmniRoute.
- Model scanner yang mengambil katalog model provider, menyimpan cache lokal, dan menyediakan manual model entry bila endpoint discovery tidak tersedia.
- Provider dan model dapat dipilih per sesi serta diganti dari header chat.
- Agent loop dengan tool calling, batas langkah, cancellation, retry terbatas, dan context budgeting.
- Tool registry dan policy engine.
- Tool inti:
  - baca file;
  - tulis file atomik;
  - patch file;
  - daftar/cari file;
  - terminal command;
  - pengelolaan proses;
  - memory;
  - session search;
  - todo/plan;
  - HTTP fetch/web extract sederhana;
  - Git/GitHub;
  - skill discovery/installation/management;
  - archive/download/artifact management;
  - structured code execution melalui process lokal;
  - local MCP client untuk tool eksternal yang dikonfigurasi.
- SQLite untuk sessions, messages, tool runs, approvals, dan full-text search.
- Memory ringkas berbasis file atau record terstruktur.
- Skills kompatibel format Agent Skills dengan UI untuk browse, install, update, enable/disable, inspect, rollback, dan uninstall.
- Sumber instalasi skill: katalog terkonfigurasi, Git URL/repository + subdirectory, local directory, dan local archive.
- Semua tool handler dan script skill dieksekusi oleh daemon/proses lokal; browser hanya mengirim intent dan merender event/result.
- Project context discovery (`AGENTS.md`, `.arka.md`, dan file konfigurasi yang disetujui).
- Structured logs, audit log, dan redaction secret.
- Approval mode: `ask`, `allow`, `deny` berdasarkan rule.
- Commands dasar: `arka start`, `arka stop`, `arka status`, `arka open`, `arka doctor`, `arka config`.
- API lokal versioned yang hanya dipakai web UI dan integrasi lokal.
- Task berjalan sebagai background job milik daemon, bukan milik halaman browser. Berpindah sesi, refresh, atau menutup tab tidak menghentikan agent.
- Web terminal lokal dengan PTY, multiple tabs, resize, reconnect, dan working directory per workspace.
- Login GitHub melalui OAuth Device Flow serta deteksi autentikasi GitHub CLI yang sudah ada.
- Git tools untuk status, diff, staging terarah, dan pembuatan commit lokal oleh agent. Push/PR/release tidak dilakukan otomatis pada MVP.

### 6.2 v0.2 — Learning loop

- Membuat dan memperbarui skill sebagai draft dari task kompleks yang berhasil.
- Memory curation: tambah, ganti, hapus, pin, dan deduplikasi.
- Nudge setelah sejumlah turn/tool call untuk mengevaluasi hal yang layak disimpan.
- Session summarization dan context compaction.
- MCP client untuk external tools.
- Import/export portable untuk session, memory, dan skill.
- Provider health analytics, configurable fallback chains, dan cost-aware routing setelah adapter MVP stabil.

### 6.3 v0.3 — Automation, delegation, dan web hardening

- Cron scheduler persisten.
- Subagent terisolasi dengan budget, toolset, dan working directory sendiri.
- Parallel delegation dengan concurrency limit.
- Delivery hasil task ke notification center di web UI.
- Optional desktop notification melalui browser dengan izin pengguna.
- Webhook/API lokal yang default-nya tetap loopback-only.

### 6.4 v1.0 — Local web agent matang

- Web UI production-ready dengan responsive desktop/mobile layout.
- Browser automation sebagai service/plugin lokal terisolasi.
- Profiles terpisah untuk work/personal/team.
- Plugin SDK stabil dan versioned.
- Container/SSH execution backend.
- Installer/updater lintas platform dengan rollback dan checksum.
- Packaging release satu binary per platform; web assets sudah embedded.

### Di luar scope awal

- Melatih foundation model sendiri.
- Menyamai seluruh fitur Hermes pada rilis pertama.
- Autonomous binary/source self-modification.
- Desktop/mobile GUI native.
- Integrasi WhatsApp, Telegram, Discord, Slack, atau messaging gateway lain.
- Hosting web UI sebagai SaaS atau menyimpan state agent di cloud.
- Menyimpan atau menampilkan private chain-of-thought model.
- Menjalankan shell tanpa policy, timeout, output cap, dan audit.
- Menjadi platform multi-tenant SaaS pada MVP.

## 7. Pengalaman pengguna utama

### 7.1 First run

1. Pengguna menjalankan `git clone ...`, masuk ke direktori, lalu menjalankan `./install.sh` (Unix) atau `./install.ps1` (Windows).
2. Installer mendeteksi OS/arsitektur, memverifikasi prerequisite, membangun atau memasang binary, menyiapkan direktori state, lalu menjalankan health check.
3. Pengguna menjalankan `arka start`; daemon lokal hidup dan browser otomatis membuka URL loopback.
4. Setup wizard web meminta pengguna memilih preset provider atau Custom OpenAI-compatible, memasukkan credential/endpoint yang diperlukan, lalu scan atau memilih model.
5. Credential disimpan lewat environment/OS keyring; tidak dikirim ke frontend kembali dan tidak ditulis plaintext ke config jika keyring tersedia.
6. Pengguna memilih workspace dan mode approval.
7. Halaman diagnostics menguji provider, DB, permission, dan tools.
8. Pengguna mulai chat dari web; refresh/menutup tab tidak menghapus sesi atau menghentikan daemon.

### 7.2 Task dengan tools

1. Pengguna memberi tujuan.
2. Agent membuat plan singkat bila task kompleks.
3. Model memilih tool melalui structured tool call.
4. Runtime memvalidasi schema dan policy.
5. Bila perlu, UI meminta approval dengan command dan dampak yang jelas.
6. Tool dieksekusi dengan timeout/cancellation.
7. Hasil dinormalisasi dan dikembalikan ke model.
8. Loop berlanjut sampai selesai, membutuhkan klarifikasi, atau mencapai batas.
9. Agent memberi hasil, perubahan yang dibuat, dan verifikasi—bukan transcript pemikiran tersembunyi.

### 7.3 Belajar dari pengalaman (v0.2)

1. Task kompleks selesai dan terverifikasi.
2. Evaluator menilai apakah prosedurnya reusable.
3. Agent membuat draft skill yang berisi trigger, prasyarat, langkah, verifikasi, dan failure modes.
4. Draft divalidasi dan, sesuai policy, diminta approval atau disimpan otomatis di area user.
5. Pada task serupa, hanya metadata skill yang dilihat lebih dahulu; isi lengkap dimuat ketika relevan.
6. Feedback penggunaan dicatat untuk usulan revisi skill berikutnya.

### 7.4 Mengelola dan melanjutkan sesi

- Pengguna melihat daftar sesi dari sidebar web, lengkap dengan judul, workspace, waktu update, status task, dan jumlah tool call.
- Judul otomatis dibuat dari percakapan awal dan dapat di-rename secara inline kapan saja.
- Sesi dapat dihapus melalui confirmation dialog. Delete awal bersifat soft delete agar kegagalan UI tidak langsung memusnahkan data; pengguna dapat memilih purge permanen dari trash.
- Sesi aktif tidak boleh dihapus ketika tool masih berjalan tanpa menghentikan task terlebih dahulu.
- Pengguna dapat resume berdasarkan ID.
- Sistem memuat transcript, summary, project context, dan memory relevan.
- Tool call yang terputus ditandai `interrupted`, tidak diam-diam dijalankan ulang.

### 7.5 Login GitHub dan membuat commit

1. Pengguna memilih **Connect GitHub** dari settings.
2. Arka menggunakan GitHub OAuth Device Flow: UI menampilkan user code dan membuka halaman verifikasi GitHub. Alternatifnya, Arka dapat memakai sesi `gh` yang sudah login setelah persetujuan pengguna.
3. Token disimpan di OS keyring dan tidak pernah dikirim kembali ke browser atau dimasukkan ke prompt model.
4. Agent memeriksa repository, `git status`, dan diff sebelum membuat commit.
5. Agent hanya men-stage file yang relevan dengan task, menjalankan verifikasi yang tersedia, lalu membuat commit dengan pesan ringkas dan attribution yang benar.
6. Commit lokal dapat diizinkan untuk seluruh trusted workspace. Push, force-push, merge, release, dan perubahan remote selalu merupakan capability terpisah dengan approval sendiri.

### 7.6 Memilih provider dan model

1. Pengguna memilih preset provider umum atau **Custom OpenAI-compatible**.
2. Untuk custom provider seperti OmniRoute, pengguna memasukkan base URL, API key reference, dan optional headers/path overrides.
3. Tombol **Scan models** memanggil endpoint katalog provider melalui daemon lokal, menormalisasi hasil, lalu menyimpan cache.
4. Jika scanning tidak didukung atau gagal, pengguna tetap dapat memasukkan model ID secara manual.
5. Header setiap sesi menampilkan provider dan model aktif. Pengguna dapat mengganti model tanpa membuat sesi baru.
6. Request yang sudah dikirim tetap diselesaikan oleh model lama. Perubahan model berlaku pada model call berikutnya dan dicatat di timeline sesi agar hasil dapat diaudit.

### 7.7 Berpindah sesi ketika agent bekerja

- Setiap turn aktif menjadi background task persisten pada daemon.
- Berpindah dari sesi A ke sesi B hanya melepas subscription UI sesi A; task sesi A tetap berjalan.
- Sidebar menampilkan status `running`, progress/tool count, `waiting approval`, `completed`, atau `failed` untuk tiap sesi.
- Pengguna dapat menjalankan task di beberapa sesi secara concurrent sampai batas concurrency; sisanya masuk queue.
- Kembali ke sesi A melakukan replay event dari cursor terakhir tanpa menggandakan model call atau tool execution.
- Hanya tombol Stop pada sesi terkait, deadline, policy failure, shutdown terkontrol, atau fatal error yang menghentikan task.

### 7.8 Web terminal

1. Pengguna membuka panel Terminal dan membuat tab baru pada workspace yang dipilih.
2. Daemon membuat local PTY dengan shell pengguna dan working directory workspace.
3. Input/output serta resize terminal dikirim melalui WebSocket yang terautentikasi.
4. Berpindah halaman tidak langsung mematikan terminal; sesi PTY bertahan sesuai idle timeout dan dapat di-reconnect.
5. Terminal manusia terpisah dari tool terminal agent. Agent tidak boleh membaca atau mengetik ke terminal manusia kecuali pengguna secara eksplisit memberikan capability tersebut.
6. Pengguna dapat menghentikan proses, menutup tab terminal, dan melihat status koneksi dari web.

### 7.9 Menginstal dan menggunakan skill

1. Pengguna membuka halaman **Skills** dan memilih katalog, memasukkan Git URL, memilih folder lokal, atau mengunggah archive lokal.
2. Daemon mengambil/membaca package, memvalidasi struktur dan metadata, menghitung checksum, lalu menampilkan source, version/commit, file, scripts, dependency, dan capability yang diminta.
3. Pengguna memilih scope instalasi: profile, global user, atau workspace tertentu.
4. Setelah approval, daemon memasang skill secara atomik ke local skill store dan menyimpan lock record untuk version/checksum/source.
5. Skill dapat di-enable/disable, di-update dengan diff, di-rollback ke versi sebelumnya, atau di-uninstall dari web.
6. Metadata skill tersedia untuk model, tetapi konten lengkap hanya dimuat saat relevan. Script skill tetap dijalankan melalui local tool runtime dan policy engine yang sama.
7. Agent boleh mencari atau menyarankan skill. Instalasi dari source baru tidak boleh dilakukan diam-diam; pengguna harus melihat provenance dan capability sebelum menyetujui.

### 7.10 Instalasi lintas platform

#### Windows

```powershell
git clone <repository>
cd <repository>
.\install.ps1
arka start
```

Installer memasang binary ke user-local path tanpa membutuhkan Administrator secara default. Daemon menggunakan Windows process primitives/Job Objects dan web terminal menggunakan ConPTY atau fallback shell yang terdeteksi.

#### macOS

```bash
git clone <repository>
cd <repository>
./install.sh
arka start
```

Mendukung Apple Silicon dan Intel. Installer user-local tidak membutuhkan `sudo`; browser dibuka melalui `open`.

#### Android — Termux

```bash
# Setelah Git tersedia di Termux
# (contoh bootstrap awal: pkg install git)
git clone <repository>
cd <repository>
./install.sh
arka start
```

Installer mengenali `$PREFIX`/Termux, menggunakan path app-private, dan membuka UI melalui `termux-open-url` bila tersedia atau menampilkan URL loopback. Akses shared storage hanya diminta bila pengguna memilih workspace di sana; Arka tidak menjalankan `termux-setup-storage` diam-diam.

#### Android — PRoot

```bash
# Di dalam distro proot (contoh Debian/Ubuntu)
git clone <repository>
cd <repository>
./install.sh
arka start --no-open
```

PRoot diperlakukan sebagai Linux tanpa systemd. Arka memakai foreground/background PID lifecycle sendiri, bind ke loopback, dan menampilkan URL untuk dibuka di browser Android. Fitur yang tidak tersedia karena kernel/proot harus dilaporkan sebagai `unsupported`, bukan dimock.

Semua platform menggunakan data lokal dan API loopback yang sama. Config/database portable secara schema, tetapi binary, shell integration, dan process backend spesifik platform.

## 8. Cara agent bekerja

### 8.1 State machine satu turn

```text
RECEIVE
  → BUILD_CONTEXT
  → MODEL_CALL
  → (FINAL | TOOL_REQUEST | CLARIFY | ERROR)
  → POLICY_CHECK
  → APPROVAL (opsional)
  → TOOL_EXECUTION
  → OBSERVE
  → MODEL_CALL ...
  → PERSIST
  → RESPOND
```

### 8.2 Context assembly

Urutan prioritas context:

1. System policy yang immutable selama turn.
2. Persona/profile pengguna.
3. Project context.
4. Daftar tool dan schema aktif.
5. Metadata skill relevan; full skill hanya on-demand.
6. Memory relevan dengan budget.
7. Summary sesi lama.
8. Pesan terbaru dan hasil tool.

Setiap bagian memiliki token budget. Runtime harus dapat memangkas output tool, merangkum history, dan mempertahankan pesan sistem serta objective terakhir.

### 8.3 Planning dan refleksi

- Plan bersifat opsional untuk task sederhana dan wajib untuk task yang diprediksi multi-step.
- Plan disimpan sebagai task checklist yang bisa diperbarui, bukan chain-of-thought.
- Setelah aksi, agent melakukan verifikasi berbasis bukti: exit code, test result, file diff, atau response status.
- Jika gagal, agent boleh memperbaiki strategi dalam batas retry dan budget.
- Setelah batas tercapai, agent berhenti secara aman, menjelaskan hambatan, dan meminta keputusan pengguna.

### 8.4 Loop guard

Konfigurasi default:

- Tidak ada batas total tool call per sesi. Satu sesi harus mampu menyimpan dan menjalankan lebih dari 500 tool calls.
- Maksimum 1.000 tool calls per turn/task secara default dan dapat dikonfigurasi. Batas tinggi ini mencegah pekerjaan panjang berhenti hanya karena hitungan tool, tetapi tetap menyediakan emergency ceiling.
- Runtime membuat checkpoint progress dan durable state secara berkala (default setiap 25 tool calls), sehingga task panjang dapat dipulihkan.
- History tool yang lama dikompaksi menjadi summary + artifact references tanpa menghapus audit record asli.
- Maksimum 3 kegagalan berulang dengan signature sama.
- Timeout per tool dan deadline per turn.
- Batas bytes output tool.
- Batas biaya/token per turn dan per hari bila provider menyediakan usage.
- Cancellation melalui tombol **Stop** di web atau `Ctrl+C` pada CLI yang mengalir ke model stream dan child process.
- Approval ditentukan oleh risiko/capability, bukan jumlah tool call. Mencapai 100, 500, atau jumlah tertentu tidak memunculkan approval baru dengan sendirinya.

## 9. Functional requirements

### FR-1 — Model provider dan model discovery

- Interface provider mencakup streaming message, structured tool call, usage, cancellation, error classification, dan model discovery opsional.
- Preset awal: OpenAI, Anthropic, Google Gemini, OpenRouter, xAI, Groq, Mistral, Azure OpenAI, dan Ollama/local. OpenAI memakai adapter API OpenAI; Anthropic dan Gemini memakai adapter native; provider yang menawarkan kontrak OpenAI-compatible memakai adapter kompatibilitas dengan preset endpoint/auth masing-masing. Agent core tetap provider-neutral.
- **Custom OpenAI-compatible** menerima nama profile, base URL, API key reference, custom headers, chat path, models path, dan compatibility flags. OmniRoute didukung melalui profile ini tanpa hard-code vendor-specific logic.
- Credential provider hanya diproses daemon dan disimpan di environment/OS keyring. API key tidak pernah dikirim kembali ke web UI setelah disimpan.
- Model scanner mencoba endpoint provider-native atau OpenAI-compatible `GET /v1/models`, menangani pagination bila ada, lalu menormalisasi `id`, display name, owner, context limit, modalities, tool support, dan availability sejauh metadata tersedia.
- Karena metadata provider sering tidak lengkap, capability yang tidak diketahui ditandai `unknown`, bukan ditebak. Pengguna dapat memasukkan model ID manual dan melakukan test connection.
- Katalog model memiliki cache lokal, timestamp, source, refresh button, loading/error state, dan fallback ke cache terakhir ketika provider sedang offline.
- Provider dan model dipilih per sesi. Setiap turn/message menyimpan provider profile ID dan model ID aktual agar pergantian model dapat diaudit.
- Pergantian model tidak membatalkan request yang sedang berjalan. Pilihan baru berlaku pada model call berikutnya; UI memperingatkan bila satu task berpotensi memakai dua model.
- Provider config tidak boleh bocor ke prompt/log.
- Retry hanya untuk error transient dan harus memakai exponential backoff + jitter. Fallback model tidak dijalankan diam-diam; harus dikonfigurasi dan ditampilkan di timeline.

### FR-2 — Agent runtime

- Agent loop deterministik pada level state transition.
- Setiap turn memiliki ID dan trace ID.
- Lifecycle turn dimiliki background task manager pada daemon dan tidak terkait lifecycle route, SSE connection, atau tab browser.
- Satu sesi hanya memiliki satu mutating turn aktif agar urutan percakapan konsisten; beberapa sesi dapat berjalan concurrent dengan queue dan global concurrency limit yang dapat dikonfigurasi.
- Event disimpan dengan monotonic sequence/cursor agar UI dapat replay setelah pindah sesi atau reconnect.
- Runtime mendukung stop, resume session, max steps hingga minimal 1.000 tool calls per task, token budget, checkpoint, dan deadline.
- Parser tool call menolak nama/argumen yang tidak sesuai registry/schema.
- Runtime tidak mengeksekusi teks biasa sebagai tool call.

### FR-3 — Tool system

Setiap tool wajib mendeklarasikan:

- nama dan versi;
- deskripsi model-facing;
- JSON Schema input/output;
- risk level (`read`, `write`, `execute`, `network`, `destructive`, `privileged`);
- timeout default;
- capability yang diperlukan;
- apakah idempotent;
- handler dan error terstruktur;
- supported OS/architecture/runtime (`windows`, `darwin`, `linux`, `termux`, `proot`) dan alasan jika unavailable.

Tool registry mendukung enable/disable per profile. Hasil tool harus memiliki status, output terpotong, metadata durasi, dan artifact reference bila output besar.

Execution boundary wajib:

- Browser tidak berisi implementasi tool produksi dan tidak menjalankan shell, filesystem, Git, memory, skill script, atau provider call.
- Browser hanya memanggil local API, mengirim approval/input, dan merender lifecycle event. JavaScript mock hanya boleh dipakai pada test/story fixture dan harus gagal build production bila mock mode aktif.
- Tool dispatcher, schema validation, policy check, credential resolution, handler, timeout, process supervision, result normalization, persistence, dan audit seluruhnya berjalan di Go daemon lokal.
- Tool eksternal melalui MCP dijalankan/dihubungkan dari daemon lokal. Jika MCP server berada di network, koneksinya tetap berasal dari daemon berdasarkan konfigurasi dan policy, bukan langsung dari browser.
- Tool dapat melakukan akses network bila memang fungsinya membutuhkan, tetapi eksekusi dan credential tetap lokal. UI harus membedakan `local execution` dari `network access`.
- Tool call mendukung sequential/parallel execution yang eksplisit, cancellation, idempotency metadata, retry policy, progress event, artifact, dan structured error.

Baseline toolsets MVP:

| Toolset | Kemampuan minimum |
|---|---|
| `files` | list, glob/search, read ranges, atomic write, patch, mkdir, move/copy, metadata |
| `terminal` | exec tanpa shell bila memungkinkan, shell opt-in, process list/status/kill, timeout |
| `git` | status, diff, log, branch info, explicit staging, unstage path, local commit |
| `github` | auth/status dan repository metadata; remote mutations tetap capability terpisah |
| `web` | HTTP fetch/download, content extraction, headers/status, size/time limits |
| `search` | web search melalui provider yang dikonfigurasi; request berasal dari daemon lokal |
| `browser` | local managed Chromium/CDP: navigate, snapshot, tabs, click, input, screenshot, download |
| `memory` | add/replace/remove/search dengan source dan audit |
| `sessions` | search/read summary dan status task tanpa bypass isolation |
| `skills` | list/inspect/install/update/disable/rollback/uninstall |
| `artifacts` | store/read/list/export dengan content hash dan size limit |
| `planning` | todo create/update/complete dan progress reporting |
| `code` | menjalankan code/script melalui supervised local process, bukan browser eval |
| `mcp` | discovery dan invocation server yang dikonfigurasi, tunduk pada policy |

Browser automation, image, voice, dan tool khusus lain dapat menjadi plugin setelah baseline stabil, tetapi tetap mengikuti local daemon boundary.

### FR-4 — Terminal safety

- Eksekusi tanpa shell bila memungkinkan; shell mode harus eksplisit.
- Working directory dibatasi ke workspace secara default.
- Environment dibangun dari allowlist.
- Command dengan pola destructive/privileged meminta approval atau ditolak.
- Secret pada command/output di-redact dari log jika terdeteksi.
- Child process harus dibersihkan saat cancel/exit.
- Sandbox/container menjadi backend tambahan, bukan syarat MVP.

### FR-5 — Persistence dan session

- SQLite memakai WAL, foreign key, migration version, dan transaksi.
- Simpan session, message, content part, tool call, tool result, usage, approval, summary, dan status turn.
- Sesi memiliki `title`, `workspace_id`, `created_at`, `updated_at`, `deleted_at`, dan revision/version untuk mencegah lost update.
- Rename sesi dapat dilakukan inline dan tidak mengubah isi transcript.
- Delete sesi menggunakan soft delete + trash terlebih dahulu; purge permanen menghapus transcript dan artifact yang tidak direferensikan setelah confirmation eksplisit.
- Delete/rename harus tetap responsif pada sesi dengan lebih dari 500 tool calls dan tidak memuat seluruh transcript ke memory.
- Pencarian history memakai SQLite FTS5 dan tidak menampilkan sesi dalam trash kecuali diminta.
- Satu writer tidak boleh menyebabkan transcript parsial yang dianggap selesai.
- Export JSON harus tersedia sebelum v1.0 agar data tidak terkunci.

### FR-6 — Memory

Tiga lapisan:

1. **Working memory** — context aktif sesi, tidak dianggap fakta permanen.
2. **Declarative memory** — fakta stabil, preferensi, keputusan, constraint, dan environment knowledge.
3. **Procedural memory** — skills/SOP reusable.

Setiap declarative memory memiliki ID, tipe, content, scope, source, confidence, created/updated time, dan status. Write harus dapat diaudit. Conflict tidak boleh diam-diam menimpa fakta pinned; sistem membuat candidate replacement.

MVP boleh memakai lexical retrieval FTS5. Embedding/semantic retrieval baru ditambahkan bila benchmark menunjukkan kebutuhan.

### FR-7 — Skills

Format skill mengikuti struktur sederhana yang kompatibel secara konsep dengan praktik Agent Skills:

```text
skills/<skill-name>/
├── SKILL.md
├── scripts/       # opsional
├── references/    # opsional
└── templates/     # opsional
```

`SKILL.md` minimal memuat nama, deskripsi/trigger, prerequisites, langkah, verifikasi, dan safety notes. Skill pihak ketiga tidak otomatis dipercaya; script-nya tunduk pada policy tool yang sama.

Skill manager wajib mendukung:

- browse/search metadata dari katalog yang dikonfigurasi;
- install dari catalog entry, Git HTTPS/SSH URL dengan ref + optional subdirectory, local directory, atau `.zip`/`.tar.gz` lokal;
- scope `user`, `profile`, dan `workspace`, dengan precedence dan conflict resolution yang terdokumentasi;
- pre-install inspection: source URL/path, publisher bila tersedia, commit/ref, checksum, license, file list, executable scripts, dependency, requested capabilities, serta compatibility OS/architecture/runtime;
- atomic install melalui staging directory, validation, lalu rename; kegagalan tidak meninggalkan half-installed skill;
- version lockfile, update check, diff sebelum update, retained previous version, rollback, enable/disable, dan uninstall;
- validation nama/path untuk mencegah traversal, symlink escape, archive bomb, executable surprise, dan overwrite skill lain;
- progressive disclosure: index hanya memuat nama/deskripsi/trigger; full `SKILL.md` dan resources dibuka on-demand;
- dependency setup melalui installer terdeklarasi dan approval terpisah. Skill tidak boleh menjalankan arbitrary install hook saat sekadar di-download/di-inspect;
- audit event untuk install/update/rollback/uninstall dan setiap script skill yang dijalankan.

Agent dapat memanggil `skills.search`, `skills.inspect`, dan mengusulkan instalasi. `skills.install` dari source yang belum dipercaya membutuhkan approval pengguna; trusted catalog policy dapat dikonfigurasi untuk update non-breaking otomatis.

### FR-8 — Project context

- Cari context dari current directory menuju root workspace dengan precedence terdokumentasi.
- Context tidak boleh dapat menurunkan system security policy.
- Ukuran, jumlah file, dan symlink traversal dibatasi.
- Pengguna dapat melihat asal semua context yang aktif.

### FR-9 — Approval

Approval tidak diminta ulang pada setiap chat. Saat onboarding atau membuka workspace baru, pengguna memilih policy profile:

- **Safe** — read otomatis; write/execute/commit meminta approval.
- **Standard (default)** — read, write dalam workspace, command development yang diizinkan, dan local git commit dapat berjalan otomatis; aksi destruktif, privileged, di luar workspace, atau remote tetap meminta approval.
- **Autonomous** — semua aksi non-destructive di trusted workspace otomatis; destructive/privileged/remote tetap mengikuti rule eksplisit.

Grant dapat berlaku `once`, `session`, atau `workspace`. Grant workspace bertahan lintas chat sehingga task panjang tidak terganggu dialog yang sama. Rule dievaluasi berdasarkan capability, target, command family, dan workspace—bukan berdasarkan jumlah tool call.

Prompt approval menampilkan:

- tool dan risk level;
- aksi/command sebenarnya;
- workspace/host target;
- alasan singkat;
- pilihan `allow once`, `allow session`, `allow workspace`, `deny`;
- untuk rule yang aman, opsi membuat allow rule terbatas.

Approval persisten harus spesifik, terlihat di Settings, dan dapat dicabut. Pola wildcard luas diberi peringatan. Perubahan policy oleh agent tidak diizinkan; hanya pengguna yang dapat menaikkan privilege.

### FR-10 — Observability

- Console menampilkan event streaming tanpa membocorkan secret.
- Structured log menggunakan JSON dan level.
- Audit record append-oriented untuk aksi eksternal.
- `arka doctor` memeriksa config, DB migration, provider connectivity, tool availability, dan permission.
- Metrics opsional: latency, token, estimated cost, tool success rate, retries.

### FR-11 — Scheduler (v0.3)

- Job tersimpan persisten dengan timezone eksplisit.
- Concurrency policy: forbid/allow/replace.
- Missed-run policy terdokumentasi.
- Job menjalankan session terisolasi dengan toolset dan budget sendiri.
- Hasil memiliki delivery target dan audit trail.

### FR-12 — Delegation (v0.3)

- Parent mengirim objective, context subset, constraints, dan output contract.
- Subagent tidak mewarisi semua secret/tool secara otomatis.
- Concurrency dan total budget dibatasi.
- Parent menerima ringkasan serta artifact, bukan seluruh hidden context.
- Cancellation parent mengalir ke child.

### FR-13 — Local web interface

- Go daemon menyajikan static frontend assets dan API dari origin/port yang sama.
- Browser berkomunikasi dengan daemon melalui HTTP lokal dan SSE atau WebSocket untuk event streaming.
- Browser tidak mengakses SQLite, shell, filesystem, Git, skill store, provider credential, atau tool process secara langsung. Tidak ada fallback tool execution di Web Worker, browser extension, IndexedDB, atau mocked frontend service pada production.
- Semua tool call dan mutation melewati local API, authentication/session check, daemon tool registry, schema validation, policy engine, approval, execution, persistence, dan audit.
- Refresh/reconnect dan perpindahan sesi melakukan replay event dari cursor tanpa menggandakan model call atau tool execution.
- Menutup stream browser hanya mengubah subscriber count; tidak membatalkan background task atau PTY.
- UI menampilkan status koneksi, status daemon, task aktif, plan, tool call, approval, artifact, progress/checkpoint task panjang, dan error recovery.
- Sidebar sesi mendukung search, rename inline, move to trash, restore, dan permanent delete. Operasi memakai pagination/virtualization agar tetap cepat pada history besar.
- API diberi prefix versi (misalnya `/api/v1`) dan tidak diekspos ke network interface selain loopback secara default.
- Tidak ada dependency pada CDN saat runtime; font, script, style, dan icon inti tersedia lokal.

### FR-14 — Visual design: claymorphism

- Design language utama adalah **claymorphism**: bentuk lembut seperti tanah liat, sudut membulat, permukaan tebal, highlight internal, dan shadow berlapis yang memberi kesan objek timbul.
- UI memiliki light dan dark theme dengan CSS custom properties/design tokens untuk warna, radius, elevation, spacing, dan motion.
- Clay effect diprioritaskan pada shell, navigation, composer, cards, dialog, dan tombol. Transcript panjang, code block, diff, terminal output, dan tabel tetap memakai surface yang lebih datar agar mudah dibaca.
- Warna dasar menggunakan palet hangat/netral dengan aksen yang dapat diganti; status success/warning/error tidak hanya dibedakan melalui warna.
- Contrast teks minimal WCAG AA. Focus ring harus jelas dan shadow dekoratif tidak boleh menjadi satu-satunya indikator state.
- Motion halus dan singkat; `prefers-reduced-motion` menonaktifkan animasi non-esensial.
- Layout responsif: sidebar collapsible di layar kecil, composer tetap mudah dijangkau, dan approval dialog tidak menyembunyikan command/diff penting.
- Semua asset visual tersedia lokal; tidak mengambil font/icon dari CDN.

### FR-15 — Git dan GitHub

- Login GitHub memakai OAuth Device Flow yang cocok untuk local app. Client ID boleh public, tetapi token/refresh token wajib berada di OS keyring dan hanya digunakan backend lokal.
- Scope dimulai minimal dan meminta scope tambahan hanya ketika fitur memerlukannya. UI menampilkan akun aktif, scope, expiry/status, dan tombol disconnect/revoke.
- Integrasi boleh mendeteksi GitHub CLI (`gh auth status`) dan menggunakan autentikasinya melalui command yang terkontrol setelah persetujuan, tanpa menyalin token ke prompt atau log.
- Tool Git memisahkan capability: `status/diff/log` (read), `add/reset-path` (index write), `commit` (local history write), serta `push/PR/merge/release` (remote write).
- Sebelum commit, agent wajib membaca status/diff, menghindari secret dan file di luar scope task, serta memastikan tidak ada unresolved conflict.
- Staging menggunakan explicit pathspec; `git add -A` tidak menjadi default.
- Commit lokal menyertakan pesan yang menjelaskan tujuan dan, bila diatur pengguna, author/signing configuration. Arka tidak mengubah global Git identity diam-diam.
- Commit buatan agent tercatat di audit log dengan hash, file yang distage, verification result, dan session ID.
- Push, force-push, merge, tag, release, perubahan branch protection, dan operasi remote destructive berada di luar auto-approved local commit dan membutuhkan capability/approval terpisah.

### FR-16 — Web terminal

- Backend terminal memakai Unix PTY pada macOS/Linux/Termux/PRoot dan ConPTY pada Windows, lalu menjalankan shell lokal sebagai child process daemon dengan working directory eksplisit.
- Frontend memakai terminal emulator yang di-bundle lokal (tanpa CDN), mendukung ANSI color, resize, copy/paste, scrollback terbatas, multiple tabs, dan reconnect.
- Transport terminal menggunakan WebSocket terautentikasi dengan terminal/session ID yang tidak dapat ditebak. Origin/Host policy sama dengan API lokal.
- Shell default dideteksi per platform: PowerShell pada Windows, login shell pengguna pada macOS/Linux, shell Termux pada Android native, dan shell distro pada PRoot.
- Terminal default berada di workspace root. Membuka shell di luar workspace membutuhkan policy/grant sesuai konfigurasi.
- PTY memiliki lifecycle state (`starting`, `running`, `disconnected`, `exited`), idle timeout yang dapat dikonfigurasi, output buffer terbatas, dan cleanup process tree saat ditutup.
- Navigasi web tidak menutup PTY. Disconnect jaringan mempertahankan PTY sampai reconnect grace period/idle timeout.
- Terminal manusia dan `terminal.exec` milik agent adalah capability berbeda. Tidak ada automatic keystroke injection dari agent ke terminal manusia.
- Audit minimal mencatat create/attach/detach/close, cwd, shell, waktu, dan exit status. Transcript terminal tidak disimpan secara default karena dapat mengandung secret; pengguna dapat mengaktifkan recording secara eksplisit.

### FR-17 — Background task manager

- Daemon menyimpan registry task aktif/queued/completed dan memulihkannya dari durable checkpoint setelah restart.
- Task tidak memiliki dependency pada browser connection dan dapat memiliki nol subscriber.
- Default global concurrency adalah 3 task agent; dapat dikonfigurasi berdasarkan resource. Per-session mutating concurrency tetap 1.
- Scheduler menggunakan fair queue agar satu sesi dengan 500+ tool calls tidak memblokir semua sesi lain.
- Tool process, approval wait, model stream, checkpoint, dan cancellation berada dalam context task masing-masing.
- Sidebar mendapatkan status ringkas seluruh task melalui event stream global; transcript/detail tetap dimuat per sesi.
- Stop bersifat scoped ke task/sesi yang dipilih dan tidak membatalkan task pada sesi lain.

### FR-18 — Installation dan lifecycle

- Fresh clone menyediakan `install.sh` untuk macOS/Linux/Termux/PRoot dan `install.ps1` untuk Windows serta quickstart per platform.
- Installer idempotent: aman dijalankan ulang dan tidak menghapus state pengguna.
- Installer mendeteksi OS, architecture, Termux, dan PRoot sebelum memilih artifact/toolchain. Platform yang tidak dikenal berhenti dengan error; tidak boleh diasumsikan sebagai Linux biasa.
- Installer mendeteksi prerequisite yang hilang dan menyiapkannya otomatis. Untuk Go, installer memakai versi sistem yang kompatibel; pada Windows/macOS/Linux/PRoot dapat mengunduh toolchain resmi pinned ke cache lokal, sedangkan pada Termux memakai package `golang` dari repository Termux yang terdeteksi atau verified Android release artifact. Semua download diverifikasi, tidak gagal diam-diam, dan tidak membutuhkan instalasi global di luar lingkungan pengguna.
- Instalasi default bersifat user-local dan tidak membutuhkan root/Administrator/sudo. Path default mengikuti konvensi platform: `%LOCALAPPDATA%`/user PATH pada Windows, `~/.local` pada macOS/Linux/PRoot, dan `$PREFIX`/Termux home pada Termux.
- Lifecycle tidak bergantung pada systemd, launchd, atau Windows Service. MVP menyediakan supervised foreground/background process + lock/PID; integration service OS dapat ditambahkan kemudian.
- `arka open` memakai `start` (Windows), `open` (macOS), `xdg-open` (Linux), dan `termux-open-url` (Termux). Pada PRoot atau command yang tidak tersedia, command mencetak URL tanpa menganggap startup gagal.
- Web UI tidak memiliki build pipeline wajib: source HTML/CSS/ES modules langsung di-embed ke binary. Runtime maupun proses build end user tidak membutuhkan Node.js.
- `arka start` mencegah daemon ganda, menulis PID/lock dengan aman, menunggu health check, lalu membuka browser bila tidak memakai `--no-open`.
- `arka stop`, `status`, `doctor`, dan `uninstall` tersedia pada semua platform; uninstall tidak menghapus data kecuali flag eksplisit diberikan.
- `arka doctor` menampilkan platform/runtime, architecture, writable paths, loopback bind, shell/PTY backend, process cancellation capability, keyring/secrets backend, browser/CDP status, Git, dan tool/skill compatibility.
- Binary release memiliki checksum; installer tidak menjalankan artifact yang gagal diverifikasi.
- Upgrade menjalankan backup/migration dan dapat rollback bila startup health check gagal.

## 10. Non-functional requirements

### Keamanan

- Secure by default; write/execute/network dipisahkan menurut capability.
- Credential dari env/keyring, tidak dari file proyek.
- Permission file state directory dibuat user-only bila OS mendukung.
- Path canonicalization mencegah traversal dan symlink escape.
- Dependency dan binary release memiliki SBOM serta checksum sebelum v1.0.
- Skill dianggap untrusted content: archive/path divalidasi, provenance/checksum disimpan, script tidak berjalan saat inspection, dan semua capability tetap melalui policy engine.
- Threat model formal diselesaikan sebelum opsi akses non-loopback.
- Server menolak Host/Origin asing, memakai CSRF protection untuk mutation, dan menghasilkan local auth token berentropi tinggi.
- API tidak boleh mempercayai UI hanya karena berasal dari browser; authorization tetap diperiksa server-side.
- Mode LAN/remote tidak tersedia pada MVP. Jika ditambahkan kelak, wajib explicit opt-in, TLS, authentication kuat, dan warning.

### Reliability

- Atomic file write: temp file + fsync/rename bila relevan.
- DB migration diuji maju; backup dibuat sebelum migration berisiko.
- Error provider/tool tidak merusak sesi.
- Interrupted turn dapat dipulihkan dan terlihat jelas.

### Portability

Target support MVP:

| Platform | Architecture | Status target | Catatan |
|---|---|---|---|
| Windows 10/11 | amd64, arm64 bila toolchain stabil | Tier 1 | PowerShell, Windows paths, Job Objects, ConPTY |
| macOS | arm64, amd64 | Tier 1 | Apple Silicon + Intel, user-local install |
| Linux | amd64, arm64 | Tier 1 | Tidak mewajibkan systemd |
| Android Termux | arm64; x86_64 best effort | Tier 1 core | Termux paths/shell, Android browser opener, no root |
| Android PRoot distro | arm64, amd64 | Tier 1 core | Linux userland tanpa systemd; beberapa kernel feature terbatas |

- **Tier 1 core** berarti daemon, web UI, provider/model scan, sessions, memory, skills, baseline file/Git/HTTP/code tools, agent terminal execution, dan web terminal wajib berjalan. Tool yang bergantung pada desktop Chromium/OS keyring boleh memakai backend alternatif atau menyatakan unsupported secara jujur.
- Fitur OS-spesifik berada di belakang interface/build tags: paths, process groups/Job Objects, PTY/ConPTY, browser opener, keyring, signals, dan file permissions.
- Gunakan SQLite driver pure-Go dan hindari dependency C toolchain agar build/install Android dan cross-platform tetap ringan.
- Core tidak membutuhkan Python, Node.js, JVM, systemd, root, atau Administrator. Skill tertentu boleh mendeklarasikan dependency tambahan dan compatibility-nya sebelum instalasi.
- File path disimpan dalam bentuk yang tidak mengasumsikan separator Unix; canonicalization dan case-sensitivity diuji per platform.
- Shared storage Android adalah opt-in. State/credential tidak ditempatkan di shared storage.
- Jika OS keyring tidak tersedia di Termux/PRoot, gunakan encrypted local secret store dengan permission paling ketat yang tersedia dan explicit unlock; jangan fallback ke plaintext config.

### Performance

- Lazy-load skill dan context besar.
- Batasi buffer streaming dan tool output.
- Tidak menjalankan embedding service pada idle.
- Benchmark startup, RSS, FTS search, dan agent-loop overhead di CI.

### Accessibility dan UX

- Web UI mendukung keyboard navigation, focus state, semantic HTML, reduced motion, dan kontras WCAG AA.
- Mendukung noninteractive JSON output CLI untuk scripting.
- Error menyebut penyebab, dampak, dan langkah perbaikan.
- Bahasa UI awal Inggris; agent dapat menjawab sesuai bahasa pengguna. Lokalisasi UI Indonesia direncanakan setelah command stabil.

## 11. Rekomendasi teknologi

### Keputusan: Go

Go direkomendasikan untuk produk ini karena:

- satu binary dan deployment mudah ke laptop/VPS;
- startup cepat dan memory footprint relatif kecil;
- concurrency, cancellation, HTTP streaming, dan process control kuat;
- interface cocok untuk provider/tool/plugin boundary;
- cross-compilation dan operasional lebih sederhana daripada runtime Python/Node yang besar.

Kekurangannya: ekosistem AI dan browser automation tidak seluas Python. Solusinya bukan mengganti core, melainkan menempatkan capability berat di MCP/external service atau sidecar. Python tetap dapat dipakai oleh skill/script pengguna, tetapi bukan dependency runtime inti.

### Stack usulan

| Area | Pilihan awal |
|---|---|
| CLI/lifecycle | Cobra atau urfave/cli |
| Web backend | Go `net/http` dengan router minimal; satu origin dengan API |
| Web UI | Standards-based HTML/CSS + vanilla ES modules; tanpa framework/runtime frontend |
| Agent events | SSE dengan durable cursor/replay |
| Web terminal | WebSocket + bundled terminal emulator; PTY/ConPTY backend |
| Asset packaging | `go:embed`; tanpa build frontend, CDN, atau runtime Node dependency |
| Config | YAML + environment override; validasi eksplisit |
| Database | SQLite dengan driver pure-Go untuk Windows/macOS/Linux/Termux/PRoot tanpa C compiler |
| Logging | `log/slog` |
| JSON Schema | Library kecil yang aktif dipelihara; schema tetap provider-neutral |
| HTTP | Standard library terlebih dahulu |
| Secret storage | OS keyring adapter + environment fallback |
| Testing | `testing`, golden tests, fake provider, fake tools |
| Skill packages | Agent Skills-style directory + local lockfile/checksum store |
| RPC/plugin | MCP melalui stdio/HTTP; hindari Go `plugin` karena portabilitas |

Dependency dipilih setelah spike dan dicatat dalam ADR, bukan dikunci hanya oleh PRD.

## 12. Arsitektur tingkat tinggi

```text
┌────────────────── Surfaces ───────────────────┐
│ Browser Web UI (chat + terminal) │ Lifecycle CLI │
└──────────────────────┬────────────────────────┘
                       │ loopback HTTP + SSE/WS
             ┌─────────▼──────────┐
             │ Background Task    │
             │ Manager + Fair Queue│
             └──────┬───────┬─────┘
                    │       │
       ┌────────────▼─┐   ┌─▼────────────────┐
       │ Agent Runtime │   │ PTY/ConPTY       │
       │ context/loop  │   │ Terminal Manager │
       └──────┬────┬───┘   └──────────────────┘
              │    │
┌─────────────▼┐  ┌▼────────────────┐
│ Providers    │  │ Tool Runtime     │
│ native +     │  │ registry/policy │
│ OpenAI-comp. │  │ approval/audit  │
└──────────────┘  └───────┬─────────┘
                           │
             ┌─────────────▼─────────────┐
             │ Built-ins │ MCP │ Backends│
             └───────────────────────────┘

┌──────────────── Persistence ─────────────────┐
│ SQLite sessions/tasks/events/models/FTS/audit│
│ memory store │ skills │ artifacts            │
└──────────────────────────────────────────────┘
```

### Paket Go yang diusulkan

```text
cmd/arka/                 daemon + lifecycle CLI entry point
internal/agent/           turn state machine
internal/context/         assembly, compaction, budgets
internal/provider/        provider interfaces/adapters
internal/tools/           registry, schema, built-ins
internal/policy/          risk and approval decisions
internal/session/         conversation persistence
internal/memory/          declarative memory/retrieval
internal/skills/          discovery, install, lock, update, rollback, loading
internal/storage/         SQLite, migrations, repositories
internal/config/          config and profiles
internal/audit/           redaction and audit records
internal/web/             local HTTP API, SSE/WS, embedded UI
internal/task/            background task registry, queue, checkpoints
internal/terminal/        PTY/ConPTY lifecycle and WebSocket bridge
internal/platform/        OS/runtime detection, paths, process, opener, secrets
web/                      frontend source and assets
internal/ui/              CLI status/diagnostic rendering
pkg/pluginapi/            public contracts only when stable
```

Core package tidak boleh mengimpor implementasi surface. Dependency direction diperiksa dalam review/CI.

## 13. Data model konseptual

Entitas minimum:

- `profiles`
- `provider_profiles`
- `model_catalog`
- `sessions`
- `turns`
- `messages`
- `content_parts`
- `tool_calls`
- `tool_results`
- `background_tasks`
- `task_events`
- `terminal_sessions`
- `approvals`
- `approval_grants`
- `task_checkpoints`
- `memories`
- `memory_sources`
- `skill_sources`
- `skills`
- `skill_versions`
- `skill_installations`
- `skill_runs`
- `artifacts`
- `github_accounts` (metadata only; token berada di OS keyring)
- `repositories`
- `git_operations`
- `usage_records`
- `audit_events`
- `schema_migrations`

ID menggunakan UUID/ULID. Timestamp disimpan UTC. Content besar disimpan sebagai artifact file dengan hash dan referensi dari DB. Retention policy dapat dikonfigurasi.

## 14. Konfigurasi konseptual

```yaml
profile: default
provider_profiles:
  omniroute:
    type: openai-compatible
    base_url: https://example.invalid/v1
    api_key_env: OMNIROUTE_API_KEY
    models_path: /models
    default_model: model-name
agent:
  max_concurrent_sessions: 3
  max_tool_calls_per_turn: 1000
  checkpoint_every_tool_calls: 25
  turn_timeout: 2h
  context_budget_tokens: 32000
  daily_cost_limit_usd: 5
workspace:
  root: .
  allow_outside_root: false
approval:
  profile: standard
  grants:
    trusted_workspace: true
    local_git_commit: allow
    remote_git_write: ask
  destructive: deny
memory:
  enabled: true
  retrieval_limit: 12
skills:
  install_scope: profile
  require_approval_for_new_source: true
  allow_auto_update: false
  catalogs: []
logging:
  level: info
  redact_secrets: true
web:
  listen: 127.0.0.1:7331
  open_browser: true
  allow_remote: false
terminal:
  enabled: true
  default_cwd: workspace
  reconnect_grace: 10m
  idle_timeout: 2h
  record_transcript: false
```

Nama dan schema final ditetapkan melalui ADR dan usability test CLI.

## 15. Acceptance criteria MVP

MVP dianggap selesai bila seluruh skenario berikut lulus otomatis atau terdokumentasi sebagai test manual:

1. **Install:** dari fresh clone pada image OS bersih tanpa Go/Node, satu installer mengunduh toolchain pinned yang diperlukan, memverifikasinya, menghasilkan aplikasi runnable tanpa langkah build manual, dan aman dijalankan ulang.
2. **Local web:** `arka start` membuka web UI loopback; assets tersedia tanpa internet/CDN dan health endpoint sukses.
3. **Setup:** dari web wizard, pengguna mengonfigurasi endpoint/model dan menerima respons streaming.
4. **Tool loop:** fake model meminta baca file lalu patch; runtime mengeksekusi keduanya dan model memberi final response.
5. **Validation:** argumen tool invalid ditolak tanpa menjalankan handler.
6. **Approval persistence:** pada profile Standard, write/command yang masuk rule trusted workspace tidak meminta approval ulang pada chat berikutnya; destructive dan remote Git write tetap meminta approval.
7. **Long task:** satu sesi berhasil mencatat lebih dari 500 tool calls, membuat checkpoint berkala, dapat reconnect dari browser, dan tidak kehilangan audit trail.
8. **Workspace boundary:** `../` dan symlink escape ditolak.
9. **Terminal:** timeout dan cancellation membunuh process tree.
10. **Persistence:** kill setelah tool result lalu restart tidak membuat turn selesai palsu atau menjalankan ulang command.
11. **Memory:** fakta yang disimpan pada sesi A dapat ditemukan relevan pada sesi B dan memiliki source.
12. **Session search:** query FTS mengembalikan pesan serta session ID yang benar.
13. **Context limit:** history besar dikompaksi tanpa membuang policy dan objective terbaru.
14. **Secret safety:** API key tidak muncul pada log, transcript, error, atau snapshot test.
15. **Portability:** unit/integration test lulus di Linux, macOS, dan Windows CI.
16. **Doctor:** config/provider/DB/tool failure menghasilkan diagnosis actionable.
17. **Resource target:** benchmark startup/RSS berada dalam target atau memiliki waiver tertulis sebelum release.
18. **Session management:** pengguna dapat rename, soft-delete, restore, dan purge sesi; operasi tetap cepat pada fixture sesi 500+ tool calls.
19. **GitHub:** Device Flow berhasil connect/disconnect tanpa token muncul di browser response, log, transcript, atau prompt.
20. **Agent commit:** agent men-stage hanya path terkait, membuat commit lokal yang dapat diaudit, dan tidak melakukan push tanpa approval remote terpisah.
21. **Provider presets:** setiap preset dapat divalidasi dan menampilkan error auth/endpoint secara actionable.
22. **OmniRoute/custom API:** profile Custom OpenAI-compatible dapat menyimpan base URL/key reference, scan model dari endpoint, menerima model ID manual, lalu menyelesaikan tool-calling chat.
23. **Model switch:** model dapat diganti per sesi; request aktif selesai dengan model lama dan call berikutnya memakai model baru, dengan kedua ID tercatat.
24. **Background sessions:** task sesi A terus berjalan ketika UI membuka sesi B; kedua sesi dapat berjalan concurrent dan reconnect tidak menggandakan eksekusi.
25. **Web terminal:** pengguna dapat membuat PTY, menjalankan command interaktif, resize, pindah halaman, reconnect, dan menutup process tree dari web tanpa CDN.
26. **Terminal isolation:** agent tidak dapat mengetik/membaca terminal manusia tanpa capability eksplisit, dan transcript terminal default-nya tidak disimpan.
27. **Local tool boundary:** production browser bundle tidak memiliki/mock handler tool; file, terminal, Git, provider, memory, skill, code, dan MCP calls terbukti melalui daemon trace + audit ID.
28. **Complete baseline tools:** contract/integration tests lulus untuk setiap toolset MVP, termasuk validation, approval, cancellation, timeout, output cap, error, dan audit path.
29. **Skill install:** catalog, Git subdirectory, local folder, dan archive dapat di-inspect lalu di-install secara atomik ke scope yang dipilih.
30. **Skill lifecycle:** enable/disable, update dengan diff, rollback, dan uninstall bekerja tanpa merusak versi atau sesi lain.
31. **Malicious skill:** traversal, symlink escape, archive bomb, checksum mismatch, undeclared executable/dependency, dan install hook otomatis ditolak.
32. **Skill execution:** script skill berjalan sebagai supervised local tool process dan tidak dapat melewati workspace/policy/approval boundary.
33. **Local browser tool:** agent mengontrol Chrome/managed Chromium nyata dari daemon untuk navigate/click/input/screenshot; tidak ada mocked browser result dari frontend.
34. **Windows:** fresh clone + `install.ps1` + `arka start` berhasil pada clean Windows test environment; file tools, process cancellation, ConPTY, Git, SQLite, dan web UI lulus smoke test.
35. **macOS:** fresh clone + `install.sh` + `arka start` berhasil pada Apple Silicon dan Intel CI/runner yang tersedia; PTY, opener, Git, SQLite, dan web UI lulus smoke test.
36. **Termux:** pada Android/Termux arm64 tanpa root, installer, daemon loopback, provider chat, SQLite, file/Git/code tools, skills, background session, dan web terminal lulus device smoke test.
37. **PRoot:** pada Debian/Ubuntu PRoot tanpa systemd, installer, PID lifecycle, daemon loopback, provider chat, SQLite, file/Git/code tools, skills, background session, dan web terminal lulus device smoke test.
38. **Android limitation honesty:** browser/keyring/kernel-dependent capability yang tidak tersedia dilaporkan `unsupported` dengan remediation; tidak ada mock success.
39. **Portable state safety:** state/secret Android tetap di app-private storage secara default dan Windows/macOS/Linux memakai user-local permission yang sesuai.

## 16. Testing strategy

- Unit test untuk state transition, policy, schema validation, context budget, redaction, dan path handling.
- Contract test untuk setiap provider menggunakan fixture streaming dan tool calls.
- Integration test dengan fake deterministic LLM server dan fake model catalog untuk semua adapter/provider profile.
- Compatibility test Custom OpenAI API mencakup OmniRoute-style base URL, `/models`, streaming, tool call, auth header, dan manual model fallback.
- Concurrency/reconnect test menjalankan task pada beberapa sesi sambil route UI berpindah dan SSE disconnect/reconnect.
- PTY integration test mencakup resize, reconnect, process cleanup, backpressure, dan terminal isolation.
- Skill package fixtures mencakup catalog/Git/local/archive, atomic failure, update diff, rollback, conflict, dan malicious archive/path cases.
- Local browser integration test menggunakan halaman fixture untuk navigate/snapshot/click/input/download/screenshot dan memverifikasi process berasal dari daemon.
- Production bundle test memastikan tidak ada mock tool transport/handler dan seluruh tool event memiliki daemon trace/audit ID.
- CI matrix minimal: Windows amd64, macOS arm64/amd64 sesuai runner, Linux amd64/arm64; compile checks untuk target tambahan yang didukung.
- Release gate mencakup smoke test perangkat Android Termux arm64 dan distro PRoot nyata/emulator yang terdokumentasi, karena container Linux biasa tidak cukup mewakili Android/PRoot.
- Platform contract tests untuk path, case sensitivity, signals/process tree, PTY/ConPTY, browser opener, secret backend, lock/PID, dan uninstall-preserves-data.
- Golden transcript test untuk agent loop.
- Property/fuzz test untuk parser stream, path, JSON args, redaction, dan migration input.
- Fault injection pada DB transaction, interrupted stream, process cancellation, dan disk-full path.
- Security test untuk command injection, prompt-based policy bypass, traversal, symlink escape, dan secret exfiltration.
- Opt-in live provider eval agar CI utama tetap deterministik dan murah.

## 17. Release plan

### Milestone 0 — Foundation (1–2 minggu)

- ADR awal, Go module, CI, config, logging, SQLite migration, domain interfaces.
- Installer skeleton Windows/macOS/Linux/Termux/PRoot, platform abstraction, local HTTP server, embedded static proof-of-concept, dan health check.
- Fake provider dan test harness.

### Milestone 1 — Usable local web agent (2–3 minggu)

- Provider presets, Custom OpenAI-compatible/OmniRoute, model scan, dan streaming/tool calling.
- Web chat claymorphism minimal, provider/model picker, session persistence + rename/delete, background task manager, event replay/reconnect, cancellation, context builder.
- Web terminal MVP dengan PTY, bundled emulator, multiple tabs, dan reconnect.
- Local tool registry/dispatcher dan contract tests; browser hanya menggunakan API/event transport.
- Lifecycle CLI dan browser auto-open.
- File tools read-only.

### Milestone 2 — Safe action (2–3 minggu)

- File write/patch, terminal/process, policy profiles + persistent workspace grants, audit, redaction.
- GitHub Device Flow dan tool local Git commit.
- Recovery dan acceptance tests kritis.

### Milestone 3 — Memory and skills (2–3 minggu)

- FTS session search, declarative memory, skill catalog/Git/local/archive installation, lifecycle management, discovery/loading, dan compaction.
- Release `v0.1.0` setelah security review.

Estimasi adalah urutan perencanaan, bukan komitmen tanggal. Scope dipotong sebelum quality/safety.

## 18. Risiko dan mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Model mengarang tool/argumen | Aksi salah | Registry allowlist, JSON schema, policy, approval |
| Prompt injection dari web/file | Exfiltration/aksi berbahaya | Tandai untrusted content, capability boundary, jangan izinkan content mengubah policy |
| Shell terlalu kuat | Kerusakan host | Risk-based policy, trusted workspace boundary, sandbox backend, deny patterns |
| Memory menyimpan fakta salah | Personalisasi memburuk | Source/confidence, pin, candidate update, inspect/delete UI |
| Context membengkak | Mahal/lambat | Budget, lazy skill load, truncate artifacts, compaction |
| Provider API berbeda | Fragile adapter/model scan gagal | Contract interface, provider test suite, capability `unknown`, cache, manual model ID |
| Pindah sesi membatalkan task | Pekerjaan panjang hilang | Daemon-owned background task, durable cursor/checkpoint, disconnect test |
| Banyak sesi/terminal membebani mesin | Resource exhaustion | Global concurrency, fair queue, idle timeout, output cap, process cleanup |
| Web terminal disalahgunakan | Host compromise/secret leak | Loopback auth, workspace policy, PTY isolation, no recording default, agent separation |
| Plugin/skill supply-chain | Remote code execution | Provenance, checksum/lock, inspect-before-install, archive validation, no install hook, policy sandbox |
| Browser/tool mock masuk production | Tool tampak sukses tetapi tidak melakukan aksi nyata | Production build guard, daemon audit ID requirement, real local browser E2E tests |
| Local web server terekspos ke LAN | Akses data/tool tanpa izin | Loopback-only, Host/Origin validation, auth token, remote mode off |
| Installer supply-chain/partial failure | Binary berbahaya atau instalasi rusak | Pinned versions, checksum, atomic install, health check, rollback |
| Android/PRoot dianggap Linux desktop biasa | Fitur rusak atau fake success | Runtime detection, Tier 1 core matrix, capability probe, real-device smoke tests |
| Process cleanup berbeda antar-OS | Orphan process | Job Objects/process groups, PRoot fallback, fault-injection per platform |
| Credential tanpa OS keyring di Android | Secret bocor | Encrypted local store, app-private path, strict permissions, explicit unlock |
| Self-learning menurunkan kualitas | Skill rusak | Draft/approval, version history, rollback, eval sebelum promote |
| Approval terlalu sering | Task panjang terganggu | Session/workspace grants, Standard profile, approval berdasarkan risk bukan tool count |
| Agent membuat commit yang salah | History tercemar/secret ikut stage | Explicit path staging, pre-commit diff, secret scan, audit; push terpisah |
| Task 500+ tool calls boros/looping | Biaya dan waktu tak terkendali | Checkpoint, compaction, repeated-failure guard, token/cost/deadline budget, tombol Stop |
| Scope meniru semua Hermes | Proyek tidak selesai | Milestone sempit, parity bukan target MVP |

## 19. Keputusan yang sudah diambil

1. Implementasi baru dimulai dari repository kosong.
2. Bahasa core: **Go**.
3. Produk **web-first, local-runtime, dan local-storage**; web hanya interface.
4. Tidak ada integrasi WhatsApp, Telegram, Discord, atau messaging gateway.
5. MVP hanya satu user dan satu daemon lokal.
6. SQLite adalah source of truth lokal.
7. OpenAI-compatible adalah compatibility contract utama; preset native OpenAI, Anthropic, dan Gemini serta preset umum lainnya tersedia dari MVP.
8. Tool execution memakai policy + approval + audit.
9. Fresh clone memiliki one-command installer; release production memakai embedded web assets.
10. “Berpikir seperti agent” diwujudkan sebagai state machine, plan/checklist, tool feedback, dan verification—not penyimpanan chain-of-thought.
11. Learning loop tidak boleh mengubah security policy atau binary secara mandiri.
12. Web UI menggunakan visual language claymorphism yang accessible.
13. Session mendukung rename, soft delete/restore, dan permanent purge.
14. Agent dapat login GitHub dan membuat local Git commit; remote writes memiliki approval terpisah.
15. Session mendukung lebih dari 500 tool calls; default ceiling per task 1.000 dan total session tidak dibatasi hitungan tool.
16. Approval bersifat risk-based dengan grant session/workspace, bukan diminta ulang setiap chat.
17. Provider umum tersedia sebagai preset; OmniRoute menggunakan Custom OpenAI-compatible profile.
18. Model dapat di-scan, dimasukkan manual, dipilih per sesi, dan diganti tanpa membuat chat baru.
19. Background task tidak berhenti ketika pengguna berpindah sesi atau menutup tab.
20. Web terminal lokal adalah fitur MVP dan terpisah dari terminal tool milik agent.
21. Skill dapat dipasang dari katalog, Git, local folder, atau archive dan dikelola penuh dari web.
22. Seluruh production tool calling berjalan di daemon/proses lokal; browser hanya interface dan tidak memiliki mock/fallback executor.
23. Repository dan one-command installer mendukung Windows, macOS, Android Termux, dan distro Android PRoot selain Linux biasa.
24. Core tidak bergantung pada systemd, root/Admin, Node.js, Python, JVM, atau C compiler.
25. Pada platform yang tidak mendukung capability tertentu, UI menyatakan `unsupported`; tidak pernah mengganti aksi nyata dengan mock.

## 20. Open questions untuk owner

Hal berikut perlu diputuskan sebelum implementasi di luar foundation:

1. Target penggunaan pertama: coding agent, personal assistant umum, atau automation server?
2. Provider preset mana yang wajib diuji live saat MVP selain Custom OpenAI-compatible/OmniRoute?
3. Lisensi: MIT, Apache-2.0, atau proprietary?
4. Versi minimum Windows, macOS, Android, dan Termux apa yang akan dijadikan support floor final?
5. Default execution: host lokal dengan approval atau container-first?
6. Apakah skill hasil belajar auto-save, selalu approval, atau berbeda per profile?
7. Apakah akses harus selamanya loopback-only, atau LAN mode dibutuhkan setelah v1.0?
8. Bahasa web UI utama: Inggris, Indonesia, atau bilingual?
9. Batas biaya default dan telemetry: sepenuhnya off atau anonymous opt-in?
10. Apakah Git commit wajib signed bila GPG/SSH signing tersedia, atau mengikuti config repository?
11. Apakah trash memiliki retention otomatis (misalnya 30 hari) atau hanya purge manual?
12. Berapa global concurrency default yang cocok untuk target mesin minimum: 2, 3, atau adaptif?
13. Apakah terminal recording perlu tersedia pada MVP atau ditunda setelah security review?
14. Katalog skill default mana yang akan dipercaya/ditampilkan saat first run?
15. Apakah dependency installer skill boleh aktif di MVP atau skill dengan dependency eksternal hanya diberi petunjuk manual terlebih dahulu?
16. Untuk Android, apakah arm64 menjadi satu-satunya release resmi awal atau x86_64 emulator juga wajib dirilis?
17. Apakah browser automation Android wajib melalui remote-debugging browser pengguna atau cukup `unsupported` sampai backend stabil?

## 21. Referensi riset

Referensi hanya untuk memahami pola produk dan kebutuhan, bukan untuk menyalin implementasi:

- Hermes Agent repository dan feature overview: https://github.com/NousResearch/hermes-agent
- Hermes architecture/development guide: https://github.com/NousResearch/hermes-agent/blob/main/AGENTS.md
- Hermes tools/toolsets: https://github.com/NousResearch/hermes-agent/blob/main/website/docs/user-guide/features/tools.md
- Hermes persistent memory: https://github.com/NousResearch/hermes-agent/blob/main/website/docs/user-guide/features/memory.md
- Model Context Protocol: https://modelcontextprotocol.io/
- Agent Skills specification: https://agentskills.io/

---

## Appendix A — Contoh event agent

```json
{
  "type": "tool.requested",
  "trace_id": "01...",
  "turn_id": "01...",
  "tool": "terminal.exec",
  "risk": "execute",
  "arguments_preview": {"command": ["go", "test", "./..."]},
  "timestamp": "2026-10-04T12:00:00Z"
}
```

Event publik tidak memuat reasoning privat. Surface menerima event seperti `message.delta`, `plan.updated`, `tool.requested`, `approval.required`, `tool.started`, `tool.completed`, `usage.updated`, `turn.completed`, dan `turn.failed`.

## Appendix B — Definition of done per fitur

Sebuah fitur dianggap selesai jika:

- requirement dan threat considerations ditulis;
- API/interface memiliki contract test;
- happy path, error path, cancel, dan timeout diuji;
- log dan error bebas secret;
- docs CLI/config diperbarui;
- migration/backward compatibility diperiksa;
- benchmark terkait tidak mengalami regresi tanpa alasan;
- acceptance criteria dapat ditelusuri ke test atau bukti manual.
