# Go

[![CircleCI](https://circleci.com/gh/remove-bg/go.svg?style=shield)](https://circleci.com/gh/remove-bg/go)

## Download

**[Download latest stable release](https://github.com/remove-bg/go/releases)** (Windows, Mac, Linux)



## CLI usage

```
removebg [options] <file>...
```

### API key

To process images you'll need your [remove.bg API key][api-key].

[api-key]: https://www.remove.bg/profile#api-key

To use the API key for all requests you can export the following environment
variable in your shell profile (e.g. `~/.bashrc` / `~/.zshrc`):

```sh
export REMOVE_BG_API_KEY=xyz
```

Alternatively you can specify the API key per command:

```sh
removebg --api-key xyz images/image1.jpg
```

### Resolution limit

Please note that the CLI tool is currently limited to images up to 10 megapixels (larger images are resized). You can get higher-resolution images of up to 25 megapixels via `--format zip` and then converting the ZIP file to the format of your choice (see [here](https://www.remove.bg/api#zip-format)). In the future we plan to extend the CLI tool to compose the final picture automatically.

### Processing a directory of images

#### Saving to the same directory (default)

If you want to remove the background from all the PNG and JPG images in a
directory, and save the transparent images in the same directory:

```sh
removebg images/*.{png,jpg}
```

Given the following input:

```
images/
├── dog.jpg
└── cat.png
```

The result would be:

```
images/
├── dog.jpg
├── cat.png
├── dog-removebg.png
└── cat-removebg.png
```

#### Saving to a different directory (`--output-directory`)

If you want to remove the background from all the PNG and JPG images in a
directory, and save the transparent images in a different directory:

```sh
mkdir processed
removebg --output-directory processed originals/*.{png,jpg}
```

Given the following input:

```
originals/
├── dog.jpg
└── cat.png
```

The result would be:

```
originals/
├── dog.jpg
└── cat.png

processed/
├── dog.png
└── cat.png
```

### CLI options

- `--api-key` or `REMOVE_BG_API_KEY` environment variable (required).

- `--output-directory` (optional) - The output directory for processed images.

- `--reprocess-existing` - Images which have already been processed are skipped
by default to save credits. Specify this flag to force reprocessing.

- `--confirm-batch-over` (default `50`) - Prompt for confirmation before
processing batches over this size. Specify `-1` to disable this safeguard.

#### Image processing options

Please see the [API documentation][api-docs] for further details.

[api-docs]: https://www.remove.bg/api#operations-tag-Background%20Removal

- `--size` (default `auto`)
- `--type`
- `--channels`
- `--bg-color`
- `--format` (default: `png`)

## Development

Prerequisites:

- `go 1.14`
- [`dep`](https://golang.github.io/dep/)

Getting started:

```
git clone git@github.com:remove-bg/go.git $GOPATH/github.com/remove-bg/go
cd $GOPATH/github.com/remove-bg/go
bin/setup
bin/test
```

To build & try out locally:

```
go build -o removebg main.go
./removebg --help
```
