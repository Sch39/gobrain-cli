# GoBrain CLI Tool (gob)
  
GoBrain adalah CLI berfokus pada proyek (project-scoped) untuk ekosistem Go. Tujuannya: menjaga lingkungan build dan tooling tetap konsisten di dalam folder proyek, mempermudah inisialisasi template, menjalankan skrip, memasang tool lokal, membuat kode dari template, serta menjalankan pipeline verifikasi — semuanya tanpa bergantung pada tool global pengguna.
  
## Fitur Utama
- Inisialisasi proyek (`gob init`): membuat `gob.yaml`, memilih sumber template (none/preset/url), mengisi placeholder, mengatur `go.mod` dan `toolchain`, serta menyiapkan `.gitignore` terkelola.
- Generator template (`gob make`): merender file dari template menggunakan data dan fungsi bawaan (snake_case, pascal_case, camel_case), mengelola dependensi generator via `go get` ke cache lokal.
- Eksekusi perintah (`gob exec`): menjalankan perintah dengan ENV yang telah diinjeksi untuk proyek (.gob/bin, .gob/mod, GOTOOLCHAIN=auto).
- Skrip proyek (`gob scripts run|list`): mengeksekusi urutan perintah dari `gob.yaml` dengan dukungan chaining `&&`, `cd`, dan `pwd`.
- Manajemen tool (`gob tools install|run|list`): memasang tool ke `.gob/bin` dan menjalankannya dengan PATH yang diinjeksi.
- Verifikasi (`gob verify`): pipeline pemeriksaan/tes berdasarkan langkah-langkah di `gob.yaml`, dengan opsi `fail_fast`.
- Deteksi root proyek: mencari `gob.yaml` pada direktori saat ini dan ke atas agar perintah selalu berjalan pada akar proyek yang tepat.
- Mode debug: flag `--debug` atau `debug: true` pada `gob.yaml` untuk log rinci.
  
## Prasyarat
- Go ≥ 1.21.x (untuk support toolchain).
- Git (untuk meng-clone template via URL).
- Windows/Linux/macOS.
  
## Instalasi & Build
### Instalasi Binary (Release)
- **Windows**
  - Download `gob_windows_amd64.zip` (Intel/AMD) atau `gob_windows_arm64.zip` (ARM).
  - Extract lalu tambahkan ke PATH, contoh: `C:\gobrain\bin`.
- **macOS**
  - Download `gob_darwin_amd64.tar.gz` (Intel) atau `gob_darwin_arm64.tar.gz` (Apple Silicon).
  - Extract dan pindahkan binary ke `/usr/local/bin` atau folder PATH lain.
- **Linux**
  - Download `gob_linux_amd64.tar.gz` atau `gob_linux_arm64.tar.gz`.
  - Extract dan pindahkan binary ke `/usr/local/bin` atau folder PATH lain.
  
Semua paket tersedia di [**Release Page**](https://github.com/Sch39/gobrain-cli-tool/releases).

### Instalasi via go install
```bash
go install github.com/sch39/gobrain-cli/cmd/gob@latest
```

Pastikan `GOBIN` atau `$GOPATH/bin` ada di PATH agar perintah `gob` bisa dijalankan.

### Build Lokal
#### macOS / Linux
  
```bash
go build -trimpath -o bin/gob ./cmd/gob
```
  
#### Windows (PowerShell)
```powershell
go build -trimpath -o bin\gob.exe .\cmd\gob
```
  
### Build dengan Makefile
#### macOS / Linux
- Build dengan Makefile:
  
```bash
  make build
```

#### Windows (Git Bash / MSYS2)
```bash
make build
```

  
## Quick Start
1) Inisialisasi proyek baru:
  
```bash
gob init
```
  
Selama proses, Anda dapat memilih:
- source-type: `none`, `preset` (menggunakan daftar embedded), atau `url` (repo git).
- toolchain: otomatis mendeteksi versi lokal dan menawarkan versi stabil yang lebih baru.
  
2) Menjalankan generator:
  
```bash
gob make list
gob make handler User
gob make handler --name=User
```
  
3) Menjalankan perintah dengan ENV proyek:
  
```bash
gob exec go test ./...
```
  
4) Skrip proyek:
  
```bash
gob scripts list
gob scripts run build
```
  
5) Mengelola tool lokal:
  
```bash
gob tools install
gob tools list
gob tools run golangci-lint --version
```
  
6) Pipeline verifikasi:
  
```bash
gob verify
```
  
## Struktur Konfigurasi (`gob.yaml`)
Konfigurasi dibaca dari `gob.yaml`/`gob.yml` pada root proyek. Bagian-bagian penting:
  
```yaml
version: "1.0.0"
debug: false
project:
  name: my-app
  module: github.com/user/my-app
  toolchain: 1.24.2
meta:
  template_origin: https://github.com/Sch39/gobrain-presets.git
  author: ""
  path: minimal
tools:
  - name: golangci-lint
    pkg: github.com/golangci/golangci-lint/cmd/golangci-lint@latest
scripts:
  build:
    - go build ./...
  test:
    - go test ./...
generators:
  requires:
    - github.com/Masterminds/sprig@latest
  commands:
    handler:
      desc: "Generate HTTP handler"
      args: ["name"]
      templates:
        - src: .template/handler.tmpl
          dest: internal/handlers/{{ pascal_case .Name }}.go
verify:
  fail_fast: true
  pipeline:
    - name: fmt
      run: go fmt ./...
    - name: vet
      run: go vet ./...
layout:
  dirs:
    - internal/handlers
    - internal/services
```
  
Penjelasan singkat:
- `project`: nama, module, dan versi `toolchain` Go yang diharapkan. Toolchain akan diselaraskan ke `go.mod`.
- `meta`: asal template, author, dan path subdirektori template (bila ada).
- `tools`: daftar tool yang akan di-install ke `.gob/bin`.
- `scripts`: peta nama→urutan langkah (perintah shell). Mendukung `&&`, `cd`, `pwd`.
- `generators`: daftar command generator dan template yang akan dirender; `requires` akan di-`go get` ke cache lokal `.gob/mod`.
- `verify`: pipeline langkah verifikasi; jika `fail_fast: true`, proses berhenti pada kegagalan pertama.
- `layout`: direktori yang akan dibuat saat init dari template/preset.
  
## Lingkungan & Toolchain
- Injeksi ENV proyek (`GOBIN`, `GOMODCACHE`, `PATH`, `GOTOOLCHAIN=auto`) memastikan tool dan modul pihak ketiga tersimpan dalam `.gob/` sehingga tidak mengotori global.
- `go.mod` dapat berisi baris `toolchain goX.Y.Z`. Proyek memvalidasi agar `toolchain` ≥ direktif `go` di `go.mod`.
- Deteksi versi Go lokal dan daftar rilis stabil dilakukan saat `init` untuk membantu memilih toolchain.
  
## Preset Template
- Preset bawaan di-embed: lihat `internal/presets/presets.yaml` (berisi preset default dari repo [**gobrain-presets**](https://github.com/Sch39/gobrain-presets) atau **https://github.com/Sch39/gobrain-presets**). Anda bisa menyediakan file preset khusus via ENV `GOB_PRESET_FILE`.
- `source-type: preset` akan menampilkan daftar preset, lalu meng-clone/menyalin isi sesuai `repo` dan `path`.
  
## Keamanan & Token Git
- Untuk repo privat, gunakan ENV:
  - `GOB_GIT_TOKEN` atau `GOB_GITLAB_TOKEN` (Bearer),
  - `GOB_GIT_BASIC_USER` dan `GOB_GIT_BASIC_TOKEN` (Basic),
  - `GOB_GIT_SSH_KEY` untuk kunci SSH,
  - `GOB_GIT_EXTRA_HEADERS` untuk header tambahan.
- Argumen sensitif disanitasi saat logging sehingga token tidak tercetak mentah.
- `.gob/` otomatis dikelola pada `.gitignore` agar artefak lokal tidak ikut ter-commit.
  
## Catatan OS
- Di Windows, `tools run` akan menambahkan `.exe` bila diperlukan saat menjalankan binary dari `.gob/bin`.
  
## Lisensi
- Proyek ini dirilis di bawah **MIT License**.
- Lihat file lisensi di [**LICENSE**](/gobrain/LICENSE) untuk detail lengkap.
  
