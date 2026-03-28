# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Kaktus is a Kowabunga HCI (Hyper-Converged Infrastructure) node agent that manages computing (KVM/libvirt) and storage (Ceph) resources on commodity hardware. It is part of the broader Kowabunga ecosystem.

## Commands

```bash
make all          # mod + fmt + vet + lint + build
make build        # build kaktus binary and ceph plugin
make tests        # run test suite with coverage (go test ./... -count=1)
make fmt          # gofmt
make vet          # go vet
make lint         # golangci-lint
make sec          # gosec security scanner
make vuln         # govulncheck vulnerability scanner
make mod          # go mod download + tidy
make update       # update all Go modules
make clean        # remove binaries and plugins
```

**Single test:**
```bash
go test ./internal/kaktus/... -run TestName -count=1
```

**Packaging:**
```bash
make deb          # Debian package (Ubuntu 24.04 LTS)
make apk          # Alpine package
```

## Architecture

**Entry point:** `cmd/kaktus/main.go` calls `kaktus.Daemonize()`.

**Core components in `internal/kaktus/`:**

- `kaktus.go` — Main agent. Extends `KowabungaAgent` from `github.com/kowabunga-cloud/common`, coordinates all subsystems, loads the Ceph plugin dynamically at runtime.
- `config.go` — YAML configuration parsing (libvirt protocol/address/TLS, Ceph plugin path/monitor settings, global agent settings).
- `libvirt.go` — TCP/TLS connections to libvirt daemons; KVM/QEMU VM lifecycle management; host capability detection (CPU, memory, NUMA).
- `kaktus_services.go` — RPC service layer: exposes node capabilities, instance CRUD, instance state operations, volume management, NFS exports.
- `nfs.go` — REST client (via `go-resty`) for managing NFS exports across multiple backends.
- `disk_image_qcow2.go` / `zstd.go` — QCOW2 image reading with Zstd decompression support.
- `plugins/ceph/` — Dynamically loaded Linux-only plugin (`//go:build linux`):
  - `main.go` — Plugin entry point and type definitions
  - `connection.go` — Ceph cluster connection management
  - `rbd.go` — RADOS Block Device volume operations (create, delete, clone, resize)
  - `fs.go` — CephFS subvolume operations

**Plugin loading:** The Ceph plugin is compiled as a separate shared object and loaded at runtime via Go's `plugin` package. This is why `make build` builds two separate artifacts.

**Key external dependencies:**
- `github.com/kowabunga-cloud/common` — shared agent framework (versioned in sync with kaktus)
- `github.com/digitalocean/go-libvirt` — libvirt wire protocol client
- `github.com/ceph/go-ceph` — Ceph client (requires `librados-dev`, `librbd-dev` at build time)
- `github.com/go-resty/resty/v2` — HTTP client for NFS backend
- `github.com/lima-vm/go-qcow2reader` — QCOW2 format support

## Important Notes

- The Ceph plugin is **Linux-only**. `macos.go` in the plugin directory is a stub. CI runs on Ubuntu 24.04 LTS.
- `common` and `kaktus` share the same version number (currently `0.64.1`). When updating `common`, update both.
- gosec exclusions: G101, G115, G602 (configured in `.github/workflows/sec.yml`).
- Go version required: 1.25.0.
