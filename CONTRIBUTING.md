# Contributing

## Development
See `README.md` for local development and Docker workflows.

## Releases
Releases are handled by GitHub Actions + GoReleaser and are split per module.

### Tagging
Use module-prefixed tags to trigger the correct pipeline:

- `local-client/vX.Y.Z`
- `signalling-server/vX.Y.Z`

Example:
```bash
git tag local-client/v0.1.0
git push origin local-client/v0.1.0
```

### What happens
On tag push:
1. GitHub Actions selects the matching job in `.github/workflows/release.yml`.
2. GoReleaser runs in the module directory.
3. Release artifacts are built for `linux`, `darwin`, and `windows` on `amd64` and `arm64`.
4. A GitHub Release is created with archives and checksums.

### Config
Each module has its own GoReleaser config:
- `local-client/.goreleaser.yml`
- `signalling-server/.goreleaser.yml`

If you want to change target platforms, archive formats, or add Docker images, update the relevant config file.
