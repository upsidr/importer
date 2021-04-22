# Importer Installation

## Install with Homebrew

<!-- == export: homebrew-install / begin == -->

The simplest approach is to get Importer via Homebrew.

```bash
$ brew install upsidr/tap/importer
```

You can also find the relevent binary files under [releases](https://github.com/upsidr/importer/releases).

<!-- == export: homebrew-install / end == -->

## Install with Go

<!-- == export: go-get / begin == -->

You can also use Go to install.

```bash
$ go get github.com/upsidr/importer/cmd/importer@v0.0.1-rc2
```

<!-- == export: go-get / end == -->

## Build from Code

<!-- == export: build-from-code / begin == -->

Building from code is straightforward, and is done for CI in this repo.

For building with the latest code, you can run the following command.

```bash
{
    temp_importer_clone=$(mktemp -d)
    pushd "$temp_importer_clone" > /dev/null
    git clone https://github.com/upsidr/importer
    pushd importer > /dev/null

    go build ./cmd/importer/

    popd > /dev/null
    popd > /dev/null
    cp "$temp_importer_clone"/importer/importer .
}
```

You can find how this is done in GitHub Action CI in [.github/workflows/importer-markdown-ci.yaml](https://github.com/upsidr/importer/blob/main/.github/workflows/importer-markdown-ci.yaml)

<!-- == export: build-from-code / end == -->
