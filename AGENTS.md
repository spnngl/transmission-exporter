# transmission-exporter

Prometheus exporter for Transmission, forked from metalmatze/transmission-exporter. The root package is the Transmission RPC client; `cmd/transmission-exporter` is the exporter binary.

## Rules

- Write agent-facing documents (every `AGENTS.md`, code comments) per the `writing-for-agents` skill. Each `AGENTS.md` has a `CLAUDE.md` symlink beside it.
- Finish a change when `make check` passes.
- Commit with the `git-commit` skill. Pull requests fail CI unless every commit subject is a Conventional Commit (`type(scope)!: description`; types feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert).
- `dashboards/transmission.json` is a Grafana export: edit it as JSON, or in Grafana and export it again.
- The container image is built in `github.com/spnngl/images` and published as `ghcr.io/spnngl/images/transmission-exporter`; container changes go there.

## Read before you work

- Before adding or changing a metric, read `cmd/transmission-exporter/AGENTS.md`.
- Before editing `.github/` or `renovate.json5`, read `.github/AGENTS.md`.
