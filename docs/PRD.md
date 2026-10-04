# Product Requirements Document — Nadi Agent

| Atribut | Nilai |
|---|---|
| Status | Draft v0.1 untuk review |
| Tanggal | 4 Oktober 2026 |
| Codename | **Nadi** (dapat diganti) |
| Tipe produk | Personal AI agent, local-first, terminal-first |
| Implementasi utama | Go |
| Lisensi | Belum diputuskan |

## 1. Ringkasan

Nadi adalah personal AI agent yang berjalan sebagai satu runtime ringan, dapat memakai model LLM dari berbagai provider, mengeksekusi tools, mengingat informasi lintas sesi, mempelajari prosedur sebagai skill, mendelegasikan pekerjaan, dan kelak diakses dari CLI maupun kanal pesan.

Inspirasi perilakunya adalah pola produk Hermes Agent: satu inti agent untuk banyak antarmuka, loop penggunaan tool, persistent memory, reusable skills, pencarian sesi lama, delegasi subagent, dan automasi terjadwal. Nadi harus merupakan implementasi mandiri dari nol—bukan salinan kode, nama, prompt, atau identitas Hermes.

Produk tahap pertama berfokus pada inti yang kecil tetapi benar: **CLI + agent loop + tools + persistence + safety**. Fitur luas seperti gateway pesan, browser visual, voice, dan multi-agent persisten ditambahkan setelah fondasi stabil.

## 2. Masalah

Chatbot biasa memiliki beberapa kelemahan:

1. Tidak dapat melakukan tindakan nyata secara konsisten.
2. Kehilangan konteks saat sesi berakhir.
3. Mengulang cara kerja yang sama karena tidak menyimpan prosedur yang berhasil.
4. Terikat pada satu provider/model.
5. Tidak transparan ketika menjalankan perintah berisiko.
6. Sulit digunakan sebagai proses jangka panjang dari terminal, server, dan kanal pesan sekaligus.

Nadi mengatasi masalah tersebut dengan runtime agent yang persisten, provider-agnostic, tool-driven, dapat diaudit, dan aman secara default.

## 3. Visi produk

> Satu agent pribadi yang ringan, tumbuh bersama penggunanya, dapat bertindak dengan aman, dan dapat dijalankan di mana pun sebagai satu binary.

### Prinsip produk

1. **Narrow core, extensible edges** — inti hanya berisi orkestrasi universal; kemampuan khusus hidup sebagai tool, skill, atau plugin.
2. **Local-first** — sesi, konfigurasi, memory, dan audit log disimpan lokal secara default.
3. **Human control** — aksi destruktif, privilege escalation, dan akses sensitif membutuhkan kebijakan atau persetujuan.
4. **Observable** — pengguna dapat melihat tool yang dipanggil, argumen, hasil, durasi, dan perubahan file.
5. **Provider-agnostic** — mendukung API OpenAI-compatible terlebih dahulu, lalu adapter native.
6. **Progressive disclosure** — skill dan memory dimuat sesuai kebutuhan agar context window hemat.
7. **Useful learning, not uncontrolled self-modification** — agent boleh membuat draft memory/skill, tetapi tidak mengubah binary atau security policy sendiri.
8. **No hidden-reasoning dependency** — produk tidak meminta atau menyimpan chain-of-thought privat model. Keputusan dijelaskan lewat plan singkat, tindakan, hasil, dan ringkasan alasan yang dapat diaudit.

## 4. Pengguna sasaran

### Persona utama: developer/power user

- Nyaman menggunakan terminal.
- Ingin agent mengedit file, menjalankan command, mencari web, dan mengotomasi pekerjaan.
- Memerlukan kontrol atas model, biaya, data, dan izin.
- Menjalankan agent di laptop, VPS kecil, atau container.

### Persona sekunder: operator/pemilik bisnis kecil

- Ingin automasi terjadwal dan akses melalui Telegram/Discord.
- Membutuhkan agent yang mengingat SOP dan preferensi.
- Tidak ingin mengelola stack berat.

## 5. Sasaran dan metrik

### Sasaran MVP

1. Instalasi menghasilkan satu binary Go dan dapat mulai chat dalam kurang dari 5 menit setelah konfigurasi provider.
2. Agent menyelesaikan task multi-step dengan loop tool tanpa restart proses.
3. Sesi dan memory bertahan setelah aplikasi ditutup.
4. Semua tool call tervalidasi, dibatasi timeout, dan tercatat di audit log.
5. Command berisiko meminta approval secara konsisten.
6. Model/provider dapat diganti tanpa mengubah agent core.

### Metrik keberhasilan

| Metrik | Target MVP |
|---|---:|
| Cold start CLI di mesin umum | < 250 ms, di luar latency jaringan/model |
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

- CLI interaktif dengan streaming output.
- Setup/config noninteraktif melalui command dan environment variable.
- Adapter API OpenAI-compatible.
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
  - HTTP fetch/web extract sederhana.
- SQLite untuk sessions, messages, tool runs, approvals, dan full-text search.
- Memory ringkas berbasis file atau record terstruktur.
- Skills berbasis direktori dan Markdown, read-only saat eksekusi MVP.
- Project context discovery (`AGENTS.md`, `.nadi.md`, dan file konfigurasi yang disetujui).
- Structured logs, audit log, dan redaction secret.
- Approval mode: `ask`, `allow`, `deny` berdasarkan rule.
- Commands dasar: `nadi`, `nadi setup`, `nadi doctor`, `nadi sessions`, `nadi tools`, `nadi config`.

### 6.2 v0.2 — Learning loop

- Membuat dan memperbarui skill sebagai draft dari task kompleks yang berhasil.
- Memory curation: tambah, ganti, hapus, pin, dan deduplikasi.
- Nudge setelah sejumlah turn/tool call untuk mengevaluasi hal yang layak disimpan.
- Session summarization dan context compaction.
- MCP client untuk external tools.
- Import/export portable untuk session, memory, dan skill.
- Provider adapter Anthropic serta model lokal/Ollama bila adapter kompatibilitas tidak cukup.

### 6.3 v0.3 — Automation dan delegation

- Cron scheduler persisten.
- Subagent terisolasi dengan budget, toolset, dan working directory sendiri.
- Parallel delegation dengan concurrency limit.
- Delivery hasil task ke kanal asal.
- Webhook/API server yang kompatibel dengan format chat umum.

### 6.4 v1.0 — Multi-surface agent

- Gateway Telegram dan Discord; kanal lain melalui plugin.
- TUI matang.
- Browser automation sebagai service/plugin terisolasi.
- Profiles terpisah untuk work/personal/team.
- Plugin SDK stabil dan versioned.
- Container/SSH execution backend.
- Dashboard opsional, bukan dependency core.

### Di luar scope awal

- Melatih foundation model sendiri.
- Menyamai seluruh fitur Hermes pada rilis pertama.
- Autonomous binary/source self-modification.
- Desktop/mobile GUI native.
- Menyimpan atau menampilkan private chain-of-thought model.
- Menjalankan shell tanpa policy, timeout, output cap, dan audit.
- Menjadi platform multi-tenant SaaS pada MVP.

## 7. Pengalaman pengguna utama

### 7.1 First run

1. Pengguna menjalankan `nadi setup`.
2. Nadi meminta endpoint, model, dan metode credential.
3. Credential disimpan lewat environment/OS keyring; tidak ditulis plaintext ke config jika keyring tersedia.
4. Pengguna memilih workspace dan mode approval.
5. `nadi doctor` menguji provider, DB, permission, dan tools.
6. Pengguna menjalankan `nadi` dan mulai chat.

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

### 7.4 Melanjutkan sesi

- Pengguna melihat sesi dengan `nadi sessions list`.
- Pengguna dapat resume berdasarkan ID.
- Sistem memuat transcript, summary, project context, dan memory relevan.
- Tool call yang terputus ditandai `interrupted`, tidak diam-diam dijalankan ulang.

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

- Maksimum 20 tool calls per turn.
- Maksimum 3 kegagalan berulang dengan signature sama.
- Timeout per tool dan deadline per turn.
- Batas bytes output tool.
- Batas biaya/token per turn dan per hari bila provider menyediakan usage.
- Cancellation melalui `Ctrl+C` yang mengalir ke model stream dan child process.

## 9. Functional requirements

### FR-1 — Model provider

- Interface provider mencakup streaming message, structured tool call, usage, cancellation, dan error classification.
- MVP mendukung endpoint OpenAI-compatible dengan custom base URL dan headers aman.
- Provider config tidak boleh bocor ke prompt/log.
- Retry hanya untuk error transient dan harus memakai exponential backoff + jitter.

### FR-2 — Agent runtime

- Agent loop deterministik pada level state transition.
- Setiap turn memiliki ID dan trace ID.
- Runtime mendukung stop, resume session, max steps, token budget, dan deadline.
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
- handler dan error terstruktur.

Tool registry mendukung enable/disable per profile. Hasil tool harus memiliki status, output terpotong, metadata durasi, dan artifact reference bila output besar.

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
- Pencarian history memakai SQLite FTS5.
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

### FR-8 — Project context

- Cari context dari current directory menuju root workspace dengan precedence terdokumentasi.
- Context tidak boleh dapat menurunkan system security policy.
- Ukuran, jumlah file, dan symlink traversal dibatasi.
- Pengguna dapat melihat asal semua context yang aktif.

### FR-9 — Approval

Prompt approval menampilkan:

- tool dan risk level;
- aksi/command sebenarnya;
- workspace/host target;
- alasan singkat;
- pilihan `allow once`, `allow session`, `deny`;
- untuk rule yang aman, opsi membuat allow rule terbatas.

Approval persisten harus spesifik dan dapat dicabut. Pola wildcard luas diberi peringatan.

### FR-10 — Observability

- Console menampilkan event streaming tanpa membocorkan secret.
- Structured log menggunakan JSON dan level.
- Audit record append-oriented untuk aksi eksternal.
- `nadi doctor` memeriksa config, DB migration, provider connectivity, tool availability, dan permission.
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

### FR-13 — Gateway (v1.0)

- Semua kanal menggunakan agent core yang sama.
- Identity mapping dan DM pairing mencegah pengguna asing mengakses agent.
- Session continuity dapat dikonfigurasi lintas kanal.
- Rate limit, message size limit, attachment policy, dan allowlist wajib tersedia.

## 10. Non-functional requirements

### Keamanan

- Secure by default; write/execute/network dipisahkan menurut capability.
- Credential dari env/keyring, tidak dari file proyek.
- Permission file state directory dibuat user-only bila OS mendukung.
- Path canonicalization mencegah traversal dan symlink escape.
- Dependency dan binary release memiliki SBOM serta checksum sebelum v1.0.
- Threat model formal diselesaikan sebelum gateway publik.

### Reliability

- Atomic file write: temp file + fsync/rename bila relevan.
- DB migration diuji maju; backup dibuat sebelum migration berisiko.
- Error provider/tool tidak merusak sesi.
- Interrupted turn dapat dipulihkan dan terlihat jelas.

### Portability

- Target awal: Linux amd64/arm64, macOS amd64/arm64, Windows amd64.
- Fitur OS-spesifik di belakang interface dan build tags.
- Core tidak membutuhkan Python, Node.js, atau JVM.

### Performance

- Lazy-load skill dan context besar.
- Batasi buffer streaming dan tool output.
- Tidak menjalankan embedding service pada idle.
- Benchmark startup, RSS, FTS search, dan agent-loop overhead di CI.

### Accessibility dan UX

- Warna dapat dimatikan dan output tetap bermakna tanpa warna.
- Mendukung noninteractive JSON output untuk scripting.
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
| CLI | Cobra atau urfave/cli |
| Interactive UI | Bubble Tea + Lip Gloss (setelah loop dasar stabil) |
| Config | YAML + environment override; validasi eksplisit |
| Database | SQLite; driver pure-Go bila stabilitas/benchmark memadai |
| Logging | `log/slog` |
| JSON Schema | Library kecil yang aktif dipelihara; schema tetap provider-neutral |
| HTTP | Standard library terlebih dahulu |
| Secret storage | OS keyring adapter + environment fallback |
| Testing | `testing`, golden tests, fake provider, fake tools |
| RPC/plugin | MCP melalui stdio/HTTP; hindari Go `plugin` karena portabilitas |

Dependency dipilih setelah spike dan dicatat dalam ADR, bukan dikunci hanya oleh PRD.

## 12. Arsitektur tingkat tinggi

```text
┌──────────── Surfaces ────────────┐
│ CLI/TUI │ API │ Gateway (later)  │
└────────────────┬─────────────────┘
                 │ normalized input/events
┌────────────────▼─────────────────┐
│ Agent Runtime                    │
│ turn state machine               │
│ context builder + budget         │
│ model/tool orchestration         │
└───────┬───────────────┬──────────┘
        │               │
┌───────▼──────┐  ┌─────▼──────────┐
│ Providers    │  │ Tool Runtime    │
│ OpenAI-comp. │  │ registry/policy│
│ native later │  │ approval/audit │
└──────────────┘  └─────┬──────────┘
                        │
          ┌─────────────▼─────────────┐
          │ Built-ins │ MCP │ Backends│
          └───────────────────────────┘

┌──────────── Persistence ─────────────┐
│ SQLite sessions/FTS/audit            │
│ memory store │ skills │ artifacts    │
└──────────────────────────────────────┘
```

### Paket Go yang diusulkan

```text
cmd/nadi/                 entry point
internal/agent/           turn state machine
internal/context/         assembly, compaction, budgets
internal/provider/        provider interfaces/adapters
internal/tools/           registry, schema, built-ins
internal/policy/          risk and approval decisions
internal/session/         conversation persistence
internal/memory/          declarative memory/retrieval
internal/skills/          discovery and loading
internal/storage/         SQLite, migrations, repositories
internal/config/          config and profiles
internal/audit/           redaction and audit records
internal/ui/              terminal rendering
pkg/pluginapi/            public contracts only when stable
```

Core package tidak boleh mengimpor implementasi surface. Dependency direction diperiksa dalam review/CI.

## 13. Data model konseptual

Entitas minimum:

- `profiles`
- `sessions`
- `turns`
- `messages`
- `content_parts`
- `tool_calls`
- `tool_results`
- `approvals`
- `memories`
- `memory_sources`
- `skills`
- `skill_runs`
- `artifacts`
- `usage_records`
- `audit_events`
- `schema_migrations`

ID menggunakan UUID/ULID. Timestamp disimpan UTC. Content besar disimpan sebagai artifact file dengan hash dan referensi dari DB. Retention policy dapat dikonfigurasi.

## 14. Konfigurasi konseptual

```yaml
profile: default
provider:
  type: openai-compatible
  base_url: https://example.invalid/v1
  model: model-name
  api_key_env: NADI_API_KEY
agent:
  max_tool_calls: 20
  turn_timeout: 15m
  context_budget_tokens: 32000
  daily_cost_limit_usd: 5
workspace:
  root: .
  allow_outside_root: false
approval:
  default: ask
  read: allow
  write: ask
  execute: ask
  destructive: deny
memory:
  enabled: true
  retrieval_limit: 12
logging:
  level: info
  redact_secrets: true
```

Nama dan schema final ditetapkan melalui ADR dan usability test CLI.

## 15. Acceptance criteria MVP

MVP dianggap selesai bila seluruh skenario berikut lulus otomatis atau terdokumentasi sebagai test manual:

1. **Setup:** dari mesin bersih, pengguna mengonfigurasi endpoint/model dan menerima respons streaming.
2. **Tool loop:** fake model meminta baca file lalu patch; runtime mengeksekusi keduanya dan model memberi final response.
3. **Validation:** argumen tool invalid ditolak tanpa menjalankan handler.
4. **Approval:** read diizinkan; write meminta approval; destructive ditolak pada default policy.
5. **Workspace boundary:** `../` dan symlink escape ditolak.
6. **Terminal:** timeout dan cancellation membunuh process tree.
7. **Persistence:** kill setelah tool result lalu restart tidak membuat turn selesai palsu atau menjalankan ulang command.
8. **Memory:** fakta yang disimpan pada sesi A dapat ditemukan relevan pada sesi B dan memiliki source.
9. **Session search:** query FTS mengembalikan pesan serta session ID yang benar.
10. **Context limit:** history besar dikompaksi tanpa membuang policy dan objective terbaru.
11. **Secret safety:** API key tidak muncul pada log, transcript, error, atau snapshot test.
12. **Portability:** unit/integration test lulus di Linux, macOS, dan Windows CI.
13. **Doctor:** config/provider/DB/tool failure menghasilkan diagnosis actionable.
14. **Resource target:** benchmark startup/RSS berada dalam target atau memiliki waiver tertulis sebelum release.

## 16. Testing strategy

- Unit test untuk state transition, policy, schema validation, context budget, redaction, dan path handling.
- Contract test untuk setiap provider menggunakan fixture streaming dan tool calls.
- Integration test dengan fake deterministic LLM server.
- Golden transcript test untuk agent loop.
- Property/fuzz test untuk parser stream, path, JSON args, redaction, dan migration input.
- Fault injection pada DB transaction, interrupted stream, process cancellation, dan disk-full path.
- Security test untuk command injection, prompt-based policy bypass, traversal, symlink escape, dan secret exfiltration.
- Opt-in live provider eval agar CI utama tetap deterministik dan murah.

## 17. Release plan

### Milestone 0 — Foundation (1–2 minggu)

- ADR awal, Go module, CI, config, logging, SQLite migration, domain interfaces.
- Fake provider dan test harness.

### Milestone 1 — Usable agent loop (2–3 minggu)

- OpenAI-compatible streaming/tool calling.
- CLI minimal, session persistence, cancellation, context builder.
- File tools read-only.

### Milestone 2 — Safe action (2–3 minggu)

- File write/patch, terminal/process, policy, approval, audit, redaction.
- Recovery dan acceptance tests kritis.

### Milestone 3 — Memory and skills (2–3 minggu)

- FTS session search, declarative memory, skill discovery/loading, compaction.
- Release `v0.1.0` setelah security review.

Estimasi adalah urutan perencanaan, bukan komitmen tanggal. Scope dipotong sebelum quality/safety.

## 18. Risiko dan mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Model mengarang tool/argumen | Aksi salah | Registry allowlist, JSON schema, policy, approval |
| Prompt injection dari web/file | Exfiltration/aksi berbahaya | Tandai untrusted content, capability boundary, jangan izinkan content mengubah policy |
| Shell terlalu kuat | Kerusakan host | Ask-by-default, workspace boundary, sandbox backend, deny patterns |
| Memory menyimpan fakta salah | Personalisasi memburuk | Source/confidence, pin, candidate update, inspect/delete UI |
| Context membengkak | Mahal/lambat | Budget, lazy skill load, truncate artifacts, compaction |
| Provider API berbeda | Fragile adapter | Contract interface dan provider test suite |
| Plugin supply-chain | Remote code execution | MCP isolation, manifest/capabilities, checksums/signing roadmap |
| Self-learning menurunkan kualitas | Skill rusak | Draft/approval, version history, rollback, eval sebelum promote |
| Scope meniru semua Hermes | Proyek tidak selesai | Milestone sempit, parity bukan target MVP |

## 19. Keputusan yang sudah diambil

1. Implementasi baru dimulai dari repository kosong.
2. Bahasa core: **Go**.
3. Produk terminal-first dan local-first.
4. MVP hanya satu user dan satu proses lokal.
5. SQLite adalah source of truth lokal.
6. OpenAI-compatible adalah provider pertama.
7. Tool execution memakai policy + approval + audit.
8. “Berpikir seperti agent” diwujudkan sebagai state machine, plan/checklist, tool feedback, dan verification—not penyimpanan chain-of-thought.
9. Learning loop tidak boleh mengubah security policy atau binary secara mandiri.

## 20. Open questions untuk owner

Hal berikut perlu diputuskan sebelum implementasi di luar foundation:

1. Nama produk final: Nadi atau nama lain?
2. Target penggunaan pertama: coding agent, personal assistant umum, atau automation server?
3. Provider wajib saat MVP: hanya OpenAI-compatible atau juga Anthropic/Ollama?
4. Lisensi: MIT, Apache-2.0, atau proprietary?
5. Apakah v0.1 harus mendukung Windows, atau Windows boleh menyusul?
6. Default execution: host lokal dengan approval atau container-first?
7. Apakah skill hasil belajar auto-save, selalu approval, atau berbeda per profile?
8. Kanal pertama setelah CLI: Telegram, Discord, atau HTTP API?
9. Bahasa UI utama: Inggris, Indonesia, atau bilingual?
10. Batas biaya default dan telemetry: sepenuhnya off atau anonymous opt-in?

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
