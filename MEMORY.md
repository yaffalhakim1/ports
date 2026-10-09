# Project memory

Facts below are recalled into every session.

## Project identity: `ports` is a Go desktop app built with `github.com/egoist/mygo
Project identity: `ports` is a Go desktop app built with `github.com/egoist/mygo` (see mygo-maintenance skill).

## Vet command: run `go tool mygo vet .` (not plain `go vet`) to lint this project.
Vet command: run `go tool mygo vet .` (not plain `go vet`) to lint this project.

## Repo: public GitHub repo at https://github.com/yaffalhakim1/ports.
Repo: public GitHub repo at https://github.com/yaffalhakim1/ports.

## Appearance control convention: theme selector is a one-click menu (`System / Lig
Appearance control convention: theme selector is a one-click menu (`System / Light / Dark`) at `view_header.go:169` — no cycle/double-click behavior.

## Build artifacts to clean before commit: stray `ports.exe`, `out.txt`, `err.txt` 
Build artifacts to clean before commit: stray `ports.exe`, `out.txt`, `err.txt` are junk and must not be committed.

## CI/release workflow: intentionally omitted; add a GitHub Actions workflow only w
CI/release workflow: intentionally omitted; add a GitHub Actions workflow only when auto-published tagged installers are wanted.
